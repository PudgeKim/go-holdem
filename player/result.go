package player

type CardCompareResult int

const (
	Player1Win CardCompareResult = iota
	Player2Win
	Draw
)
