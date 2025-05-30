package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"fppd-jogo/common"
)

type GameServer struct {
	state *StateGame
}

func (g *GameServer) RegisterPlayer(req *common.JoinRequest, res *common.JoinResponse) error {
	id, err := g.state.RegisterPlayer(req.Name)
	if err != nil {
		return err
	}
	res.ID = id
	return nil
}

func main() {
	server := &GameServer{state: NewStateGame()}

	err := rpc.Register(server)
	if err != nil {
		log.Fatal("Erro ao registrar servidor:", err)
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}

	fmt.Println("Servidor ouvindo na porta 8080...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Erro na conexão:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
