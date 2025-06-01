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
	mapaBase  [][]rune
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
	var mapa [][]rune
	lines := []string{}
	maxWidth := 0

	file, err := os.Open("server/mapa.txt")
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
		lines = append(lines, line)
	}

	s.mapHeight = len(lines)
	s.mapWidth = maxWidth

	for _, line := range lines {
		runes := []rune(line)
		for len(runes) < maxWidth {
			runes = append(runes, ' ') // completa
		}
		mapa = append(mapa, runes)
	}

	s.mapaBase = mapa // salva o mapa base

	// depois: parseia armadilhas e tesouros
	for y, row := range mapa {
		for x, ch := range row {
			switch ch {
			case '☠':
				s.traps = append(s.traps, common.Element{X: x, Y: y, Symbol: ch})
			case '$':
				s.treasures = append(s.treasures, common.Element{X: x, Y: y, Symbol: ch})
			}
		}
	}
}
