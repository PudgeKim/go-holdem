package gameerror

import (
	"errors"
	"fmt"
)

var (
	PlayerLimitationError = errors.New(fmt.Sprintf("max number of player is %d", 7))
	NoPlayerExists        = errors.New("no player exists")
	PlayerLeft            = errors.New("player left the game")
	NoPlayersLeft         = errors.New("all players left or dead or not ready")
	PlayerAlreadyExists   = errors.New("player is already in the gameroom")
	GameAlreadyStarted    = errors.New("game is already started you can't change ready status")
	InvalidHostId         = errors.New("requested host id is invalid")
)
