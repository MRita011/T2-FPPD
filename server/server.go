package main

import (
	"fppd-jogo/common"
	"log"
	"net"
	"net/rpc"
	"time"
)

// GameServer expõe métodos RPC
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
		Symbol: "P",
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
		Players:   s.state.getPlayersSlice(),
		Traps:     s.state.traps,
		Treasures: s.state.treasures,
	}
	return nil
}

func main() {
	state := NewStateGame()
	state.LoadMapFromFile("mapa.txt") // carrega traps e treasures

	server := &GameServer{state: state}
	rpc.Register(server)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Erro ao escutar:", err)
	}
	log.Println("Servidor iniciado na porta 8080")
	rpc.Accept(listener)
}
