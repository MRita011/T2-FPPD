package main

import (
	"net/rpc"
	"sync"
	"time"
)

// essa struct é o cliente do jogo, que se conecta no servidor rpc
type GameClient struct {
	client         *rpc.Client   // conexão com o servidor via rpc
	config         NetworkConfig // configurações de rede (ex: ip, porta)
	sequenceNumber int64         // número de sequência pra garantir ordem dos comandos
	estadoAtual    *EstadoJogo   // estado atual do jogo que o cliente tem
	jogadorID      string        // id único do jogador nesse cliente
	mutex          sync.RWMutex  // trava de leitura/escrita pra acessar dados com segurança
	sincronizando  bool          // flag que diz se a sincronização tá rolando
	stopSync       chan bool     // canal pra mandar sinal de parar a sincronização
}

// cria um novo cliente com a config padrão
func NewGameClient() (*GameClient, error) {
	return NewGameClientWithConfig(DefaultConfig)
}

// cria um novo cliente com uma config específica
func NewGameClientWithConfig(config NetworkConfig) (*GameClient, error) {
	client, err := rpc.Dial("tcp", config.GetAddress()) // tenta conectar no servidor
	if err != nil {
		return nil, err
	}
	return &GameClient{
		client:   client,
		config:   config,
		stopSync: make(chan bool),
	}, nil
}

// fecha o cliente (desconecta e para a sync)
func (gc *GameClient) Close() error {
	gc.PararSincronizacao()
	return gc.client.Close()
}

// pega o nome do arquivo de mapa padrão que veio da config
func (gc *GameClient) GetDefaultMapFile() string {
	return gc.config.DefaultMapFile
}

// conecta o jogador no jogo, usando um arquivo de mapa
func (gc *GameClient) ConectarJogo(mapaFile string) (string, error) {
	req := ConectarRequest{MapaFile: mapaFile}
	var resp ConectarResponse
	err := gc.client.Call("GameService.ConectarJogo", req, &resp) // chama o servidor
	if err != nil {
		return "", err
	}

	// guarda o id do jogador e o estado do jogo retornado
	gc.mutex.Lock()
	gc.jogadorID = resp.JogadorID
	gc.estadoAtual = &resp.Estado
	gc.mutex.Unlock()

	return resp.JogadorID, nil
}

// envia um movimento pro servidor (tipo pressionar uma tecla)
func (gc *GameClient) Mover(jogadorID string, tecla rune) error {
	gc.sequenceNumber++ // incrementa o número de sequência
	req := MoverRequest{
		JogadorID:      jogadorID,
		SequenceNumber: gc.sequenceNumber,
		Tecla:          tecla,
	}

	var estado EstadoJogo
	err := gc.client.Call("GameService.Mover", req, &estado) // manda pro servidor
	if err != nil {
		return err
	}

	// atualiza o estado local com o que o servidor retornou
	gc.mutex.Lock()
	gc.estadoAtual = &estado
	gc.mutex.Unlock()

	return nil
}

// pega o estado atual do jogo direto do servidor
func (gc *GameClient) ObterEstado() (*EstadoJogo, error) {
	var estado EstadoJogo
	err := gc.client.Call("GameService.ObterEstado", struct{}{}, &estado)
	if err != nil {
		return nil, err
	}

	// salva esse estado no cliente
	gc.mutex.Lock()
	gc.estadoAtual = &estado
	gc.mutex.Unlock()

	return &estado, nil
}

// começa a sincronização automática com o servidor
func (gc *GameClient) IniciarSincronizacao(jogadorID string) {
	gc.mutex.Lock()
	if gc.sincronizando { // se já tá sincronizando, não faz nada
		gc.mutex.Unlock()
		return
	}
	gc.sincronizando = true
	gc.mutex.Unlock()

	// inicia a goroutine que sincroniza o estado
	go gc.loopSincronizacao()
}

// para a sincronização automática
func (gc *GameClient) PararSincronizacao() {
	gc.mutex.Lock()
	if !gc.sincronizando { // se já tava parado, só sai
		gc.mutex.Unlock()
		return
	}
	gc.sincronizando = false
	gc.mutex.Unlock()

	// manda sinal pra parar a goroutine
	select {
	case gc.stopSync <- true:
	default:
	}
}

// essa função roda numa goroutine e fica pedindo o estado do servidor o tempo todo
func (gc *GameClient) loopSincronizacao() {
	ticker := time.NewTicker(100 * time.Millisecond) // faz isso 10 vezes por segundo
	defer ticker.Stop()

	for {
		select {
		case <-gc.stopSync:
			// recebeu sinal pra parar, então sai
			return
		case <-ticker.C:
			// hora de sincronizar: pega o estado novo do servidor
			estado, err := gc.ObterEstado()
			if err != nil {
				continue // se der erro, só ignora e tenta dnv na próxima
			}

			// checa se o jogador ainda tá presente no estado novo
			gc.mutex.RLock()
			jogadorAtual := estado.Jogadores[gc.jogadorID]
			gc.mutex.RUnlock()

			if jogadorAtual != nil {
				// se o jogador ainda tá no jogo, desenha o estado novo na tela
				DesenharEstadoJogo(estado)
			}
		}
	}
}
