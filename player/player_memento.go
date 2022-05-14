package player

import "github.com/PudgeKim/go-holdem/card"

type PlayerMemento struct {
	Id 			 int64 
	Nickname     string
	IsReady      bool 
	IsDead       bool
	IsLeft       bool 
	IsAllIn      bool
	TotalBalance uint64        
	GameBalance  uint64         
	TotalBet     uint64        
	CurrentBet   uint64         
	Hands        []card.Card    
	HandsRank    card.HandsRank 
	HighCard     card.Rank      
	BestCards    []card.Card    
}

func NewPlayerMemento(player Player) *PlayerMemento {
	return &PlayerMemento{
		Id: player.Id,
		Nickname: player.Nickname,
		IsReady: player.IsReady,
		IsDead: player.IsDead,
		IsLeft: player.IsLeft,
		IsAllIn: player.IsAllIn,
		TotalBalance: player.TotalBalance,
		GameBalance: player.GameBalance,
		TotalBet: player.TotalBet,
		CurrentBet: player.CurrentBet,
	}
}