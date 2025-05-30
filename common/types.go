package common

import (
	"fmt"
	"sync"
)

type JoinRequest struct {
	Name string
}

type JoinResponse struct {
	ID int
}

type Player struct {
	ID   int
	Name string
}

// Estado global do jogo no servidor
type StateGame struct {
	sync.Mutex
	NextID  int
	Players map[int]*Player
}

func NewStateGame() *StateGame {
	return &StateGame{
		NextID:  0,
		Players: make(map[int]*Player),
	}
}

func (s *StateGame) RegisterPlayer(name string) (int, error) {
	s.Lock()
	defer s.Unlock()

	s.NextID++
	id := s.NextID

	s.Players[id] = &Player{
		ID:   id,
		Name: name,
	}

	fmt.Printf("Jogador registrado: %s (ID %d)\n", name, id)
	return id, nil
}
