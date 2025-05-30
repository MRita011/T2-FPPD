# Caminhos
CLIENT_DIR=client
SERVER_DIR=server

# Comandos
all: build-client build-server

build-client:
	go build -o bin/cliente $(CLIENT_DIR)/main.go $(CLIENT_DIR)/rpc_client.go

build-server:
	go build -o bin/servidor $(SERVER_DIR)/main.go

run-client:
	go run $(CLIENT_DIR)/main.go

run-server:
	go run $(SERVER_DIR)/main.go

clean:
	rm -rf bin
