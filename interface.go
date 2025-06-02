package main

import "github.com/nsf/termbox-go"

// inicia a interface do terminal
func IniciarInterface() {
	if err := termbox.Init(); err != nil {
		panic(err) // se der erro, o programa para
	}
}

// finaliza e limpa a interface do terminal
func FinalizarInterface() {
	termbox.Close()
}

// lê uma tecla pressionada e retorna um evento traduzido
func LerEvento() EventoTeclado {
	ev := termbox.PollEvent() // espera uma tecla ser pressionada

	if ev.Type != termbox.EventKey {
		return EventoTeclado{} // se não for tecla, retorna vazio
	}

	// tecla ESC para sair
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}

	// tecla 'e' para interagir
	if ev.Ch == 'e' {
		return EventoTeclado{Tipo: "interagir"}
	}

	// qualquer outra tecla é movimento
	return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

// desenha o estado do jogo (mapa + jogadores + status)
func DesenharEstadoJogo(estado *EstadoJogo) {
	termbox.Clear(CorPadrao, CorPadrao) // limpa a tela

	// desenha o mapa
	if estado.Mapa != nil {
		for y, linha := range estado.Mapa {
			for x, elem := range linha {
				termbox.SetCell(x, y, elem.Simbolo, elem.Cor, elem.CorFundo)
			}
		}
	}

	// desenha todos os jogadores conectados
	if estado.Jogadores != nil {
		for _, jogador := range estado.Jogadores {
			if jogador.Conectado {
				termbox.SetCell(jogador.PosX, jogador.PosY, jogador.Simbolo, jogador.Cor, CorPadrao)
			}
		}
	}

	// mostra a mensagem de status
	statusY := 0
	if estado.Mapa != nil {
		statusY = len(estado.Mapa) + 1
	}
	for i, c := range estado.StatusMsg {
		termbox.SetCell(i, statusY, c, CorTexto, CorPadrao)
	}

	// lista de jogadores conectados
	if estado.Jogadores != nil {
		infoY := statusY + 2
		info := "Jogadores conectados:"
		for i, c := range info {
			termbox.SetCell(i, infoY, c, CorTexto, CorPadrao)
		}

		linha := infoY + 1
		for _, jogador := range estado.Jogadores {
			if jogador.Conectado {
				texto := jogador.Nome + " " + string(jogador.Simbolo)
				for i, c := range texto {
					termbox.SetCell(i, linha, c, jogador.Cor, CorPadrao)
				}
				linha++
			}
		}
	}

	// mostra instruções ao jogador
	instrY := statusY + 2
	if estado.Jogadores != nil {
		instrY = statusY + 2 + len(estado.Jogadores) + 2
	}
	msg := "Use WASD para mover. ESC para sair."
	for i, c := range msg {
		termbox.SetCell(i, instrY, c, CorTexto, CorPadrao)
	}

	termbox.Flush() // atualiza a tela
}
