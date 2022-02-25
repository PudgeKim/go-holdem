package game

import "errors"

var (
	NoPlayerExists    = errors.New("no player exists")
	InvalidPlayerTurn = errors.New("invalid player's turn")
	DeadPlayer        = errors.New("player is dead")
	PlayerLeft        = errors.New("player left the game")
	OverBalance       = errors.New("betting amount is more than player's balance")
	LowBetting        = errors.New("player's betting is lower than current betting amount")
	NoPlayersLeft     = errors.New("all players left or dead")
)
