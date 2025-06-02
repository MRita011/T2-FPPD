package main

import "fmt"

// essa struct guarda as configs de rede do jogo
type NetworkConfig struct {
	Host           string // ip ou nome do host (ex: localhost)
	Port           string // porta usada pra conectar
	DefaultMapFile string // nome do arquivo de mapa que será usado por padrão
}

// aqui ficam algumas configurações prontas pra usar
var (
	// essa é a config padrão que o sistema usa se não passar outra
	DefaultConfig = NetworkConfig{
		Host:           "localhost",
		Port:           "8080",
		DefaultMapFile: "mapa.txt",
	}

	// outra config que pode ser usada localmente
	LocalConfig = NetworkConfig{
		Host:           "127.0.0.1",
		Port:           "3000",
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
