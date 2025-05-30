package main

import (
	"fppd-jogo/common"
	"sync"
)

// Estrutura que mantém o estado do jogo
type StateGame struct {
	sync.Mutex
	nextID  int
	players map[int]*common.Player
}

// Inicializa o estado do jogo
func NewStateGame() *StateGame {
	return &StateGame{
		nextID:  0,
		players: make(map[int]*common.Player),
	}
}

// Método auxiliar para retornar os jogadores como slice
func (s *StateGame) getPlayersSlice() []common.Player {
	players := []common.Player{}
	for _, p := range s.players {
		players = append(players, *p)
	}
	return players
}
