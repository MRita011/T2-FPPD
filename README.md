````markdown
# 🎮 Jogo Concorrente e Multiplayer em Go

## 🧭 Visão Geral

Este projeto evolui um jogo de aventura em Go com interface textual, inicialmente focado em concorrência, para uma versão multiplayer usando RPC. A nova implementação permite que múltiplos jogadores explorem um mapa dinâmico, enfrentando armadilhas, monstros e coletando tesouros, tudo conectado a um **servidor central**.

> 🔗 Código original da disciplina de Fundamentos de Processamento Paralelo e Distribuído:  
> [https://github.com/mvneves/fppd-jogo](https://github.com/mvneves/fppd-jogo)

---

## 🎯 Objetivos

- Adicionar elementos concorrentes interativos  
- Tornar o jogo multiplayer com **RPC cliente-servidor**  
- Criar **níveis progressivos**, desafios e critérios de vitória  
- Usar conceitos de concorrência segura (`goroutines`, `canais`, `mutexes`)  
- Garantir execução única de comandos via `sequenceNumber`

---

## 🕹️ Como Jogar

- `W`, `A`, `S`, `D`: mover personagem  
- `E`: interagir  
- `ESC`: sair  
- Cada jogador tem um símbolo próprio  
- Clientes se conectam a um **servidor central**

---
## 🌍 Mapa e Níveis

O jogo tem **4 níveis**, cada um com 40 tesouros (160 no total). Para avançar:

| Tesouros Coletados         | Resultado                         |
|----------------------------|-----------------------------------|
| ≥ 20                       | Avança normalmente                |
| 15–19 + enfrenta monstro   | Avança, perde 5 tesouros          |
| < 15 + enfrenta monstro    | Avança, perde metade dos tesouros |
| Nenhuma das condições      | Fica no nível atual               |

---

## ⚙️ Elementos Concorrentes

### 💰 Tesouros

- Coletáveis e acumuláveis  
- Protegidos por `mutex`  
- Contabilizados por jogador  

### 💣 Armadilhas

- **Espinhos Andarilhos**: móveis, perdem até 3 tesouros  
- **Espinhos Nômades**: fixos, perdem 1 tesouro  

### 👾 Monstro (`¥`)

- Um por nível  
- **Níveis 1–2**: passivo  
- **Níveis 3–4**: contra-ataca  
- Jogadores com poucos tesouros enfrentam o monstro  
- Pode ajudar na recuperação  

---

## 🌐 Multiplayer com RPC

### 🧠 Servidor

- Mantém estado global  
- Gerencia posições, vidas, tesouros  
- **Sem interface gráfica**  
- Garante execução única via `sequenceNumber`  

### 🎮 Cliente

- Interface do jogador  
- Envia ações e recebe atualizações via RPC  
- Usa `goroutine` para atualizações contínuas  

---

## 👑 Vitória

No final do **4º nível**, vence:

1. Quem tiver mais **tesouros**  
2. Se empate, quem tiver mais **vidas**  
3. Persistindo empate, quem teve menos **penalidades**

---

## 🛠️ Compilação

### 🪟 Windows

```cmd
go build -o jogo.exe
````

---

## ▶️ Execução

### 1. Servidor

```bash
go run ./server
```

### 2. Cliente

```bash
go run ./client
```

---

## 🧑‍💻 Grupo

* Amanda Wilmsen – [amanda.wilmsen@edu.pucrs.br](mailto:amanda.wilmsen@edu.pucrs.br)
* Killian D.B – [killian.d@edu.pucrs.br](mailto:killian.d@edu.pucrs.br)
* Luís Trein – [luis.trein@edu.pucrs.br](mailto:luis.trein@edu.pucrs.br)
* **Maria Rita** – [m.ritarodrigues09@gmail.com](mailto:m.ritarodrigues09@gmail.com)

---

## 📄 Relatório

📄 *\[Link do relatório]* (ainda não escrito)

```