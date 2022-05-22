package entity

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/PudgeKim/go-holdem/card"
	"github.com/PudgeKim/go-holdem/errors/gameerror"
	"github.com/PudgeKim/go-holdem/gameconst"
	"github.com/google/uuid"
)

type CardCompareResult int

const (
	Player1Win CardCompareResult = iota
	Player2Win
	Draw
)

type Game struct {
	Memento *GameMemento

	RoomId uuid.UUID
	RoomLimit uint 
	HostName string 
	Players    []*Player // 게임에 참가하고 있는 플레이어들
	MinBetAmount uint64 // SmallBlind가 걸어야할 최소 금액 
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

func NewGame(roomId uuid.UUID, roomLimit uint, hostPlayer *Player, minBetAmount uint64) *Game {
	game := Game{
		RoomId: roomId,
		RoomLimit: roomLimit,
		HostName: hostPlayer.Nickname,
		Players:    make([]*Player, 0, roomLimit),
		MinBetAmount: minBetAmount,
		TotalBet:   0,
		CurrentBet: 0,
		IsStarted:  false,
		Deck:       card.NewDeck(),
		Status:     gameconst.FreeFlop,
	}
	game.Players = append(game.Players, hostPlayer)
	memento := NewGameMemento(game)
	game.Memento = memento
	return &game
}

func (g *Game) GetReadyPlayers() []*Player {
	var readyPlayers []*Player

	for _, p := range g.Players {
		if p.IsReady {
			readyPlayers = append(readyPlayers, p)
		}
	}

	return readyPlayers
}

func (g *Game) GetFirstPlayer() *Player {
	return g.Players[g.FirstPlayerIdx]
}

func (g *Game) GetSmallBlind() *Player {
	return g.Players[g.SmallBlindIdx]
}

func (g *Game) GetBigBlind() *Player {
	return g.Players[g.BigBlindIdx]
}

func (g *Game) IsPlayerExist(nickname string) bool {
	for _, p := range g.Players {
		if p.Nickname == nickname {
			return true 
		}
	}
	return false 
}

func (g *Game) FindPlayer(nickname string) *Player {
	for _, p := range g.Players {
		if p.Nickname == nickname {
			return p
		}
	}
	return nil 
}

func (g *Game) StartGame() {
	g.setPlayers()
	g.GiveCardsToPlayers()
}

// 게임이 종료되면 초기화용
func (g *Game) InitGame() {
	g.TotalBet = 0
	g.CurrentBet = 0  
	g.IsStarted = false
	g.Deck = card.NewDeck()
	g.Status = gameconst.FreeFlop
	g.IsFirstPlayerBet = false 

	g.removeLeftPlayers()
	
	for _, p := range g.Players {
		p.CurrentBet = 0
		p.TotalBet = 0
		p.IsDead = false 
		p.IsAllIn = false 
		p.Hands = nil 
		p.HandsRank = card.HandsRank(card.None)
		p.HighCard = card.None
		p.BestCards = nil 
	}
}

func (g *Game) GiveCardsToPlayers() {
	validPlayers := g.GetValidPlayers()

	for i := 0; i < len(validPlayers); i++ {
		g.Players[i].Hands = append(g.Players[i].Hands, g.Deck.GetCard(), g.Deck.GetCard())
	}
}

// 현재 플레이어들 중 나가지도 않고 죽지도 않고 준비도 된 플레이어들만 리턴
func (g *Game) GetValidPlayers() []*Player {
	var validPlayers []*Player

	for _, p := range g.Players {
		if p.IsReady && !p.IsDead && !p.IsLeft {
			validPlayers = append(validPlayers, p)
		}
	}

	return validPlayers
}

func (g *Game) ClearPlayersCurrentBet() {
	for _, p := range g.Players {
		p.CurrentBet = 0 
	}
}

// 리턴되는 배열에는 1명 이상의 승리자가 들어가게 되는데 
// 배열에 포함되는 승리자들이 2명 이상인 경우 단순히 1/n으로 나누면 안됨 
// (돈이 상대적으로 없는 사용자가 올인하고 나머지 플레이어들은 추가 베팅이 가능하기 때문)
func (g *Game) GetWinnersAndLosers() (winners []*Player, losers []*Player, err error) {
	// 죽거나 나가지 않은 플레이어들로 승리자 계산해야함 
	validPlayers := g.GetValidPlayers()

	if len(validPlayers) == 0 {
		return nil, nil, errors.New("zero player")
	}

	if len(validPlayers) == 1 {
		return []*Player{validPlayers[0]}, nil, nil 
	}

	// 승리자/패배자는 여러 명이 나올 수 있으므로 배열을 이용함

	// 베팅을 가장 많이 한 플레이어 순으로 정렬함 
	// (중간에 올인해버린 플레이어들도 winner로 고려하기 위해)
	sort.Slice(validPlayers, func(i, j int) bool {
		return validPlayers[i].TotalBet > validPlayers[j].TotalBet
	})

	// 첫번째 플레이어를 우선 승리자로 지정해놓고 이후 플레이어들과 비교해나감
	winners = append(winners, validPlayers[0])
	
	for i := 1; i < len(validPlayers); i++ {
		winner := winners[0]
		player := validPlayers[i]

		switch compare(winner, player) {
		case Player1Win:
			continue
		case Player2Win:
			if winner.TotalBet == player.TotalBet {
				// 배열에 기존 승리자가 여러명 있을 수도 있으니 아예 비워야함
				winners = nil
				winners = append(winners, player)
			} else {
				// 지금 이긴 플레이어는 중간에 올인하여 돈이 부족한 플레이어임 
				// 즉 지금 이긴 플레이어에게 진 플레이어들도 나머지 돈을 배분 받을 수 있음 
				winners = append(winners, player)
			}
			
		case Draw:
			// 기존 승리자들과 같은 패이므로 승리자에 추가
			winners = append(winners, player)
		}
	}

	// 승리자들에 대해 돈이 적은 순으로 정렬 
	// 예를 들어 0번 인덱스의 총 베팅액이 100이고 
	// 1번 인덱스의 총 베팅액이 150이라면 
	// 0번 인덱스에게 각 플레이어가 100만큼씩 주고 
	// 1번 인덱스에게 50씩 줘야함 
	sort.Slice(winners, func(i, j int) bool {
		return winners[i].TotalBet < winners[j].TotalBet
	})

	// 패자들 생성 
	for _, p := range g.GetReadyPlayers() {
		isExist := false 
		for _, winner := range winners {
			if p.Nickname == winner.Nickname {
				isExist = true 
				break 
			}
		}

		if !isExist {
			losers = append(losers, p)
		}
	}

	return winners, losers, nil
}

func (g *Game) setPlayers() ([]string, error) {
	
	playerCnt := 0
	for _, p := range g.Players {
		if p != nil {
			playerCnt++
		}
	}

	if playerCnt < 2 {
		return nil, gameerror.LackOfPlayers
	}

	var readyPlayersName []string 
	readyCnt := 0
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].IsReady {
			readyPlayersName = append(readyPlayersName, g.Players[i].Nickname)
			readyCnt += 1
		}
	}
	if readyCnt < 2 {
		return nil, gameerror.NotEnoughPlayersReady
	}

	// 게임이 처음 시작되는 경우와 이미 몇판 진행되고 있는지에 따라 나눔
	// (이미 진행되고 있었다면 이전 게임의 플레이어들의 순서를 기준으로 세팅해야되기 때문)
	var isFirstGame bool

	if g.SmallBlindIdx == 0 && g.BigBlindIdx == 0 {
		isFirstGame = true
	} else {
		isFirstGame = false
	}

	if isFirstGame {
		if len(g.Players) > 2 {
			g.SmallBlindIdx = getReadyPlayerIdx(g.Players, 0)
			g.BigBlindIdx = getReadyPlayerIdx(g.Players, g.SmallBlindIdx + 1)
			g.FirstPlayerIdx = getReadyPlayerIdx(g.Players, g.BigBlindIdx + 1)
		} else { // 플레이어가 2명 밖에 없는 경우 (이 함수의 첫 부분 검사에 의해 최소 2명은 Ready 상태이므로 0번과 1번 인덱스가 모두 Ready 상태임)
			g.SmallBlindIdx = 0
			g.BigBlindIdx = 1
			g.FirstPlayerIdx = g.SmallBlindIdx
		}
	} else { // 기존에 진행되던 순서가 있는 경우
		if len(g.Players) > 2 {
			g.SmallBlindIdx = getReadyPlayerIdx(g.Players, g.SmallBlindIdx + 1)
			g.BigBlindIdx = getReadyPlayerIdx(g.Players, g.SmallBlindIdx + 1)
			g.FirstPlayerIdx = getReadyPlayerIdx(g.Players, g.BigBlindIdx + 1)
		} else { // 플레이어가 2명만 있는 경우
			g.SmallBlindIdx, g.BigBlindIdx = g.BigBlindIdx, g.SmallBlindIdx
			g.FirstPlayerIdx = g.SmallBlindIdx
		}
	}

	g.CurrentPlayerIdx = g.FirstPlayerIdx
	g.BetLeaderIdx = g.FirstPlayerIdx
	g.IsFirstPlayerBet = false

	return readyPlayersName, nil
}


// 나간 플레이어들 고려해서 인덱스 변경해야함
// smallBlind와 bigBlind가 바뀌었는지 리턴
func (g *Game) removeLeftPlayers() {
	var indexesToRemove []int

	for idx, p := range g.Players {
		if p.IsLeft {
			if g.SmallBlindIdx == uint(idx) {
				g.SmallBlindIdx = getNextIdx(g.Players, uint(idx))
			}
			if g.BigBlindIdx == uint(idx) {
				g.BigBlindIdx = getNextIdx(g.Players, uint(idx))
			}

			indexesToRemove = append(indexesToRemove, idx)
		}
	}

	for _, idx := range indexesToRemove {
		g.Players = removePlayerByIndex(g.Players, idx)
	}
}

// 함수 인자로 들어온 인덱스에 해당하는 플레이어가 Ready 상태면 해당 인덱스를 리턴하고
// 아니라면 다음 플레이어들 중에서 가장 빠른 순서인 Ready 상태인 플레이어에 해당 인덱스를 리턴
func getReadyPlayerIdx(players []*Player, idx uint) uint {
	readyIdx := idx
	for i := 0; i < len(players); i++ {
		if players[readyIdx].IsReady {
			break
		}
		readyIdx = getNextIdx(players, readyIdx)
	}
	return readyIdx
}

// 현재 인덱스에 +1한 인덱스 값을 리턴
// 원형 리스트처럼 적용시켜야하므로 mod 연산 이용
func getNextIdx(players []*Player, idx uint) uint {
	nextIdx := (idx + 1) % uint(len(players))
	return nextIdx
}

// 준비를 안해서 게임을 진행중이지 않거나 죽거나 나가는 사람이 있기 때문에 단순히 currentPlayerIdx를 1씩 증가하면 오류가 생김
func (g *Game) GetNextPlayerIdx() (uint, error) {
	idx := g.CurrentPlayerIdx
	for i := 0; i < len(g.Players); i++ {
		nextPlayer := g.Players[getReadyPlayerIdx(g.Players, idx+1)]
		idx = getReadyPlayerIdx(g.Players, idx + 1)

		if nextPlayer.IsReady && !nextPlayer.IsDead && !nextPlayer.IsLeft {
			break
		}
	}

	// 인덱스가 다시 돌아왔단 것은 모든 플레이어가 죽었거나 나갔다는 것
	if idx == g.CurrentPlayerIdx {
		return 0, gameerror.NoPlayersLeft
	}

	return idx, nil
}

func removePlayerByIndex(players []*Player, s int) []*Player {
	return append(players[:s], players[s+1:]...)
}

// 두 플레이어 간의 bestCards를 비교해서 이긴 플레이어를 리턴
// 둘이 같다면 Draw를 리턴
// ** 각 플레이어들의 bestCards는 정렬되어 있음 (bestCards를 만드는 과정에서 정렬 함수가 쓰임)
func compare(player1 *Player, player2 *Player) CardCompareResult {
	if player1.HandsRank > player2.HandsRank {
		return Player1Win
	} else if player1.HandsRank < player2.HandsRank {
		return Player2Win
	} else {
		if player1.HighCard > player2.HighCard {
			return Player1Win
		} else if player1.HighCard < player2.HighCard {
			return Player2Win
		} else {
			for i := len(player1.BestCards) - 1; i >= 0; i-- {
				if player1.BestCards[i].Rank > player2.BestCards[i].Rank {
					return Player1Win
				} else if player1.BestCards[i].Rank < player2.BestCards[i].Rank {
					return Player2Win
				}
			}
			return Draw
		}
	}
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