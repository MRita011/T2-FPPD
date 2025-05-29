package main

import (
	"T2-FPPD/shared"
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	cliente, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Erro ao conectar: ", err)
	}

	pedido := shared.PedidoEntrada{IDJogador: "Deus Tenha Piedade"}
	var resposta shared.RespostaEntrada

	err = cliente.Call("Servidor.Entrar", pedido, &resposta)
	if err != nil {
		log.Fatal("Erro ao chamar método: ", err)
	}
	fmt.Println("Resposta do servidor:", resposta.Mensagem)
}
