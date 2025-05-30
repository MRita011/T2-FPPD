package main

import (
	"fmt"
	"sync"

	"fppd-jogo/common"
)

// Estado global do jogo
type StateGame struct {
	sync.Mutex
	nextID  int
	players map[int]*common.Player
}

func NewStateGame() *StateGame {
	return &StateGame{
		nextID:  0,
		players: make(map[int]*common.Player),
	}
}

func (s *StateGame) RegisterPlayer(name string) (int, error) {
	s.Lock()
	defer s.Unlock()

	s.nextID++
	id := s.nextID

	s.players[id] = &common.Player{
		ID:   id,
		Name: name,
	}

	fmt.Printf("Jogador registrado: %s (ID %d)\n", name, id)
	return id, nil
}
