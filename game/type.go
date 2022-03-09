package game

type BetType int

const (
	Check BetType = iota
	Raise
	AllIn
)

const (
	FreeFlop = "FreeFlop"
	Flop     = "Flop"
	Turn     = "Turn"
	River    = "River"
	GameEnd  = "GameEnd"
)
