package main

import (
	"log"
	"net"
	"net/rpc"
)

type GameServer struct {
	manager *GameManager // gerenciador do jogo
}

// estrutura do rpc que o cliente var chamar
type GameService struct {
	manager *GameManager
}

// criando um novo servidor
func NewGameServer() *GameServer {
	return &GameServer{manager: NewGameManager()}
}

// iniciando o servidor
func (gs *GameServer) StartRPC(port string) error {
	service := &GameService{manager: gs.manager} // cria o serviço com o gerenciador do jogo
	rpc.Register(service)                        // registra esse serviço pro rpc

	// fica escutando na porta do servidor
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	log.Printf("servidor rpc iniciado na porta %s", port)

	// loop para aceitar conexões dos clientes
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("erro ao aceitar conexão: %v", err)
			continue // se der erro numa conexão, ignora e continua esperando outras
		}
		go rpc.ServeConn(conn) // atende a conexão em paralelo
	}
}

// rpc: jogador se conecta ao jogo
func (gs *GameService) ConectarJogo(req ConectarRequest, reply *ConectarResponse) error {
	jogador, estado, err := gs.manager.ConectarJogador(req.MapaFile) // tenta conectar o jogador
	if err != nil {
		return err
	}

	// se deu certo, preenche a resposta com o id do jogador e o estado atual do jogo
	reply.JogadorID = jogador.ID
	reply.Estado = *estado
	return nil
}

// rpc: o jogador tenta se mover
func (gs *GameService) Mover(req MoverRequest, reply *EstadoJogo) error {
	estado, err := gs.manager.MoverJogador(req.JogadorID, req.Tecla, req.SequenceNumber) // tenta mover o jogador
	if err != nil {
		return err
	}

	*reply = *estado // se funcionar, retorna o novo estado do jogo
	return nil
}

// pegando o estado atual do jogo (usado na sincronização)
func (gs *GameService) ObterEstado(args struct{}, reply *EstadoJogo) error {
	estado := gs.manager.ObterEstado() // pega o estado do jogo
	*reply = *estado                   // coloca o estado na resposta
	return nil
}
