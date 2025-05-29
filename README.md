# Jogo Concorrente e Multiplayer em Go

## 🧭 Visão Geral

Este projeto evolui um jogo de aventura em Go com interface textual, inicialmente focado em concorrência, para uma versão multiplayer usando RPC. A nova implementação permite que **múltiplos jogadores explorem um mapa dinâmico**, enfrentando **desafios concorrentes**, coletando tesouros e interagindo com **elementos móveis**, enquanto compartilham suas ações via servidor central.

O código-base original foi fornecido pelo professor da disciplina de Fundamentos de Processamento Paralelo e Distribuído.
Repositório original:
`https://github.com/mvneves/fppd-jogo`

## 🎯 Objetivos

* Adicionar elementos concorrentes que interajam com o jogador de forma autônoma.
* Expandir o jogo para o **modo multiplayer** com **comunicação RPC cliente-servidor**.
* Implementar **níveis progressivos de jogo**, desafios, penalidades e critérios de vitória.
* Aplicar conceitos de **concorrência segura** (goroutines, canais, mutexes, select).
* Garantir consistência via **execução única de comandos RPC (exactly-once)**.

## 🕹️ Como Jogar

* Use `W`, `A`, `S`, `D` para mover o personagem.
* Pressione `E` para interagir com o ambiente.
* Pressione `ESC` para sair do jogo.
* Cada jogador é representado por um símbolo exclusivo.
* O jogo é executado em múltiplas instâncias clientes que se conectam a um **servidor central**.

## 🌍 Mapa e Níveis

O mapa é dividido em **quatro níveis**, cada um com 40 tesouros escondidos, totalizando 160. Para avançar:

| Condição                          | Resultado                             |
| --------------------------------- | ------------------------------------- |
| ≥ 20 tesouros                     | Avança normalmente                    |
| 15–19 tesouros + enfrenta monstro | Avança, mas perde 5 tesouros          |
| < 15 tesouros + enfrenta monstro  | Avança, mas perde metade dos tesouros |
| Nenhuma das condições             | Permanece no nível                    |

## ⚙️ Elementos Concorrentes

### 💰 Tesouros

* Coletáveis e acumuláveis
* Protegidos por `mutex` para evitar condições de corrida.
* Contabilizados por jogador.

### 💣 Armadilhas

  * **Espinhos Andarilhos**: móveis, causam perda de até 3 tesouros.
  * **Espinhos Nômades**: fixos, causam perda de 1 tesouro.

### 👾 Monstro (`¥`)

* Surge em cada nível.
* Comportamento:
  * Nível 1 – 2: passivo
  * Nível 3 – 4: contra-ataca e reduz vidas
* Pode ser enfrentado por jogadores que não coletaram tesouros suficientes.
* Usado como **mecânica de recuperação**.

## 🌐 Multiplayer com RPC

### Arquitetura

#### 🧠 Servidor

* Mantém o estado global do jogo.
* Gerencia posições, vidas, tesouros e caixas.
* **Não possui interface gráfica**.
* Processa comandos com `sequenceNumber` para garantir execução única.

#### 🎮 Cliente

* Interface de jogo para o jogador.
* Envia ações e busca estado do servidor via chamadas RPC.
* Possui goroutine dedicada para atualizações periódicas de estado.

## 👑 Vitória

Ao fim do 4º nível, vence:

1. Quem tiver mais **tesouros acumulados**.
2. Em caso de empate, quem tiver mais **vidas restantes**.
3. Persistindo empate, vence quem sofreu menos penalidades pelos espinhos ao longo das partidas.

## 🛠️ Compilação

### 🪟 Windows

```cmd
go build -o jogo.exe
```

## ▶️ Execução

* Certifique-se de iniciar o **servidor** antes dos **clientes**.
* Deixe o arquivo `mapa.txt` no diretório raiz com um mapa válido.

```cmd
./servidor   # em um terminal
./jogo       # em outro terminal (cliente)
```

## 🧑‍💻 Grupo

* Amanda Wilmsen: [amanda.wilmsen@edu.pucrs.br](mailto:amanda.wilmsen@edu.pucrs.br)
* Killian D.B: [killian.d@edu.pucrs.br](mailto:killian.d@edu.pucrs.br)
* Luís Trein: [luis.trein@edu.pucrs.br](mailto:luis.trein@edu.pucrs.br)
* Maria Rita: [m.ritarodrigues09@gmail.com](mailto:m.ritarodrigues09@gmail.com)

## 📄 Relatório

O relatório detalha:

* A implementação do modo multiplayer
* A integração com RPC
* Os novos elementos e interações concorrentes
* Estratégias para garantir consistência e concorrência segura

📄 [Link do Relatório no DOCS](ainda não escrito)