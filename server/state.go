package main

import (
	"fmt"
	"sync"

	"fppd-jogo/common" // pacote com as structs comuns usadas no jogo
)

// Estado global do jogo
type StateGame struct {
	sync.Mutex                        // mutex pra garantir que só uma goroutine mexa ao mesmo tempo
	nextID     int                    // contador para gerar ids únicos pros jogadores
	players    map[int]*common.Player // mapa que guarda os jogadores cadastrados, indexados pelo id
}

// função que cria um estado novo (sem jogadores ainda)
func NewStateGame() *StateGame {
	return &StateGame{
		nextID:  0,                            // começa o id do zero
		players: make(map[int]*common.Player), // cria o mapa vazio pros jogadores
	}
}

// método para registrar um jogador novo
func (s *StateGame) RegisterPlayer(name string) (int, error) {
	s.Lock()         // não deixa outra goroutine mexer enquanto registra
	defer s.Unlock() // destrava o mutex no fim dessa função

	s.nextID++ // incrementa o id p/ o próximo jogador
	id := s.nextID

	// cria um jogador novo com o id e nome e guarda no mapa
	s.players[id] = &common.Player{
		ID:   id,
		Name: name,
	}

	fmt.Printf("Jogador registrado: %s (ID %d)\n", name, id)
	return id, nil // retorna o id do jogador e nil pra erro (tudo certo)
}
