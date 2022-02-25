package game

type BetType int

const (
	Check BetType = iota
	Raise
	AllIn
)

type Status int

const (
	FreeFlop Status = iota
	Flop
	Turn
	River
)
