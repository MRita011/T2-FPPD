package server

import (
	"T2-FPPD/shared"
)

type Servidor struct{}

func (s *Servidor) Entrar(pedido shared.PedidoEntrada, resposta *shared.RespostaEntrada) error {
	resposta.Mensagem = "Bem-vindo, " + pedido.IDJogador + "!"
	return nil
}
