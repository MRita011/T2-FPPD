# caminhos
CLIENT_DIR = client
SERVER_DIR = server
COMMON_DIR = common

# arquivos
CLIENT_FILES=$(CLIENT_DIR)/main.go $(CLIENT_DIR)/rpc_client.go
SERVER_FILES=$(SERVER_DIR)/main.go $(SERVER_DIR)/state.go

# binários
BIN_DIR=bin
CLIENT_BIN=$(BIN_DIR)/cliente
SERVER_BIN=$(BIN_DIR)/servidor

# inicialização do módulo (opcional)
init:
	go mod init fppd-jogo
	go mod tidy

# build
build-client:
	go build -o $(CLIENT_BIN) $(CLIENT_FILES)

build-server:
	go build -o $(SERVER_BIN) $(SERVER_FILES)

build: build-client build-server

# run
run-client:
	go run $(CLIENT_FILES)

run-server:
	go run $(SERVER_FILES)

# limpeza
clean:
	rm -rf $(BIN_DIR)

# cria a pasta bin (se nao existir)
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# build + pasta
build-all: $(BIN_DIR) build
