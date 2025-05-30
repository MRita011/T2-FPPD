package server

import (
	"T2-FPPD/shared"
)

// define a struct vazia que vai servir como base pro nosso servidor rpc
type Servidor struct{}

// método chamado remotamente pelo cliente quando ele quer \"entrar\" no servidor
func (s *Servidor) Entrar(pedido shared.PedidoEntrada, resposta *shared.RespostaEntrada) error {
	// monta a mensagem de resposta, dando boas-vindas pro jogador com o id que ele mandou
	resposta.Mensagem = "Bem-vindo, " + pedido.IDJogador + "!"
	// retorna nil porque deu tudo certo (nenhum erro)
	return nil
}
