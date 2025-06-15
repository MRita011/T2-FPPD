package main

import (
	"log"
	"net/rpc"
	"sync"
	"time"
)

// GameClient se conecta ao servidor e sincroniza as posições dos jogadores
type GameClient struct {
	client         *rpc.Client   // conexão com o servidor via rpc
	config         NetworkConfig // configurações de rede (ex: ip, porta)
	sequenceNumber int64         // número de sequência pra garantir ordem dos comandos
	gameManager    *GameManager  // gerenciador do jogo local
	jogadorID      string        // id único do jogador nesse cliente
	mutex          sync.RWMutex  // trava de leitura/escrita pra acessar dados com segurança
	sincronizando  bool          // flag que diz se a sincronização tá rolando
	stopSync       chan bool     // canal pra mandar sinal de parar a sincronização
}

// Cria um novo cliente com a config padrão
func NewGameClient() (*GameClient, error) {
	return NewGameClientWithConfig(DefaultConfig)
}

// Cria um novo cliente com uma config específica
func NewGameClientWithConfig(config NetworkConfig) (*GameClient, error) {
	log.Println("Conectando cliente ao servidor em", config.GetAddress())
	client, err := rpc.Dial("tcp", config.GetAddress()) // tenta conectar no servidor
	if err != nil {
		return nil, err
	}

	return &GameClient{
		client:      client,
		config:      config,
		gameManager: NewGameManager(),
		stopSync:    make(chan bool),
	}, nil
}

// Fecha o cliente (desconecta e para a sync)
func (gc *GameClient) Close() error {
	gc.PararSincronizacao()

	// Tenta desconectar o jogador do servidor
	if gc.jogadorID != "" {
		var resposta bool
		gc.client.Call("GameService.Desconectar", gc.jogadorID, &resposta)
	}

	return gc.client.Close()
}

// Retorna o gerenciador de jogo local
func (gc *GameClient) GetGameManager() *GameManager {
	return gc.gameManager
}

// Pega o nome do arquivo de mapa padrão que veio da config
func (gc *GameClient) GetDefaultMapFile() string {
	return gc.config.DefaultMapFile
}

// Conecta o jogador no jogo
func (gc *GameClient) ConectarJogo(mapaFile string) (string, error) {
	// Inicializa o jogo local com o mapa
	if err := gc.gameManager.InicializarJogo(mapaFile); err != nil {
		return "", err
	}

	// Prepara a requisição para o servidor
	req := ConectarRequest{
		MapaFile: mapaFile,
		Nome:     "Jogador" + time.Now().Format("15:04:05"),
	}

	// Chama o servidor para conectar
	var resp ConectarPosicaoResponse
	err := gc.client.Call("GameService.ConectarJogo", req, &resp)
	if err != nil {
		return "", err
	}

	// Guarda o ID do jogador
	gc.mutex.Lock()
	gc.jogadorID = resp.JogadorID
	gc.mutex.Unlock()

	// Encontra os dados do jogador local nas posições recebidas
	jogadorLocal := resp.Posicoes.Jogadores[resp.JogadorID]

	// Cria o jogador local no gerenciador de jogo
	gc.gameManager.CriarJogadorLocal(
		resp.JogadorID,
		jogadorLocal.Nome,
		jogadorLocal.PosX,
		jogadorLocal.PosY,
		jogadorLocal.Cor,
	)

	// Atualiza as posições dos outros jogadores
	gc.gameManager.AtualizarJogadoresRemotos(resp.Posicoes.Jogadores)
	gc.gameManager.AtualizarCaixas(resp.Posicoes.Caixas)

	return resp.JogadorID, nil
}

// Envia um movimento pro servidor e atualiza o estado local
func (gc *GameClient) Mover(jogadorID string, tecla rune) error {
	// Incrementa o número de sequência
	gc.sequenceNumber++

	// Atualiza o movimento localmente primeiro
	gc.gameManager.MoverJogadorLocal(tecla)

	// Prepara a requisição para o servidor
	req := MoverRequest{
		JogadorID:      jogadorID,
		SequenceNumber: gc.sequenceNumber,
		Tecla:          tecla,
	}

	// Envia o movimento para o servidor
	var posicoes PosicoesJogadores
	err := gc.client.Call("GameService.Mover", req, &posicoes)
	if err != nil {
		return err
	}

	// Atualiza o estado local com as posições recebidas do servidor
	gc.gameManager.AtualizarJogadoresRemotos(posicoes.Jogadores)

	// Atualiza o mapa de caixas recebido do servidor
	gc.gameManager.AtualizarCaixas(posicoes.Caixas)

	return nil
}

// Interagir com caixa
func (gc *GameClient) Interagir(jogadorID string) error {
	req := InteragirRequest{JogadorID: jogadorID}
	var resp InteragirResponse
	if err := gc.client.Call("GameService.InteragirCaixa", req, &resp); err != nil {
		return err
	}

	// atualiza o mapa de caixas
	gc.gameManager.AtualizarCaixas(resp.Caixas)

	// verifica o tipo de caixa encontrada e att a mensagem
	switch resp.Tipo {
	case Tesouro:
		gc.gameManager.jogo.UltimoVisitado = Elemento{
			Simbolo:   '■',
			Cor:       CorVerde,
		}
		gc.gameManager.jogo.StatusMsg = "Você encontrou um TESOURO!"
	case Armadilha:
		gc.gameManager.mutex.Lock()
		gc.gameManager.jogo.UltimoVisitado = Elemento{
			Simbolo:   '■',
			Cor:       CorVermelho,
		}
        gc.gameManager.jogo.GameOver = true
        gc.gameManager.jogo.StatusMsg = "Você ativou uma ARMADILHA! GAME OVER!"
        gc.gameManager.mutex.Unlock()
		//Perguntar se o jogador que esperar uma nova partida começar ou sair
		gc.gameManager.jogo.StatusMsg += " Pressione 'ESC' para sair ou aguarde uma nova partida."
	case "":
		gc.gameManager.jogo.StatusMsg = "Nada a interagir aqui."
	default:
		gc.gameManager.jogo.StatusMsg = "Nem eu sei o que é isso!"
	}
	return nil
}

// Obtém as posições atualizadas do servidor
func (gc *GameClient) ObterPosicoes() error {
	var posicoes PosicoesJogadores
	err := gc.client.Call("GameService.ObterPosicoes", gc.jogadorID, &posicoes)
	if err != nil {
		return err
	}

	// Atualiza o jogo local com as posições mais recentes
	gc.gameManager.AtualizarJogadoresRemotos(posicoes.Jogadores)

	// Atualiza o mapa de caixas com as informações mais recentes
	gc.gameManager.AtualizarCaixas(posicoes.Caixas)

	return nil
}

// Obtém o estado atual do jogo local
func (gc *GameClient) ObterEstado() (*EstadoJogo, error) {
	// Tenta atualizar com as posições mais recentes do servidor
	gc.ObterPosicoes()

	// Retorna o estado atual do jogo local
	return gc.gameManager.ObterEstado(), nil
}

// Começa a sincronização automática com o servidor
func (gc *GameClient) IniciarSincronizacao(jogadorID string) {
	gc.mutex.Lock()
	if gc.sincronizando { // se já tá sincronizando, não faz nada
		gc.mutex.Unlock()
		return
	}
	gc.sincronizando = true
	gc.mutex.Unlock()

	// Inicia a goroutine que sincroniza o estado
	go gc.loopSincronizacao()
}

// Para a sincronização automática
func (gc *GameClient) PararSincronizacao() {
	gc.mutex.Lock()
	if !gc.sincronizando { // se já tava parado, só sai
		gc.mutex.Unlock()
		return
	}
	gc.sincronizando = false
	gc.mutex.Unlock()

	// Manda sinal pra parar a goroutine
	select {
	case gc.stopSync <- true:
	default:
	}
}

// Loop de sincronização que atualiza as posições periodicamente
func (gc *GameClient) loopSincronizacao() {
	ticker := time.NewTicker(100 * time.Millisecond) // faz isso 10 vezes por segundo
	defer ticker.Stop()

	for {
		select {
		case <-gc.stopSync:
			// Recebeu sinal pra parar, então sai
			return
		case <-ticker.C:
			// Hora de sincronizar: pega posições do servidor
			estado, err := gc.ObterEstado()
			if err != nil {
				continue // se der erro, só ignora e tenta dnv na próxima
			}

			// Verifica se o jogador ainda está no jogo
			gc.mutex.RLock()
			jogadorAtual := estado.Jogadores[gc.jogadorID]
			gc.mutex.RUnlock()

			if jogadorAtual != nil {
				// Desenha o estado na tela
				DesenharEstadoJogo(estado)
			}
		}
	}
}
