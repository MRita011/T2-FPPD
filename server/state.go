package main

import (
	"bufio"
	"fppd-jogo/common"
	"log"
	"os"
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
		log.Fatalf("Erro ao abrir mapa: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		runes := []rune(line)
		if len(runes) > s.mapWidth {
			s.mapWidth = len(runes)
		}

		for x, ch := range runes {
			switch ch {
			case '▤': // parede (não usado no render atual, mas pode salvar)
				// ignore for now
			case '♣':
				// pode ser vegetação se quiser adicionar
			case '☠':
				s.traps = append(s.traps, common.Element{X: x, Y: y, Symbol: ch})
			case '$':
				s.treasures = append(s.treasures, common.Element{X: x, Y: y, Symbol: ch})
			}
		}
		y++
	}
	s.mapHeight = y

	if err := scanner.Err(); err != nil {
		log.Fatalf("Erro ao ler mapa: %v", err)
	}
}
