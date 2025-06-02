package main

import (
	"flag"
	"log"
)

func main() {
	// verifica se foi passado o argumento "--server"
	servidor := flag.Bool("server", false, "Servidor")
	flag.Parse() // processa os argumentos da linha de comando

	if *servidor {
		runServidor() // se for servidor, roda o servidor
	} else {
		runCliente() // se não, roda o cliente
	}
}

// iniciar o servidor
func runServidor() {
	server := NewGameServer()
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(server.StartRPC("8080")) // inicia o servidor rpc e encerra se der erro
}

// iniciar cliente
func runCliente() {
	IniciarInterface()
	defer FinalizarInterface()

	client, err := NewGameClient() // tenta criar um novo cliente
	if err != nil {
		log.Fatal("Erro ao conectar:", err)
	}
	defer client.Close()

	// conecta o jogador no servidor e pega o id dele
	jogadorID, err := client.ConectarJogo("mapa.txt")
	if err != nil {
		log.Fatal("Erro ao conectar ao jogo:", err) // se não conseguir conectar, finaliza
	}

	client.IniciarSincronizacao(jogadorID) // começa a sincronizar o estado do jogo com o servidor
	defer client.PararSincronizacao()

	// loop principal do jogo
	for {
		evento := LerEvento() // lê o que o jogador apertou

		if evento.Tipo == "sair" {
			break // se apertou esc, sai do jogo
		}
		if evento.Tipo == "mover" {
			client.Mover(jogadorID, evento.Tecla) // se apertou wasd, envia o movimento pro servidor
		}
	}
}
