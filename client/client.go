package main

import (
	"fmt"
	"fppd-jogo/common" // tipos compartilhados serv/cliente
	"log"
	"net/rpc"
	"time"
)

// registra um jogador no servidor + inicia o loop de atualização do mapa
func registerPlayer(name string) {
	client, err := rpc.Dial("tcp", "localhost:8080") // conecta ao servidor
	if err != nil {
		log.Fatal("erro ao conectar ao servidor:", err)
	}
	defer client.Close()

	req := common.JoinRequest{Name: name} // montando a req pro servidor com o nome do jogador
	var res common.JoinResponse           // variavel que recebe a resposta do servidor

	// chamada rpc para registrar o jogador
	err = client.Call("GameServer.RegisterPlayer", &req, &res)
	if err != nil {
		log.Fatal("erro na chamada rpc:", err)
	}

	fmt.Printf("jogador registrado com id: %d\n", res.ID)
	// iniciando o loop de att do jogo
	updateMap(client, res.Player.ID)
}

// loop que atualiza o mapa
func updateMap(client *rpc.Client, playerID int) {
	for {
		var state common.GameState                     // variavel que recebe o estado do jogo
		req := common.StateRequest{PlayerID: playerID} // montando a req com o id do jogador

		err := client.Call("GameServer.GetState", req, &state) // chamando o método GetState do servidor
		if err != nil {
			log.Println("erro ao obter estado do jogo:", err)
			time.Sleep(time.Second)
			continue // espera 1 segundo e tenta de novo
		}

		renderMap(&state)
		time.Sleep(500 * time.Millisecond) // espera 500ms antes de fazer a próxima atualização
	}
}

// renderiza o mapa do jogo
func renderMap(state *common.GameState) {
	fmt.Print("\033[H\033[2J") // limpa o terminal

	for y := 0; y < state.MapHeight; y++ {
		for x := 0; x < state.MapWidth; x++ {
			symbol := " " // símbolo padrão é espaço em branco
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
			fmt.Print(symbol) // imprime o símbolo na posição atual do mapa
		}
		fmt.Println() // vai para próxima depois de imprimir uma linha completa
	}
}
