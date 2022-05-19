package entity

import (
	"encoding/json"

	"github.com/PudgeKim/go-holdem/card"
	"github.com/PudgeKim/go-holdem/gameconst"
	"github.com/google/uuid"
)

type Game struct {
	Memento *GameMemento

	RoomId uuid.UUID
	RoomLimit uint 
	HostName string 
	Players    []*Player // 게임에 참가하고 있는 플레이어들
	TotalBet   uint64           // 해당 게임에서 모든 플레이어들의 베팅액 합산 (새로운 게임이 시작되면 초기화됨)
	CurrentBet uint64           // 현재 턴에서 최고 베팅액 (player1이 20을 걸었고 player2가 30을 걸었으면 currentBet을 30으로 변경해줘야함)
	IsStarted  bool             // 게임이 시작됬는지

	Deck            *card.Deck
	Status          string                    // FreeFlop인지 Turn인지 등

	SmallBlindIdx uint
	BigBlindIdx   uint

	// 처음 베팅하는 플레이어의 인덱스 (Players 배열에서 인덱싱을 하기 위함)
	// 새 게임마다 1씩 증가함
	FirstPlayerIdx   uint
	IsFirstPlayerBet bool // 첫 플레이어가 베팅을 했는지를 체크
	CurrentPlayerIdx uint
	BetLeaderIdx     uint // 누군가 베팅을 추가로하면 해당 플레이어 이전까지 다시 베팅을 돌아야하기 때문에 저장해둠
}

func NewGame(roomId uuid.UUID, roomLimit uint, hostName string) *Game {
	game := Game{
		RoomId: roomId,
		RoomLimit: roomLimit,
		HostName: hostName,
		Players:    make([]*Player, 2),
		TotalBet:   0,
		CurrentBet: 0,
		IsStarted:  false,
		Deck:       card.NewDeck(),
		Status:     gameconst.FreeFlop,
	}
	memento := NewGameMemento(game)
	game.Memento = memento
	return &game
}

func (g *Game) Undo() {
	memento := g.Memento

	g.RoomId = memento.RoomId
	g.RoomLimit = memento.RoomLimit
	g.TotalBet = memento.TotalBet
	g.CurrentBet = memento.CurrentBet
	g.IsStarted = memento.IsStarted
	g.Deck = memento.Deck
	g.Status = memento.Status
	g.SmallBlindIdx = memento.SmallBlindIdx
	g.BigBlindIdx = memento.BigBlindIdx
	g.FirstPlayerIdx = memento.FirstPlayerIdx
	g.IsFirstPlayerBet = memento.IsFirstPlayerBet
	g.CurrentPlayerIdx = memento.CurrentPlayerIdx
	g.BetLeaderIdx = memento.BetLeaderIdx
}

func (g *Game) SetMemento() {
	memento := g.Memento

	memento.RoomId = g.RoomId
	memento.RoomLimit = g.RoomLimit
	memento.TotalBet = g.TotalBet
	memento.CurrentBet = g.CurrentBet
	memento.IsStarted = g.IsStarted
	memento.Deck = g.Deck
	memento.Status = g.Status
	memento.SmallBlindIdx = g.SmallBlindIdx
	memento.BigBlindIdx = g.BigBlindIdx
	memento.FirstPlayerIdx = g.FirstPlayerIdx
	memento.IsFirstPlayerBet = g.IsFirstPlayerBet
	memento.CurrentPlayerIdx = g.CurrentPlayerIdx
	memento.BetLeaderIdx = g.BetLeaderIdx
}

// redis에 struct를 저장/가져오기위해 구현해야함
func (g Game)MarshalBinary() ([]byte, error) {
    return json.Marshal(g)
}

func (g *Game) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &g); err != nil {
		return err
	}
 
	return nil
}

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