package main

import (
	"fmt"
	"fppd-jogo/common"
	"sync"
)

// estado global do jogo no servidor
type StateGame struct {
	sync.Mutex
	nextID  int
	players map[int]*common.Player
}

// cria uma nova instância do estado
func NewStateGame() *StateGame {
	return &StateGame{
		nextID:  0,
		players: make(map[int]*common.Player),
	}
}

// registra um novo jogador
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
