 // goroutine de inicialização da conexão RPC
package client

func InicializarRPC(clientAddr string) *rpc.Client {
    client, err := rpc.Dial("tcp", clientAddr)
    if err != nil {
        log.Fatal("Erro ao conectar:", err)
    }

    go func() {
        req := shared.JoinRequest{PlayerID: "algumID"}
        var res shared.JoinResponse
        client.Call("Game.Join", req, &res)
    }()

    return client
}
