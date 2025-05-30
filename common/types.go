package common

// tipos compartilhados entre o servidor e o cliente
type JoinRequest struct {
	Name string
}

type JoinResponse struct {
	ID int
}

type Player struct {
	ID   int
	Name string
}
