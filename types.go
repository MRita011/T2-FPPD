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
}

// estrutura com o estado atual do jogo que é compartilhado com os clientes
type EstadoJogo struct {
	Mapa      [][]Elemento             // o mapa atual com todos os elementos
	Jogadores map[string]*Jogador      // todos os jogadores conectados
	StatusMsg string                   // mensagem de status que aparece na tela
	Caixas    map[Coordenada]TipoCaixa // mapa das caixas (tesouros e armadilhas)
}

// Nova estrutura para armazenar apenas as posições dos jogadores
type PosicoesJogadores struct {
	Jogadores        map[string]PosicaoJogador // posições de todos os jogadores
	JogadorID        string                    // id do jogador atual
	UltimoProcessado int64                     // último comando processado
	Caixas           map[Coordenada]TipoCaixa  // caixas no mapa
}

// Interação com caixas
type InteragirRequest struct {
	JogadorID string
}

// Resposta traz de volta o mapa de caixas atualizado e o tipo que o jogador acabou de revelar
type InteragirResponse struct {
	Caixas map[Coordenada]TipoCaixa
	Tipo   TipoCaixa // Tesouro ou Armadilha
}

// Estrutura minimalista para representar a posição de um jogador
type PosicaoJogador struct {
	ID        string // id único do jogador
	Nome      string // nome do jogador
	PosX      int    // posição x no mapa
	PosY      int    // posição y no mapa
	Cor       Cor    // cor do jogador
	Simbolo   rune   // símbolo que representa o jogador
	Conectado bool   // se está conectado ou não
}

// estrutura que representa o jogo no servidor
type Jogo struct {
	ID             string
	Mapa           [][]Elemento
	Jogadores      map[string]*Jogador
	UltimoVisitado Elemento // guarda o último elemento que o jogador pisou
	StatusMsg      string
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
	Personagem = Elemento{'☺', CorBranco, CorPadrao, true}      // jogador
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}    // inimigo
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true} // parede
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}      // vegetação (não colide)
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}     // espaço vazio
)

// cores que os jogadores podem ter
var CoresJogadores = []Cor{
	CorBranco, CorVermelho, CorVerde, CorAzul,
	CorAmarelo, CorMagenta, CorCyan,
}
