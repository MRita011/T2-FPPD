package main

import "fmt"

// essa struct guarda as configs de rede do jogo
type NetworkConfig struct {
	Host           string // ip ou nome do host (ex: localhost)
	Port           string // porta usada pra conectar
	DefaultMapFile string // nome do arq/mapa padrão
}

// multiplayer: utilizamos o ip de uma das maquinas
var (
	DefaultConfig = NetworkConfig{
		Host:           "192.168.29.120", // ip do luis
		Port:           "8080",
		DefaultMapFile: "mapa.txt",
	}

	// singleplayer: usamos o localhost
	LocalConfig = NetworkConfig{
		Host:           "localhost",
		Port:           "8080",
		DefaultMapFile: "mapa.txt",
	}
)

// pega o endereço completo (ip + porta) pra se conectar no servidor
func (nc *NetworkConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s", nc.Host, nc.Port)
}

// pega o endereço que o servidor usa pra escutar conexões
func (nc *NetworkConfig) GetListenAddress() string {
	return ":" + nc.Port
}

// função pra criar uma nova configuração personalizada de forma rápida
func NewConfig(host, port, mapFile string) NetworkConfig {
	return NetworkConfig{
		Host:           host,
		Port:           port,
		DefaultMapFile: mapFile,
	}
}
