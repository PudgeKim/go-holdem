package entity

import (
	"github.com/PudgeKim/go-holdem/card"
)

type Player struct {
	Memento *PlayerMemento
	Id 			 int64 // User struct의 id
	Nickname     string
	IsReady      bool // 게임준비
	IsDead       bool
	IsLeft       bool // 게임 중간에 나간 경우 여기에 우선 체크를 해두고 게임이 종료되면 실제로 나가게 처리함 (인덱스가 꼬이는거 방지하기 위해)
	IsAllIn      bool
	TotalBalance uint64         // 매 게임 또는 플레이어가 죽거나 나가는 경우 갱신
	GameBalance  uint64         // 게임 참가시에 들고갈 돈 (매 게임 또는 플레이어가 죽거나 나가는 경우 갱신)
	TotalBet     uint64         // 해당 게임에서 누적 베팅액
	CurrentBet   uint64         // 현재 턴에서 베팅한 금액
	Hands        []card.Card    // 처음 받는 2장의 카드
	HandsRank    card.HandsRank // 족보 (fullHouse인지 onePair인지.. 등)
	HighCard     card.Rank      // 예를 들어 33322 fullHouse면 highCard는 3
	BestCards    []card.Card    // 필드에 카드가 모두 오픈되었을 때 hands까지 합쳐서 가장 좋은 5장의 카드들
}

func New(id int64, nickname string, totalBalance, gameBalance uint64) *Player {
	player :=  &Player{
		Id: 			id,
		Nickname:     nickname,
		IsReady:      false,
		IsDead:       false,
		IsLeft:       false,
		IsAllIn:      false,
		TotalBalance: totalBalance,
		GameBalance:  gameBalance,
		TotalBet:     0,
		CurrentBet:   0,
	}
	memento := NewPlayerMemento(*player)
	player.Memento = memento

	return player
}

func (p *Player) Undo() {
	memento := p.Memento

	p.Id=  memento.Id			 
	p.Nickname= memento.Nickname
	p.IsReady=        memento.IsReady
	p.IsDead=        memento.IsDead
	p.IsLeft=        memento.IsLeft
	p.IsAllIn=       memento.IsAllIn
	p.TotalBalance=        memento.TotalBalance
	p.GameBalance=       memento.GameBalance
	p.TotalBet=          memento.TotalBet
	p.CurrentBet=       memento.CurrentBet
	p.Hands=        memento.Hands
	p.HandsRank=     memento.HandsRank
	p.HighCard=         memento.HighCard
	p.BestCards=        memento.BestCards
}

func (p *Player) SetMemento() {
	memento := p.Memento

	memento.Id=  p.Id			 
	memento.Nickname= p.Nickname
	memento.IsReady=        p.IsReady
	memento.IsDead=        p.IsDead
	memento.IsLeft=        p.IsLeft
	memento.IsAllIn=       p.IsAllIn
	memento.TotalBalance=        p.TotalBalance
	memento.GameBalance=       p.GameBalance
	memento.TotalBet=          p.TotalBet
	memento.CurrentBet=       p.CurrentBet
	memento.Hands=        p.Hands
	memento.HandsRank=     p.HandsRank
	memento.HighCard=         p.HighCard
	memento.BestCards=        p.BestCards
}

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