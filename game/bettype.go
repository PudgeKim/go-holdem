package game

type BetType int

const (
	Check BetType = iota
	Raise
	AllIn
)
