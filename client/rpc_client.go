package main

import (
	"fmt"
	"log"
	"net/rpc"

	"fppd-jogo/common"
)

func registerPlayer(name string) {
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Erro ao conectar ao servidor: ", err)
	}
	defer client.Close()

	req := common.JoinRequest{Name: name}
	var res common.JoinResponse

	err = client.Call("GameServer.RegisterPlayer", &req, &res)
	if err != nil {
		log.Fatal("Erro na chamada RPC: ", err)
	}

	fmt.Printf("Jogador registrado com ID: %d\n", res.ID)
}
