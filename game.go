package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

// GameManager gerencia o estado do jogo localmente
type GameManager struct {
	jogo                *Jogo
	jogadorID           string                    // ID do jogador local
	jogadoresRemotos    map[string]PosicaoJogador // Jogadores remotos
	comandosProcessados map[string]int64          // jogadorID -> último sequence number processado
	caixas              map[Coordenada]TipoCaixa  // Caixas no mapa
	mutex               sync.RWMutex
}

// Cria um novo gerenciador de jogo local
func NewGameManager() *GameManager {
	return &GameManager{
		jogadoresRemotos:    make(map[string]PosicaoJogador),
		comandosProcessados: make(map[string]int64),
		caixas:              make(map[Coordenada]TipoCaixa),
	}
}

func (gm *GameManager) AtualizarCaixas(novas map[Coordenada]TipoCaixa) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	gm.caixas = novas
}

func (gm *GameManager) InteragirCaixa(posX, posY int) (TipoCaixa, bool) {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()

	coord := Coordenada{X: posX, Y: posY}
	log.Printf("Interagindo com caixa em (%d, %d)", posX, posY)
	tipoCaixa, existe := gm.caixas[coord]

	if existe {
		if tipoCaixa == Armadilha {
			gm.mutex.Lock()
			gm.jogo.GameOver = true
			gm.jogo.StatusMsg = "Você abriu uma armadilha! GAME OVER!"
			gm.mutex.Unlock()
		} else if tipoCaixa == Tesouro {
			gm.mutex.Lock()
			gm.jogo.StatusMsg = "Você encontrou um tesouro!"
			gm.mutex.Unlock()
		}

		return tipoCaixa, true
	}

	// Verifica caixas adjacentes
	adjacentes := []Coordenada{
		{X: posX + 1, Y: posY},
		{X: posX - 1, Y: posY},
		{X: posX, Y: posY + 1},
		{X: posX, Y: posY - 1},
	}

	for _, adj := range adjacentes {
		if tipo, ok := gm.caixas[adj]; ok {
			return tipo, true
		}
	}

	return "", false
}

// Inicializa o jogo local com o mapa fornecido
func (gm *GameManager) InicializarJogo(mapaFile string) error {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if gm.jogo != nil {
		return nil // Jogo já inicializado
	}

	id := uuid.New().String()
	jogo := &Jogo{
		ID:             id,
		Jogadores:      make(map[string]*Jogador),
		UltimoVisitado: Vazio,
		StatusMsg:      "Jogo multiplayer iniciado",
	}

	if err := CarregarMapa(mapaFile, jogo); err != nil {
		return fmt.Errorf("erro ao carregar mapa: %v", err)
	}

	gm.jogo = jogo
	return nil
}

// Cria um jogador local com as informações recebidas do servidor
func (gm *GameManager) CriarJogadorLocal(jogadorID string, nome string, posX, posY int, cor Cor) *Jogador {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Verifica se o jogo foi inicializado
	if gm.jogo == nil {
		gm.InicializarJogo("mapa.txt") // Usa o mapa padrão se não tiver sido inicializado
	}

	jogador := &Jogador{
		ID:        jogadorID,
		Nome:      nome,
		PosX:      posX,
		PosY:      posY,
		Cor:       cor,
		Simbolo:   '♟',
		Conectado: true,
	}

	// Armazena o jogador no jogo local
	gm.jogo.Jogadores[jogadorID] = jogador
	gm.jogadorID = jogadorID
	gm.comandosProcessados[jogadorID] = 0

	return jogador
}

// Atualiza as posições dos jogadores remotos no jogo local
func (gm *GameManager) AtualizarJogadoresRemotos(posicoes map[string]PosicaoJogador) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Guarda as posições remotas
	gm.jogadoresRemotos = posicoes

	// Se o jogo não estiver inicializado, não faz nada
	if gm.jogo == nil {
		return
	}

	// Atualiza os jogadores no jogo local
	for id, posicao := range posicoes {
		// Não atualiza o jogador local
		if id == gm.jogadorID {
			continue
		}

		// Cria ou atualiza o jogador remoto no jogo local
		jogador, existe := gm.jogo.Jogadores[id]
		if !existe {
			jogador = &Jogador{
				ID:        posicao.ID,
				Nome:      posicao.Nome,
				Cor:       posicao.Cor,
				Simbolo:   posicao.Simbolo,
				Conectado: posicao.Conectado,
			}
			gm.jogo.Jogadores[id] = jogador
		}

		// Atualiza a posição do jogador remoto
		jogador.PosX = posicao.PosX
		jogador.PosY = posicao.PosY
		jogador.Conectado = posicao.Conectado
	}
}

// Atualiza a posição do jogador local no jogo
func (gm *GameManager) MoverJogadorLocal(tecla rune) (*EstadoJogo, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if gm.jogo == nil || gm.jogadorID == "" {
		return nil, fmt.Errorf("jogo não inicializado ou jogador local não definido")
	}

	jogador, existe := gm.jogo.Jogadores[gm.jogadorID]
	if !existe || !jogador.Conectado {
		return nil, fmt.Errorf("jogador não encontrado ou desconectado")
	}

	dx, dy := 0, 0
	switch tecla {
	case 'w':
		dy = -1
	case 'a':
		dx = -1
	case 's':
		dy = 1
	case 'd':
		dx = 1
	default:
		return gm.obterEstadoAtual(), nil
	}

	nx, ny := jogador.PosX+dx, jogador.PosY+dy
	if gm.podeMover(nx, ny, gm.jogadorID) {
		jogador.PosX = nx
		jogador.PosY = ny
		gm.jogo.StatusMsg = fmt.Sprintf("Você moveu para (%d, %d)", nx, ny)
	} else {
		gm.jogo.StatusMsg = "Movimento bloqueado!"
	}

	return gm.obterEstadoAtual(), nil
}

// Verifica se o jogador pode se mover para a posição especificada
func (gm *GameManager) podeMover(x, y int, jogadorID string) bool {
	if gm.jogo == nil || gm.jogo.Mapa == nil {
		return false
	}

	// Verifica limites do mapa
	if y < 0 || y >= len(gm.jogo.Mapa) || x < 0 || x >= len(gm.jogo.Mapa[y]) {
		return false
	}

	// Verifica colisão com elementos do mapa
	if gm.jogo.Mapa[y][x].Tangivel {
		return false
	}

	// Verifica colisão com outros jogadores
	for id, jogador := range gm.jogo.Jogadores {
		if id != jogadorID && jogador.PosX == x && jogador.PosY == y && jogador.Conectado {
			return false
		}
	}

	return true
}

// Obtém o estado atual do jogo local
func (gm *GameManager) ObterEstado() *EstadoJogo {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()
	return gm.obterEstadoAtual()
}

// obtém o estado atual do jogo local (sem lock)
func (gm *GameManager) obterEstadoAtual() *EstadoJogo {
	// Caso o jogo ainda não tenha sido inicializado
	if gm.jogo == nil {
		return &EstadoJogo{
			Jogadores: make(map[string]*Jogador),
			StatusMsg: "Jogo não inicializado",
			Caixas:    make(map[Coordenada]TipoCaixa),
		}
	}

	// Caso o jogo esteja carregado, retorna mapa, jogadores, status e caixas
	return &EstadoJogo{
		Mapa:      gm.jogo.Mapa,
		Jogadores: gm.copiarJogadores(),
		StatusMsg: gm.jogo.StatusMsg,
		Caixas:    gm.caixas,
	}
}

// Copia os jogadores para evitar problemas de concorrência
func (gm *GameManager) copiarJogadores() map[string]*Jogador {
	copia := make(map[string]*Jogador)
	for id, jogador := range gm.jogo.Jogadores {
		if jogador.Conectado {
			copia[id] = &Jogador{
				ID:        jogador.ID,
				Nome:      jogador.Nome,
				PosX:      jogador.PosX,
				PosY:      jogador.PosY,
				Cor:       jogador.Cor,
				Simbolo:   jogador.Simbolo,
				Conectado: jogador.Conectado,
			}
		}
	}
	return copia
}

// Carrega o mapa do arquivo especificado
func CarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	// Ensure the map is initialized
	jogo.Mapa = make([][]Elemento, 0)

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for _, ch := range linha {
			e := Vazio
			switch ch {
			// case '■':
			// 	e = Caixa
			case '▤':
				e = Parede
			case '♙':
				e = Inimigo
			case '♣':
				e = Vegetacao
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	return scanner.Err()
}
