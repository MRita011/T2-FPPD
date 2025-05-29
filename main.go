package main

import (
	"log"
	"net"
	"net/rpc"
	"T2-FPPD/server"
)

func main() {
	s := new(server.Servidor)
	rpc.Register(s)

	ouvinte, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Erro ao escutar: ", err)
	}
	defer ouvinte.Close()

	log.Println("Servidor iniciado na porta 8080...")
	rpc.Accept(ouvinte)
}
