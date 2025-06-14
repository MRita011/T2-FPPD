package main

import (
	"flag"
	"log"
)

func main() {
	// Verifica se foi passado o argumento "--server"
	servidor := flag.Bool("server", false, "Servidor")
	flag.Parse() // processa os argumentos da linha de comando

	if *servidor {
		runServidor() // se for servidor, roda o servidor
	} else {
		runCliente() // se não, roda o cliente
	}
}

// Iniciar o servidor de posições dos jogadores
func runServidor() {
	server := NewGameServer()
	log.Println("Servidor de posições iniciado na porta 8080")
	log.Fatal(server.StartRPC("8080")) // inicia o servidor e encerra se der erro
}

// Iniciar cliente
func runCliente() {

	log.Println("Iniciando cliente...")
	client, err := NewGameClient() // tenta criar um novo cliente
	if err != nil {
		log.Fatal("Erro ao conectar:", err)
	}
	defer client.Close()

	// Conecta ao servidor e carrega o jogo local
	log.Println("Conectando ao jogo...")
	jogadorID, err := client.ConectarJogo("mapa.txt")
	if err != nil {
		log.Fatal("Erro ao conectar ao jogo:", err) // se não conseguir conectar, finaliza
	}
	log.Println("Conectado com sucesso! ID:", jogadorID)

	IniciarInterface()
	defer FinalizarInterface()

	// Começa a sincronizar estado com o servidor
	client.IniciarSincronizacao(jogadorID)
	defer client.PararSincronizacao()

	// Loop principal do jogo
	for {
		evento := LerEvento() // lê o que o jogador apertou

		if evento.Tipo == "sair" {
			break // se apertou esc, sai do jogo
		}
		if evento.Tipo == "mover" {
			client.Mover(jogadorID, evento.Tecla) // envia o movimento pro servidor
		}
		if evento.Tipo == "interagir" {
			client.Interagir(jogadorID)
		}
	}
}
