package main

import "github.com/nsf/termbox-go"

type Cor = termbox.Attribute

const (
	CorPadrao      Cor = termbox.ColorDefault
	CorCinzaEscuro     = termbox.ColorDarkGray
	CorVermelho        = termbox.ColorRed
	CorBranco          = termbox.ColorWhite
	CorVerde           = termbox.ColorGreen
	CorAzul            = termbox.ColorBlue
	CorAmarelo         = termbox.ColorYellow
	CorMagenta         = termbox.ColorMagenta
	CorCyan            = termbox.ColorCyan
	CorParede          = termbox.ColorBlack | termbox.AttrBold
	CorFundoParede     = termbox.ColorDarkGray
	CorTexto           = termbox.ColorDarkGray
)

type TipoCaixa string

const (
	Tesouro   TipoCaixa = "tesouro"
	Armadilha TipoCaixa = "armadilha"
	Vazia     TipoCaixa = "vazia"
)

// limitando as caixas no mapa
const (
	TOTAL_CAIXAS     = 20
	TESOUROS_COUNT   = 12
	ARMADILHAS_COUNT = 4
	VAZIAS_COUNT     = 4
)

type EstadoPartida string

const (
	EstadoAguardando EstadoPartida = "aguardando"
	EstadoJogando    EstadoPartida = "jogando"
	EstadoFinalizado EstadoPartida = "finalizado"
)

type Coordenada struct {
	X, Y int
}

type Elemento struct {
	Simbolo  rune
	Cor      Cor
	CorFundo Cor
	Tangivel bool
}

// estrutura que representa um jogador
type Jogador struct {
	ID        string // id único do jogador
	Nome      string // nome do jogador
	PosX      int    // posição x no mapa
	PosY      int    // posição y no mapa
	Cor       Cor    // cor do jogador
	Simbolo   rune   // símbolo que representa o jogador
	Conectado bool   // se está conectado ou não
	GameOver  bool   // se o jogo acabou para esse jogador
	Pontuacao int    // número de tesouros coletados
	Morto     bool   // se o jogador morreu por armadilha
}

// estrutura com o estado atual do jogo que é compartilhado com os clientes
type EstadoJogo struct {
	Mapa              [][]Elemento             // o mapa atual com todos os elementos
	Jogadores         map[string]*Jogador      // todos os jogadores conectados
	StatusMsg         string                   // mensagem de status que aparece na tela
	Caixas            map[Coordenada]TipoCaixa // mapa das caixas (tesouros e armadilhas)
	GameOver          bool                     // indica se o jogo terminou para este jogador
	EstadoPartida     EstadoPartida            // estado atual da partida
	Vencedor          string                   // ID do jogador vencedor
	TesourosRestantes int                      // número de tesouros restantes
	JogadoresVivos    int                      // número de jogadores ainda vivos
}

// Nova estrutura para armazenar apenas as posições dos jogadores
type PosicoesJogadores struct {
	Jogadores         map[string]PosicaoJogador // posições de todos os jogadores
	JogadorID         string                    // id do jogador atual
	UltimoProcessado  int64                     // último comando processado
	Caixas            map[Coordenada]TipoCaixa  // caixas no mapa
	EstadoPartida     EstadoPartida             // estado atual da partida
	Vencedor          string                    // ID do jogador vencedor, se houver
	TesourosRestantes int                       // tesouros restantes no jogo
	JogadoresVivos    int                       // jogadores ainda vivos
}

// Interação com caixas
type InteragirRequest struct {
	JogadorID string
}

// Resposta traz de volta o mapa de caixas atualizado e o tipo que o jogador acabou de revelar
type InteragirResponse struct {
	Caixas            map[Coordenada]TipoCaixa
	Tipo              TipoCaixa     // Tesouro, Armadilha ou Vazia
	GameOver          bool          // se o jogador morreu
	EstadoPartida     EstadoPartida // estado atual da partida
	Vencedor          string        // ID do vencedor, se houver
	Pontuacao         int           // pontuação atual do jogador
	TesourosRestantes int           // tesouros restantes
	JogadoresVivos    int           // jogadores vivos
}

// representar a posição de um jogador
type PosicaoJogador struct {
	ID        string // id único do jogador
	Nome      string // nome do jogador
	PosX      int    // posição x no mapa
	PosY      int    // posição y no mapa
	Cor       Cor    // cor do jogador
	Simbolo   rune   // símbolo que representa o jogador
	Conectado bool   // se está conectado ou não
	Pontuacao int    // pontuação do jogador
	Morto     bool   // se o jogador está morto
}

// estrutura que representa o jogo no servidor
type Jogo struct {
	ID             string
	Mapa           [][]Elemento
	Jogadores      map[string]*Jogador
	UltimoVisitado Elemento // guarda o último elemento que o jogador pisou
	StatusMsg      string
	GameOver       bool // indica se o jogo terminou para este jogador
}

type EventoTeclado struct {
	Tipo  string // tipo do evento (ex: mover, sair)
	Tecla rune   // tecla apertada
}

// mandar comandos do cliente pro servidor
type ComandoRequest struct {
	JogadorID      string
	SequenceNumber int64       // número da ação pra evitar fora de ordem
	Comando        string      // tipo de comando
	Dados          interface{} // dados do comando (pode ser qualquer coisa)
}

// estrutura usada quando o jogador quer se mover
type MoverRequest struct {
	JogadorID      string
	SequenceNumber int64
	Tecla          rune
}

// estrutura usada quando o jogador se conecta
type ConectarRequest struct {
	MapaFile string // arquivo do mapa que o cliente quer usar
	Nome     string // nome do jogador que está se conectando
}

// resposta do servidor quando o jogador se conecta
type ConectarResponse struct {
	JogadorID string     // id que o servidor gerou pro jogador
	Estado    EstadoJogo // estado atual do jogo que o cliente vai receber
}

// Nova estrutura para resposta do servidor com apenas as posições
type ConectarPosicaoResponse struct {
	JogadorID string            // id que o servidor gerou pro jogador
	Posicoes  PosicoesJogadores // posições dos jogadores
}

var (
	// Caixa      = Elemento{'■', CorPadrao, CorPadrao, false} // caixa genérica (tesouro ou armadilha)
	Personagem = Elemento{'♟', CorBranco, CorPadrao, true}      // jogador
	Inimigo    = Elemento{'♙', CorVermelho, CorPadrao, true}    // inimigo
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true} // parede
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}      // vegetação (não colide)
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}     // espaço vazio
)

// cores que os jogadores podem ter
var CoresJogadores = []Cor{
	CorBranco, CorVermelho, CorVerde, CorAzul,
	CorAmarelo, CorMagenta, CorCyan,
}
