package main

import (
	"log"
	"net"
	"net/rpc"
	"T2-FPPD/server"
)

func main() {
	// cria uma instância do nosso servidor rpc
	s := new(server.Servidor)

	// registra esse servidor pra ficar disponível nas chamadas rpc
	rpc.Register(s) 

	// tenta abrir uma escuta na porta 8080 via tcp
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		// se der ruim pra escutar, imprime o erro e termina o programa
		log.Fatal("Erro ao escutar: ", err)
	}
	// fecha o listener antes de sair do main
	defer listener.Close()

	// mostra no console que o servidor subiu certinho na porta 8080
	log.Println("Servidor iniciado na porta 8080...")
	rpc.Accept(listener) // fica aceitando conexões rpc nesse listener até o programa ser interrompido
}
