package main

import (
	"fmt"
	"fppd-jogo/common" // pacote com as strucs/tipos comuns entre servidor e cliente
	"log"
	"net/rpc"
)

// registra um jogador no servidor
func registerPlayer(name string) {
	client, err := rpc.Dial("tcp", "localhost:8080") // conectando ao servidor RPC que roda na porta 8080

	if err != nil { // se der erro na conexão, encerra o programa
		log.Fatal("erro ao conectar ao servidor:", err)
	}
	defer client.Close() // fecha a conexão quando terminar

	req := common.JoinRequest{Name: name} // montando a req pro servidor com o nome do jogador

	var res common.JoinResponse // variavel que recebe a resposta do servidor

	// chamando o método RegisterPlayer do servidor, passando a requisição e esperando a resposta
	err = client.Call("GameServer.RegisterPlayer", &req, &res)
	if err != nil {
		log.Fatal("erro na chamada rpc:", err) // se der ruim na chamada, encerra o programa
	}

	// se deu tudo certo, imprime o id que o servidor deu pro jogador registrado
	fmt.Printf("jogador registrado com id: %d\n", res.ID)
}
