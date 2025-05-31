package common

type JoinRequest struct {
	Name string
}

type JoinResponse struct {
	ID     int
	Player Player
}

type StateRequest struct {
	PlayerID int
}

type GameState struct {
	MapWidth  int
	MapHeight int
	Players   []Player
	Traps     []Element
	Treasures []Element
}

type Player struct {
	ID     int
	Name   string
	X, Y   int
	Symbol string
}

type Element struct {
	X, Y   int
	Symbol string
}
