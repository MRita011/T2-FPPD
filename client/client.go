package main

import (
	"bufio"
	"fmt"
	"fppd-jogo/common"
	"log"
	"net/rpc"
	"os"
	"strings"
	"time"
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
	atualizarMapa(client, res.ID)
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

func atualizarMapa(client *rpc.Client, playerID int) {
	for {
		var state common.GameState
		err := client.Call("GameServer.GetState", common.StateRequest{PlayerID: playerID}, &state)
		if err != nil {
			log.Println("Erro ao buscar estado:", err)
			time.Sleep(1 * time.Second)
			continue
		}
		renderizar(&state)
		time.Sleep(200 * time.Millisecond)
	}
}

func renderizar(state *common.GameState) {
	fmt.Print("\033[H\033[2J")
	for y := 0; y < state.MapHeight; y++ {
		for x := 0; x < state.MapWidth; x++ {
			symbol := ' '

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

			fmt.Printf("%c", symbol)
		}
		fmt.Println()
	}
}
