package game

import (
	"github.com/PudgeKim/go-holdem/card"
	"github.com/google/uuid"
)

type GameMemento struct {
	RoomId uuid.UUID
	RoomLimit uint 
	TotalBet   uint64           
	CurrentBet uint64           
	IsStarted  bool             
	Deck            *card.Deck
	Status          string                    
	SmallBlindIdx uint
	BigBlindIdx   uint
	FirstPlayerIdx   uint
	IsFirstPlayerBet bool 
	CurrentPlayerIdx uint
	BetLeaderIdx     uint 
}

func NewGameMemento(game Game) *GameMemento {
	return &GameMemento{
		RoomId: game.RoomId,
		RoomLimit: game.RoomLimit,
		TotalBet: game.TotalBet,
		CurrentBet: game.CurrentBet,
		IsStarted: game.IsStarted,
		Deck: game.Deck,
		Status: game.Status,
	}
}

