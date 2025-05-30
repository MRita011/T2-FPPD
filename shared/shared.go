package shared

// define a estrutura do pedido que o cliente vai enviar pro servidor
type PedidoEntrada struct {
	IDJogador string // o id ou nome do jogador que quer entrar
}

// define a estrutura da resposta que o servidor vai mandar de volta pro cliente
type RespostaEntrada struct {
	Mensagem string // a mensagem que o servidor quer enviar como resposta
}
