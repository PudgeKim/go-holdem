package gameerror

import "errors"

var (
	InvalidPlayerTurn     = errors.New("invalid player's turn")
	DeadPlayer            = errors.New("player is dead")
	PlayerNotReady        = errors.New("player is not ready")
	NotEnoughPlayersReady = errors.New("equal or more than two players should be ready")
	LackOfPlayers         = errors.New("players should be equal or more than two")
	OverBalance           = errors.New("betting amount is more than player's balance")
	LowBetting            = errors.New("player's betting is lower than current betting amount")
	AlreadyStarted        = errors.New("game is already started")
	GiveCardsError 		  = errors.New("players couldn't hand out the cards")
	NotEnoughBalance      = errors.New("game balance must be equal or lower than user's balance")
)
