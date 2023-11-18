package p2p

type GameStatus uint32


func (g GameStatus) String() string {
	return []string{"WAITING","DEALING","PREFLOP","FLOP","TURN","RIVER"}[g]
}
const (
	GameStatusWating GameStatus = iota
	GameStatusDealing 
	GameStatusPreFlop
	GameStatusFlop
	GameStatusTurn
	GameStatusRiver
)

type GameState struct {
	isDealer bool //should be atomic accessable

	GameStatus GameStatus //should be atomic accessable
}

func NewGameState() *GameState {
	return &GameState{}
}

func (g *GameState) loop() {
	for {
		select {}
	}
}
