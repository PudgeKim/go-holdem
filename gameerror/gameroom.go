package gameerror

import (
	"errors"
	"fmt"
	"github.com/PudgeKim/go-holdem/gameroom"
)

var (
	PlayerLimitationError = errors.New(fmt.Sprintf("max number of player is %d", gameroom.RoomLimit))
	NoPlayerExists        = errors.New("no player exists")
	PlayerLeft            = errors.New("player left the game")
	NoPlayersLeft         = errors.New("all players left or dead or not ready")
	PlayerAlreadyExists   = errors.New("player is already in the gameroom")
	GameAlreadyStarted    = errors.New("game is already started you can't change ready status")
)
