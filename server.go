package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/google/uuid"
)

// armazena globalmente as caixas (tesouros, armadilhas e vazias)
// var caixasGlobal = make(map[Coordenada]TipoCaixa)

// Servidor que gerencia o estado do jogo, sendo as posições dos jogadores e caixas no jogo
type GameServer struct {
	jogadores         map[string]PosicaoJogador // mapa com todas as posições dos jogadores
	processados       map[string]int64          // jogadorID -> último sequence number processado
	caixas            map[Coordenada]TipoCaixa  // mapa de caixas (tesouros, armadilhas e vazias)
	mutex             sync.RWMutex              // trava de sincronização
	estadoPartida     EstadoPartida             // estado atual da partida
	vencedor          string                    // ID do jogador vencedor, se houver
	tesourosRestantes int                       // contador de tesouros restantes
	jogadoresVivos    int                       // contador de jogadores vivos
}

// Serviço RPC para comunicação com clientes
type GameService struct {
	servidor *GameServer // referência ao servidor de posições
}

// Cria um novo servidor de posições de jogadores
func NewGameServer() *GameServer {
	return &GameServer{
		jogadores:         make(map[string]PosicaoJogador),
		processados:       make(map[string]int64),
		caixas:            make(map[Coordenada]TipoCaixa),
		mutex:             sync.RWMutex{},
		estadoPartida:     EstadoAguardando,
		tesourosRestantes: TESOUROS_COUNT,
		jogadoresVivos:    0,
	}
}

// Inicia o servidor RPC na porta especificada
func (gs *GameServer) StartRPC(port string) error {
	service := &GameService{servidor: gs} // cria o serviço RPC
	rpc.Register(service)                 // registra o serviço

	// Inicia o listener TCP na porta especificada
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	log.Printf("Servidor RPC de posições iniciado na porta %s", port)

	// Loop para aceitar conexões
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Erro ao aceitar conexão: %v", err)
			continue
		}
		go rpc.ServeConn(conn) // atende cada conexão em paralelo
	}
}

// Distribui as caixas no mapa seguindo a regra: 12 tesouros, 4 armadilhas, 4 vazias
func (gs *GameServer) distribuirCaixas() error {
	// caixasGlobal = make(map[Coordenada]TipoCaixa)
	gs.caixas = make(map[Coordenada]TipoCaixa)

	var srvJogo Jogo
	if err := CarregarMapa("mapa.txt", &srvJogo); err != nil {
		return fmt.Errorf("falha ao carregar mapa: %v", err)
	}

	// Coleta todas as posições vazias
	var vagas []Coordenada
	for y, linha := range srvJogo.Mapa {
		for x, elem := range linha {
			if elem == Vazio {
				vagas = append(vagas, Coordenada{x, y})
			}
		}
	}

	// Verifica se há posições suficientes
	if len(vagas) < TOTAL_CAIXAS {
		return fmt.Errorf("mapa tem apenas %d posições vazias, mas precisa de %d caixas", len(vagas), TOTAL_CAIXAS)
	}

	// Embaralha as posições
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(vagas), func(i, j int) { vagas[i], vagas[j] = vagas[j], vagas[i] })

	// Distribui as caixas: primeiro tesouros, depois armadilhas, depois vazias
	idx := 0

	// 12 tesouros
	for i := 0; i < TESOUROS_COUNT && idx < len(vagas); i++ {
		gs.caixas[vagas[idx]] = Tesouro
		idx++
	}

	// 4 armadilhas
	for i := 0; i < ARMADILHAS_COUNT && idx < len(vagas); i++ {
		gs.caixas[vagas[idx]] = Armadilha
		idx++
	}

	// 4 vazias
	for i := 0; i < VAZIAS_COUNT && idx < len(vagas); i++ {
		gs.caixas[vagas[idx]] = Vazia
		idx++
	}

	gs.tesourosRestantes = TESOUROS_COUNT
	log.Printf("Caixas distribuídas: %d tesouros, %d armadilhas, %d vazias", TESOUROS_COUNT, ARMADILHAS_COUNT, VAZIAS_COUNT)
	return nil
}

// Verifica as condições de vitória/derrota
func (gs *GameServer) verificarEstadoPartida() {
	if gs.estadoPartida != EstadoJogando {
		return
	}

	// Conta jogadores vivos
	jogadoresVivos := 0
	var jogadorComMaisTesouros string
	maxTesouros := -1

	for _, jogador := range gs.jogadores {
		if jogador.Conectado && !jogador.Morto {
			jogadoresVivos++
			if jogador.Pontuacao > maxTesouros {
				maxTesouros = jogador.Pontuacao
				jogadorComMaisTesouros = jogador.ID
			}
		}
	}

	gs.jogadoresVivos = jogadoresVivos

	// Condições de fim de jogo
	if jogadoresVivos == 0 {
		// Todos morreram - derrota coletiva
		gs.estadoPartida = EstadoFinalizado
		gs.vencedor = ""
		log.Println("Todos os jogadores morreram - Derrota coletiva!")
	} else if gs.tesourosRestantes == 0 {
		// Todos os tesouros foram coletados - jogador com mais tesouros vence
		gs.estadoPartida = EstadoFinalizado
		gs.vencedor = jogadorComMaisTesouros
		log.Printf("Todos os tesouros coletados! Vencedor: %s com %d tesouros", jogadorComMaisTesouros, maxTesouros)
	}
}

// RPC: Jogador se conecta ao servidor de posições
func (gs *GameService) ConectarJogo(req ConectarRequest, reply *ConectarPosicaoResponse) error {
	gs.servidor.mutex.Lock()
	defer gs.servidor.mutex.Unlock()

	// Se o jogo já terminou, não permite novos jogadores
	if gs.servidor.estadoPartida == EstadoFinalizado {
		return fmt.Errorf("jogo já finalizado")
	}

	// Cria um novo ID para o jogador
	jogadorID := uuid.New().String()

	// Determina a cor do jogador baseado na quantidade atual de jogadores
	corIndex := len(gs.servidor.jogadores) % len(CoresJogadores)

	// Encontra uma posição inicial (simplificada)
	posX, posY := 5+len(gs.servidor.jogadores)*2, 5+len(gs.servidor.jogadores)*2

	// Cria novo jogador com as informações básicas
	novoJogador := PosicaoJogador{
		ID:        jogadorID,
		Nome:      req.Nome,
		PosX:      posX,
		PosY:      posY,
		Cor:       CoresJogadores[corIndex],
		Simbolo:   '♟',
		Conectado: true,
		Pontuacao: 0,
		Morto:     false,
	}

	// Adiciona o jogador ao mapa de posições
	gs.servidor.jogadores[jogadorID] = novoJogador
	gs.servidor.processados[jogadorID] = 0

	// Se é o primeiro jogador ou o jogo ainda não começou, redistribui as caixas
	if gs.servidor.estadoPartida == EstadoAguardando {
		if err := gs.servidor.distribuirCaixas(); err != nil {
			return fmt.Errorf("erro ao distribuir caixas: %v", err)
		}
		gs.servidor.estadoPartida = EstadoJogando
		log.Println("Partida iniciada!")
	}

	gs.servidor.jogadoresVivos = len(gs.servidor.jogadores)

	// Prepara a resposta para o cliente
	posicoes := PosicoesJogadores{
		Jogadores:         gs.servidor.copiarPosicoes(),
		JogadorID:         jogadorID,
		UltimoProcessado:  0,
		Caixas:            gs.servidor.copiarCaixas(),
		EstadoPartida:     gs.servidor.estadoPartida,
		Vencedor:          gs.servidor.vencedor,
		TesourosRestantes: gs.servidor.tesourosRestantes,
		JogadoresVivos:    gs.servidor.jogadoresVivos,
	}

	reply.JogadorID = jogadorID
	reply.Posicoes = posicoes

	log.Printf("Jogador %s conectado (%s)", novoJogador.Nome, jogadorID)
	return nil
}

// RPC: Jogador tenta se mover e recebe posições atualizadas
func (gs *GameService) Mover(req MoverRequest, reply *PosicoesJogadores) error {
	gs.servidor.mutex.Lock()
	defer gs.servidor.mutex.Unlock()

	// Se o jogo terminou, não permite movimentos
	if gs.servidor.estadoPartida == EstadoFinalizado {
		*reply = gs.servidor.criarResposta(req.JogadorID, req.SequenceNumber)
		return nil
	}

	// Verifica se esse comando já foi processado
	if gs.servidor.processados[req.JogadorID] >= req.SequenceNumber {
		*reply = gs.servidor.criarResposta(req.JogadorID, gs.servidor.processados[req.JogadorID])
		return nil
	}

	// Busca o jogador que está se movendo
	jogador, existe := gs.servidor.jogadores[req.JogadorID]
	if !existe || !jogador.Conectado || jogador.Morto {
		*reply = gs.servidor.criarResposta(req.JogadorID, req.SequenceNumber)
		return nil
	}

	// Calcula direção do movimento
	dx, dy := 0, 0
	switch req.Tecla {
	case 'w':
		dy = -1
	case 'a':
		dx = -1
	case 's':
		dy = 1
	case 'd':
		dx = 1
	default:
		*reply = gs.servidor.criarResposta(req.JogadorID, gs.servidor.processados[req.JogadorID])
		return nil
	}

	// Calcula nova posição
	nx, ny := jogador.PosX+dx, jogador.PosY+dy

	// Faz o movimento
	jogador.PosX = nx
	jogador.PosY = ny
	gs.servidor.jogadores[req.JogadorID] = jogador

	// Atualiza o processamento
	gs.servidor.processados[req.JogadorID] = req.SequenceNumber

	// Prepara a resposta com posições atualizadas
	*reply = gs.servidor.criarResposta(req.JogadorID, req.SequenceNumber)
	return nil
}

// RPC: Cliente solicita posições atuais de todos jogadores
func (gs *GameService) ObterPosicoes(jogadorID string, reply *PosicoesJogadores) error {
	gs.servidor.mutex.RLock()
	defer gs.servidor.mutex.RUnlock()

	// Prepara a resposta com todas as posições atuais
	*reply = gs.servidor.criarResposta(jogadorID, gs.servidor.processados[jogadorID])
	return nil
}

// Interagir com caixa
func (gs *GameService) InteragirCaixa(req InteragirRequest, reply *InteragirResponse) error {
	gs.servidor.mutex.Lock()
	defer gs.servidor.mutex.Unlock()

	// Se o jogo terminou, não permite interações
	if gs.servidor.estadoPartida == EstadoFinalizado {
		reply.Caixas = gs.servidor.copiarCaixas()
		reply.EstadoPartida = gs.servidor.estadoPartida
		reply.Vencedor = gs.servidor.vencedor
		return nil
	}

	// encontra posição do jogador
	pj, existe := gs.servidor.jogadores[req.JogadorID]
	if !existe || pj.Morto {
		return fmt.Errorf("jogador não encontrado ou já morreu")
	}

	// Apenas verifica a posição atual do jogador
	posAtual := Coordenada{pj.PosX, pj.PosY}
	if t, ok := gs.servidor.caixas[posAtual]; ok {
		reply.Tipo = t
		// Remove a caixa do mapa após a interação
		delete(gs.servidor.caixas, posAtual)

		// Atualiza o jogador baseado no tipo de caixa
		switch t {
		case Tesouro:
			pj.Pontuacao++
			gs.servidor.tesourosRestantes--
			log.Printf("Jogador %s encontrou um tesouro! Pontuação: %d", pj.Nome, pj.Pontuacao)
		case Armadilha:
			pj.Morto = true
			reply.GameOver = true
			log.Printf("Jogador %s ativou uma armadilha e morreu!", pj.Nome)
		case Vazia:
			log.Printf("Jogador %s abriu uma caixa vazia", pj.Nome)
		}
		// Atualiza o jogador no servidor
		gs.servidor.jogadores[req.JogadorID] = pj
		reply.Pontuacao = pj.Pontuacao
	} else {
		// Se não encontrou caixa na posição exata, a resposta é vazia.
		reply.Tipo = ""
		log.Printf("Jogador %s tentou interagir em (%d,%d), mas não há nada aqui.", pj.Nome, posAtual.X, posAtual.Y)
	}

	// Verifica condições de fim de jogo
	gs.servidor.verificarEstadoPartida()

	// Prepara resposta
	reply.Caixas = gs.servidor.copiarCaixas() // Agora copia as caixas atualizadas do servidor
	reply.EstadoPartida = gs.servidor.estadoPartida
	reply.Vencedor = gs.servidor.vencedor
	reply.TesourosRestantes = gs.servidor.tesourosRestantes
	reply.JogadoresVivos = gs.servidor.jogadoresVivos

	return nil
}

// RPC: Jogador se desconecta
func (gs *GameService) Desconectar(jogadorID string, reply *bool) error {
	gs.servidor.mutex.Lock()
	defer gs.servidor.mutex.Unlock()

	if jogador, existe := gs.servidor.jogadores[jogadorID]; existe {
		log.Printf("Jogador %s (%s) desconectado", jogador.Nome, jogadorID)
		delete(gs.servidor.jogadores, jogadorID)   // remove o jogador do mapa
		delete(gs.servidor.processados, jogadorID) // remove o processamento

		// Verifica se ainda há jogadores e se o jogo deve continuar
		gs.servidor.verificarEstadoPartida()
	}

	*reply = true
	return nil
}

// Copia as posições dos jogadores para evitar problemas de concorrência
func (gs *GameServer) copiarPosicoes() map[string]PosicaoJogador {
	copia := make(map[string]PosicaoJogador)
	for id, jogador := range gs.jogadores {
		if jogador.Conectado {
			copia[id] = jogador
		}
	}
	return copia
}

// Copia o mapa de caixas para evitar problemas de concorrência
func (gs *GameServer) copiarCaixas() map[Coordenada]TipoCaixa {
	copia := make(map[Coordenada]TipoCaixa)
	for coord, tipo := range gs.caixas {
		copia[coord] = tipo
	}
	return copia
}

// Cria uma resposta padronizada
func (gs *GameServer) criarResposta(jogadorID string, sequenceNumber int64) PosicoesJogadores {
	return PosicoesJogadores{
		Jogadores:         gs.copiarPosicoes(),
		JogadorID:         jogadorID,
		UltimoProcessado:  sequenceNumber,
		Caixas:            gs.copiarCaixas(),
		EstadoPartida:     gs.estadoPartida,
		Vencedor:          gs.vencedor,
		TesourosRestantes: gs.tesourosRestantes,
		JogadoresVivos:    gs.jogadoresVivos,
	}
}
