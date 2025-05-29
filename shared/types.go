# Diretórios
BIN_DIR = bin

# Arquivos
SERVER_SRC = server.go
CLIENT_SRC = main.go
SERVER_BIN = $(BIN_DIR)/server
CLIENT_BIN = $(BIN_DIR)/client

# Compila o servidor e o cliente
all: $(SERVER_BIN) $(CLIENT_BIN)

$(SERVER_BIN): $(SERVER_SRC)
	go build -o $(SERVER_BIN) $(SERVER_SRC)

$(CLIENT_BIN): $(CLIENT_SRC)
	go build -o $(CLIENT_BIN) $(CLIENT_SRC)

# Executa o servidor
run-server: $(SERVER_BIN)
	./$(SERVER_BIN)

# Executa o cliente
run-client: $(CLIENT_BIN)
	./$(CLIENT_BIN)

# Remove os binários
clean:
	rm -f $(SERVER_BIN) $(CLIENT_BIN)
