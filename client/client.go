package main

import (
	"T2-FPPD/shared"
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	// tenta conectar com o servidor rpc que tá rodando no localhost na porta 8080
	cliente, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		// se der erro na conexão, mostra no log e encerra o programa
		log.Fatal("erro ao conectar: ", err)
	}

	// cria um pedido com o id do jogador que vai tentar entrar no servidor
	pedido := shared.PedidoEntrada{IDJogador: "Deus Tenha Piedade"}
	// cria uma variável pra guardar a resposta que o servidor vai mandar de volta
	var resposta shared.RespostaEntrada

	// chama o método remoto 'Entrar' do servidor, passando o pedido e esperando a resposta
	err = cliente.Call("Servidor.Entrar", pedido, &resposta)
	if err != nil {
		// se tiver erro na chamada do método rpc, mostra o erro e para tudo
		log.Fatal("erro ao chamar método: ", err)
	}
	// imprime no terminal a mensagem que veio do servidor
	fmt.Println("resposta do servidor:", resposta.Mensagem)
}
