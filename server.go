package main

import (
	"log"
	"net"
	"net/rpc"
	"sync"

	"github.com/google/uuid"
)

// Servidor que gerencia apenas as posições dos jogadores
type GameServer struct {
	jogadores   map[string]PosicaoJogador // mapa com todas as posições dos jogadores
	processados map[string]int64          // jogadorID -> último sequence number processado
	mutex       sync.RWMutex              // trava de sincronização
}

// Serviço RPC para comunicação com clientes
type GameService struct {
	servidor *GameServer // referência ao servidor de posições
}

// Cria um novo servidor de posições de jogadores
func NewGameServer() *GameServer {
	return &GameServer{
		jogadores:   make(map[string]PosicaoJogador),
		processados: make(map[string]int64),
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

// RPC: Jogador se conecta ao servidor de posições
func (gs *GameService) ConectarJogo(req ConectarRequest, reply *ConectarPosicaoResponse) error {
	gs.servidor.mutex.Lock()
	defer gs.servidor.mutex.Unlock()

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
		Simbolo:   '☺',
		Conectado: true,
	}

	// Adiciona o jogador ao mapa de posições
	gs.servidor.jogadores[jogadorID] = novoJogador
	gs.servidor.processados[jogadorID] = 0

	// Prepara a resposta para o cliente
	posicoes := PosicoesJogadores{
		Jogadores:        gs.servidor.copiarPosicoes(),
		JogadorID:        jogadorID,
		UltimoProcessado: 0,
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

	// Verifica se esse comando já foi processado
	if gs.servidor.processados[req.JogadorID] >= req.SequenceNumber {
		*reply = PosicoesJogadores{
			Jogadores:        gs.servidor.copiarPosicoes(),
			JogadorID:        req.JogadorID,
			UltimoProcessado: gs.servidor.processados[req.JogadorID],
		}
		return nil
	}

	// Busca o jogador que está se movendo
	jogador, existe := gs.servidor.jogadores[req.JogadorID]
	if !existe || !jogador.Conectado {
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
		*reply = PosicoesJogadores{
			Jogadores:        gs.servidor.copiarPosicoes(),
			JogadorID:        req.JogadorID,
			UltimoProcessado: gs.servidor.processados[req.JogadorID],
		}
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
	*reply = PosicoesJogadores{
		Jogadores:        gs.servidor.copiarPosicoes(),
		JogadorID:        req.JogadorID,
		UltimoProcessado: req.SequenceNumber,
	}

	return nil
}

// RPC: Cliente solicita posições atuais de todos jogadores
func (gs *GameService) ObterPosicoes(jogadorID string, reply *PosicoesJogadores) error {
	gs.servidor.mutex.RLock()
	defer gs.servidor.mutex.RUnlock()

	// Prepara a resposta com todas as posições atuais
	*reply = PosicoesJogadores{
		Jogadores:        gs.servidor.copiarPosicoes(),
		JogadorID:        jogadorID,
		UltimoProcessado: gs.servidor.processados[jogadorID],
	}

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

// RPC: Jogador se desconecta
func (gs *GameService) Desconectar(jogadorID string, reply *bool) error {
	gs.servidor.mutex.Lock()
	defer gs.servidor.mutex.Unlock()

	if jogador, existe := gs.servidor.jogadores[jogadorID]; existe {
		log.Printf("Jogador %s (%s) desconectado", jogador.Nome, jogadorID)
		delete(gs.servidor.jogadores, jogadorID)   // remove o jogador do mapa
		delete(gs.servidor.processados, jogadorID) // remove o processamento
	}

	*reply = true
	return nil
}
