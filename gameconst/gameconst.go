package gameconst

const (
	FreeFlop = "FreeFlop"
	Flop     = "Flop"
	Turn     = "Turn"
	River    = "River"
	GameEnd  = "GameEnd"
)

type BetType int

const (
	Check BetType = iota
	Raise
	AllIn
)
