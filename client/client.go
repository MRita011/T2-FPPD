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
		log.Fatal("erro ao conectar ao servidor:", err)
	}
	defer client.Close()

	req := common.JoinRequest{Name: name}
	var res common.JoinResponse

	err = client.Call("GameServer.RegisterPlayer", &req, &res)
	if err != nil {
		log.Fatal("erro na chamada rpc:", err)
	}

	fmt.Printf("jogador registrado com id: %d\n", res.ID)

	// Goroutine para escutar o teclado
	go listenInput(client, res.Player.ID)

	// Loop para atualização do mapa
	updateMap(client, res.Player.ID)
}

func listenInput(client *rpc.Client, playerID int) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Digite comando (ex: w, a, s, d): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Aqui você pode enviar comandos ao servidor se desejar
		fmt.Println("Comando recebido:", input)
	}
}

func updateMap(client *rpc.Client, playerID int) {
	for {
		var state common.GameState
		req := common.StateRequest{PlayerID: playerID}

		err := client.Call("GameServer.GetState", req, &state)
		if err != nil {
			log.Println("erro ao obter estado do jogo:", err)
			time.Sleep(time.Second)
			continue
		}

		renderMap(&state)
		time.Sleep(500 * time.Millisecond)
	}
}

func renderMap(state *common.GameState) {
	fmt.Print("\033[H\033[2J")

	for y := 0; y < state.MapHeight; y++ {
		for x := 0; x < state.MapWidth; x++ {
			symbol := " "

			for _, p := range state.Players {
				if p.X == x && p.Y == y {
					symbol = p.Symbol
					break
				}
			}

			for _, t := range state.Traps {
				if t.X == x && t.Y == y {
					symbol = t.Symbol
					break
				}
			}

			for _, t := range state.Treasures {
				if t.X == x && t.Y == y {
					symbol = t.Symbol
					break
				}
			}

			fmt.Print(symbol)
		}
		fmt.Println()
	}
}
