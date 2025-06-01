package main

import (
	"bufio"
	"fmt"
	"fppd-jogo/common"
	"log"
	"net/rpc"
	"os"
	"strings"

	"github.com/inancgumus/screen"
)

func registerPlayer(name string) {
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Erro ao conectar:", err)
	}
	defer client.Close()

	req := common.JoinRequest{Name: name}
	var res common.JoinResponse

	err = client.Call("GameServer.RegisterPlayer", &req, &res)
	if err != nil {
		log.Fatal("Erro ao registrar:", err)
	}

	fmt.Printf("Jogador registrado: %s (ID %d)\n", name, res.ID)

	go escutarTeclado(client, res.ID)
	renderMapa(client, res.ID)
}

func escutarTeclado(client *rpc.Client, playerID int) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("→ ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToUpper(input))

		if input == "W" || input == "A" || input == "S" || input == "D" {
			var res common.MoveResponse
			req := common.MoveRequest{PlayerID: playerID, Direction: input}
			err := client.Call("GameServer.MovePlayer", &req, &res)
			if err != nil {
				log.Println("Erro movimentando:", err)
			}
		}
	}
}

func renderMapa(client *rpc.Client, playerID int) {
	var state common.GameState
	err := client.Call("GameServer.GetState", common.StateRequest{PlayerID: playerID}, &state)
	if err != nil {
		log.Println("Erro ao buscar estado:", err)
		return
	}
	renderizar(&state)
}

func renderizar(state *common.GameState) {
	screen.Clear()
	screen.MoveTopLeft()

	for y := 0; y < state.MapHeight; y++ {
		var lineBuilder strings.Builder

		for x := 0; x < state.MapWidth; x++ {
			symbol := state.MapBase[y][x] // fundo do mapa

			for _, e := range state.Treasures {
				if e.X == x && e.Y == y {
					symbol = e.Symbol
					break
				}
			}
			for _, e := range state.Traps {
				if e.X == x && e.Y == y {
					symbol = e.Symbol
					break
				}
			}
			for _, p := range state.Players {
				if p.X == x && p.Y == y {
					symbol = p.Symbol
					break
				}
			}

			lineBuilder.WriteRune(symbol)
		}

		fmt.Println(lineBuilder.String())
	}
}
