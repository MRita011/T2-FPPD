package main

import (
	"fppd-jogo/common"
	"log"
	"net"
	"net/rpc"
	"time"
)

type GameServer struct {
	state *StateGame
}

func (s *GameServer) RegisterPlayer(req *common.JoinRequest, res *common.JoinResponse) error {
	s.state.Lock()
	defer s.state.Unlock()

	s.state.nextID++
	id := s.state.nextID

	s.state.players[id] = &common.Player{
		ID:     id,
		Name:   req.Name,
		X:      1,
		Y:      1,
		Symbol: '☺', // ← símbolo como RUNE
	}

	log.Printf("[%s] Jogador registrado: %s (ID %d)\n", time.Now().Format(time.RFC3339), req.Name, id)

	*res = common.JoinResponse{
		ID:     id,
		Player: *s.state.players[id],
	}
	return nil
}

func (s *GameServer) GetState(req common.StateRequest, res *common.GameState) error {
	s.state.Lock()
	defer s.state.Unlock()

	log.Printf("[%s] Estado solicitado por jogador %d\n", time.Now().Format(time.RFC3339), req.PlayerID)

	*res = common.GameState{
		MapWidth:  s.state.mapWidth,
		MapHeight: s.state.mapHeight,
		MapBase:   s.state.mapaBase,
		Players:   s.state.getPlayersSlice(),
		Traps:     s.state.traps,
		Treasures: s.state.treasures,
	}

	return nil
}

func (s *GameServer) MovePlayer(req *common.MoveRequest, res *common.MoveResponse) error {
	s.state.Lock()
	defer s.state.Unlock()

	player, exists := s.state.players[req.PlayerID]
	if !exists {
		*res = common.MoveResponse{Success: false, Message: "Jogador não encontrado"}
		return nil
	}

	newX, newY := player.X, player.Y

	switch req.Direction {
	case "W":
		newY--
	case "S":
		newY++
	case "A":
		newX--
	case "D":
		newX++
	default:
		*res = common.MoveResponse{Success: false, Message: "Direção inválida"}
		return nil
	}

	if newX < 0 || newX >= s.state.mapWidth || newY < 0 || newY >= s.state.mapHeight {
		*res = common.MoveResponse{Success: false, Message: "Movimento fora dos limites"}
		return nil
	}

	player.X = newX
	player.Y = newY

	*res = common.MoveResponse{Success: true, Message: "Movido com sucesso"}
	return nil
}

func main() {
	state := NewStateGame()
	state.LoadMapFromFile("mapa.txt")

	server := &GameServer{state: state}
	rpc.Register(server)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Erro ao escutar:", err)
	}
	log.Println("Servidor iniciado na porta 8080")
	rpc.Accept(listener)
}
