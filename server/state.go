package main

import (
	"bufio"
	"fppd-jogo/common"
	"log"
	"os"
	"strings"
	"sync"
)

type StateGame struct {
	sync.Mutex
	nextID    int
	players   map[int]*common.Player
	traps     []common.Element
	treasures []common.Element
	mapWidth  int
	mapHeight int
}

func NewStateGame() *StateGame {
	return &StateGame{
		nextID:  0,
		players: make(map[int]*common.Player),
	}
}

func (s *StateGame) getPlayersSlice() []common.Player {
	players := []common.Player{}
	for _, p := range s.players {
		players = append(players, *p)
	}
	return players
}

func (s *StateGame) LoadMapFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Erro ao abrir mapa: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		s.mapWidth = len(line)
		for x, char := range strings.TrimSpace(line) {
			switch char {
			case 'X':
				s.traps = append(s.traps, common.Element{X: x, Y: y, Symbol: "X"})
			case '$':
				s.treasures = append(s.treasures, common.Element{X: x, Y: y, Symbol: "$"})
			}
		}
		y++
	}
	s.mapHeight = y

	if err := scanner.Err(); err != nil {
		log.Printf("Erro ao ler mapa: %v", err)
	}
}
