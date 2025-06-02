package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

type GameManager struct {
	jogo                *Jogo
	comandosProcessados map[string]int64 // jogadorID -> último sequence number processado
	mutex               sync.RWMutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		comandosProcessados: make(map[string]int64),
	}
}

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

func (gm *GameManager) ConectarJogador(mapaFile string) (*Jogador, *EstadoJogo, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Inicializar jogo se necessário
	if gm.jogo == nil {
		if err := gm.inicializarJogoInterno(mapaFile); err != nil {
			return nil, nil, err
		}
	}

	// Criar novo jogador
	jogadorID := uuid.New().String()
	corIndex := len(gm.jogo.Jogadores) % len(CoresJogadores)

	// Encontrar posição inicial livre
	posX, posY := gm.encontrarPosicaoLivre()

	jogador := &Jogador{
		ID:        jogadorID,
		Nome:      fmt.Sprintf("Jogador%d", len(gm.jogo.Jogadores)+1),
		PosX:      posX,
		PosY:      posY,
		Cor:       CoresJogadores[corIndex],
		Simbolo:   '☺',
		Conectado: true,
	}

	gm.jogo.Jogadores[jogadorID] = jogador
	gm.comandosProcessados[jogadorID] = 0

	estado := &EstadoJogo{
		Mapa:      gm.jogo.Mapa,
		Jogadores: gm.copiarJogadores(),
		StatusMsg: fmt.Sprintf("%s conectou-se ao jogo", jogador.Nome),
	}

	gm.jogo.StatusMsg = estado.StatusMsg
	return jogador, estado, nil
}

func (gm *GameManager) inicializarJogoInterno(mapaFile string) error {
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

func (gm *GameManager) encontrarPosicaoLivre() (int, int) {
	if gm.jogo == nil || len(gm.jogo.Mapa) == 0 {
		return 1, 1
	}

	// Procurar por uma posição livre no mapa
	for y := 1; y < len(gm.jogo.Mapa)-1; y++ {
		for x := 1; x < len(gm.jogo.Mapa[y])-1; x++ {
			if !gm.jogo.Mapa[y][x].Tangivel && !gm.posicaoOcupada(x, y) {
				return x, y
			}
		}
	}
	return 1, 1 // Fallback
}

func (gm *GameManager) posicaoOcupada(x, y int) bool {
	for _, jogador := range gm.jogo.Jogadores {
		if jogador.PosX == x && jogador.PosY == y && jogador.Conectado {
			return true
		}
	}
	return false
}

func (gm *GameManager) MoverJogador(jogadorID string, tecla rune, sequenceNumber int64) (*EstadoJogo, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Verificar execução única
	if gm.comandosProcessados[jogadorID] >= sequenceNumber {
		// Comando já processado, retornar estado atual
		return gm.obterEstadoAtual(), nil
	}

	jogador, exists := gm.jogo.Jogadores[jogadorID]
	if !exists || !jogador.Conectado {
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
	if gm.podeMover(nx, ny, jogadorID) {
		jogador.PosX, jogador.PosY = nx, ny
		gm.jogo.StatusMsg = fmt.Sprintf("%s moveu para (%d, %d)", jogador.Nome, nx, ny)
		gm.comandosProcessados[jogadorID] = sequenceNumber
	} else {
		gm.jogo.StatusMsg = fmt.Sprintf("%s: movimento bloqueado!", jogador.Nome)
	}

	return gm.obterEstadoAtual(), nil
}

func (gm *GameManager) podeMover(x, y int, jogadorID string) bool {
	if y < 0 || y >= len(gm.jogo.Mapa) || x < 0 || x >= len(gm.jogo.Mapa[y]) {
		return false
	}
	if gm.jogo.Mapa[y][x].Tangivel {
		return false
	}
	// Verificar se há outro jogador na posição
	for id, jogador := range gm.jogo.Jogadores {
		if id != jogadorID && jogador.PosX == x && jogador.PosY == y && jogador.Conectado {
			return false
		}
	}
	return true
}

func (gm *GameManager) ObterEstado() *EstadoJogo {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()
	return gm.obterEstadoAtual()
}

func (gm *GameManager) obterEstadoAtual() *EstadoJogo {
	if gm.jogo == nil {
		return &EstadoJogo{
			Jogadores: make(map[string]*Jogador),
			StatusMsg: "Jogo não inicializado",
		}
	}

	return &EstadoJogo{
		Mapa:      gm.jogo.Mapa,
		Jogadores: gm.copiarJogadores(),
		StatusMsg: gm.jogo.StatusMsg,
	}
}

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

func (gm *GameManager) DesconectarJogador(jogadorID string) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if jogador, exists := gm.jogo.Jogadores[jogadorID]; exists {
		jogador.Conectado = false
		gm.jogo.StatusMsg = fmt.Sprintf("%s desconectou-se", jogador.Nome)
	}
}

func CarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for _, ch := range linha {
			e := Vazio
			switch ch {
			case '▤':
				e = Parede
			case '☠':
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

func PodeMover(jogo *Jogo, x, y int) bool {
	if y < 0 || y >= len(jogo.Mapa) || x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}
	return !jogo.Mapa[y][x].Tangivel
}
