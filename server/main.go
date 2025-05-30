package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"fppd-jogo/common" // pacote comum com structs usadas no rpc
)

// servidor do jogo, que guarda o estado da partida
type GameServer struct {
	state *StateGame // estado atual do jogo (onde estão os jogadores, monstros, etc)
}

// método que registra o jogador, recebe uma requisição e responde com o id do jogador
func (g *GameServer) RegisterPlayer(req *common.JoinRequest, res *common.JoinResponse) error {
	// tenta registrar o jogador no estado do jogo
	id, err := g.state.RegisterPlayer(req.Name)
	if err != nil {
		// se deu ruim, retorna o erro pra rpc
		return err
	}
	// coloca o id retornado na resposta pra enviar pro cliente
	res.ID = id
	return nil
}

func main() {
	// cria o servidor do jogo com o estado inicial zerado
	server := &GameServer{state: NewStateGame()}

	// registra o servidor pra rpc conseguir expor os métodos
	err := rpc.Register(server)
	if err != nil {
		log.Fatal("Erro ao registrar servidor:", err)
	}

	// cria um listener na porta 8080 pra aceitar conexões tcp
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}

	fmt.Println("Servidor ouvindo na porta 8080...")
	for {
		// fica esperando uma nova conexão chegar
		conn, err := listener.Accept()
		if err != nil {
			// se der erro na conexão, só loga e continua esperando outras conexões
			log.Println("Erro na conexão:", err)
			continue
		}
		// para cada conexão aceita, cria uma goroutine que vai servir essa conexão com rpc
		go rpc.ServeConn(conn)
	}
}
