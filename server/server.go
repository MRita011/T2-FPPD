package main

import (
	"fmt"
	"fppd-jogo/common"
	"log"
	"net"
	"net/rpc"
)

// GameServer é o objeto RPC com referência ao estado do jogo
type GameServer struct {
	state *StateGame
}

// Método RPC para registrar jogador
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

	fmt.Printf("Jogador registrado: %s (ID %d)\n", req.Name, id)

	*res = common.JoinResponse{
		ID:     id,
		Player: *s.state.players[id],
	}
	return nil
}

// Método RPC para retornar o estado do jogo
func (s *GameServer) GetState(req common.StateRequest, res *common.GameState) error {
	s.state.Lock()
	defer s.state.Unlock()

	*res = common.GameState{
		MapWidth:  20,
		MapHeight: 10,
		Players:   s.state.getPlayersSlice(),
		Traps:     []common.Element{}, // Exemplo: preenchido depois
		Treasures: []common.Element{}, // Idem
	}
	return nil
}

func main() {
	server := &GameServer{state: NewStateGame()}

	rpc.Register(server)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Erro ao escutar:", err)
	}
	log.Println("Servidor iniciado na porta 8080")
	rpc.Accept(listener)
}
