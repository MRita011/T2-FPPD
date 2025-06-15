package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	// define uma flag booleano "server"
	// se passar "--server" na linha de comando, 'servidor' = true
	servidor := flag.Bool("server", false, "Servidor")
	flag.Parse() // processa todos as flags

	// se a flag 'servidor' = true, roda como servidor;
	// caso contrário, como cliente
	if *servidor {
		runServidor()
	} else {
		runCliente()
	}
}

// inicia o servidor de posições dos jogadores
func runServidor() {
	// cria uma nova instância do GameServe
	server := NewGameServer()
	log.Println("Servidor de posições iniciado na porta 8080")
	// inicia o servidor RPC na porta 8080
	// log.Fatal encerra o programa se der erro
	log.Fatal(server.StartRPC("8080"))
}

// inicia o cliente
func runCliente() {
	// inicializa a interface
	IniciarInterface()
	defer FinalizarInterface() // fecha a interface ao sair

	log.Println("Iniciando cliente...")
	// tenta criar um novo cliente
	client, err := NewGameClient()
	if err != nil {
		log.Fatal("Erro ao conectar:", err) // se não conectar, fecha
	}
	defer client.Close() // fecha o cliente ao sair

	// conecta o jogador ao jogo no servidor, carregando um mapa
	log.Println("Conectando ao jogo...")
	jogadorID, err := client.ConectarJogo("mapa.txt")
	if err != nil {
		log.Fatal("Erro ao conectar ao jogo:", err) // se der erro, fecha
	}
	log.Println("Conectado com sucesso! ID:", jogadorID)

	// inicia a sincronizacao de estado com o servidor em uma goroutine separada
	client.IniciarSincronizacao(jogadorID)
	defer client.PararSincronizacao()
	var tempoGameOver time.Duration

	// loop principal do jogo no cliente.
	for {
		evento := LerEvento() // lê eventos do teclado

		if evento.Tipo == "sair" {
			break // sai do loop
		}
		if evento.Tipo == "mover" {
			// se o evento for um movimento, manda pro servidor
			err := client.Mover(jogadorID, evento.Tecla)
			if err != nil {
				log.Printf("Erro ao mover: %v", err)
			}
		}
		if evento.Tipo == "interagir" {
			err := client.Interagir(jogadorID)
			if err != nil {
				log.Printf("Erro ao interagir: %v", err)
			}
			// Força uma atualização imediata da tela após interagir
			estado, err := client.ObterEstado()
			if err == nil {
				DesenharEstadoJogo(estado)
			}
		}
		estado := client.gameManager.ObterEstado()
		if estado.GameOver {
			// Continua desenhando, mas não permite mais movimentos
			DesenharEstadoJogo(estado)
			time.Sleep(time.Second) // Pequena pausa para mostrar o GAMEOVER

			// Incrementa o contador de tempo
			tempoGameOver += time.Second

			// Após 5 segundos, encerra o jogo
			if tempoGameOver >= 5*time.Second {
				break // Sai do loop do jogo
			}
		} else {
			// Reseta o contador se o jogo não estiver em game over
			tempoGameOver = 0
		}
	}
}
