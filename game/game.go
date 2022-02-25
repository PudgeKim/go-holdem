package game

import (
	"github.com/PudgeKim/card"
	"github.com/PudgeKim/player"
)

type Game struct {
	players    []*player.Player // 게임에 참가하고 있는 플레이어들
	totalBet   uint64           // 해당 게임에서 모든 플레이어들의 베팅액 합산 (새로운 게임이 시작되면 초기화됨)
	currentBet uint64           // 현재 턴에서 최고 베팅액 (player1이 20을 걸었고 player2가 30을 걸었으면 currentBet을 30으로 변경해줘야함)
	isStarted  bool             // 게임이 시작됬는지

	deck    *card.Deck
	status  Status       // FreeFlop인지 Turn인지 등
	betChan chan Request // 베팅 등을 포함해서 게임에 필요한 요청들을 받고 처리함

	smallBlindIdx uint
	bigBlindIdx   uint

	// 처음 베팅하는 플레이어의 인덱스 (players 배열에서 인덱싱을 하기 위함)
	// 새 게임마다 1씩 증가함
	firstPlayerIdx   uint
	isFirstPlayerBet bool // 첫 플레이어가 베팅을 했는지를 체크
	currentPlayerIdx uint
	betLeaderIdx     uint // 누군가 베팅을 추가로하면 해당 플레이어 이전까지 다시 베팅을 돌아야하기 때문에 저장해둠
}

func New(players []*player.Player) Game {
	var startPlayerIdx uint

	if len(players) > 2 {
		startPlayerIdx = 2
	} else {
		startPlayerIdx = 0
	}

	// 각 플레이어들의 순서 지정해줌
	for i := 0; i < len(players); i++ {
		players[i].Turn = uint(i)
	}

	return Game{
		players:          players,
		totalBet:         0,
		currentBet:       0,
		isStarted:        false,
		deck:             card.NewDeck(),
		status:           FreeFlop,
		betChan:          make(chan Request),
		smallBlindIdx:    0,
		bigBlindIdx:      1,
		firstPlayerIdx:   startPlayerIdx,
		isFirstPlayerBet: false,
		currentPlayerIdx: startPlayerIdx,
		betLeaderIdx:     startPlayerIdx,
	}
}

// Start 게임 시작 요청이 들어오면 해당 함수 실행
func (g *Game) Start() {
	for {
		req := <-g.betChan

		if req.IsLeave {
			p, err := g.findPlayer(req.PlayerName)
			if err != nil {
				// 어딘가로 보내야함
				// TODO
			}
			p.IsLeft = true
		}

		nextPlayerName, isBetEnd, err := g.HandleBet(req)
		if err != nil {
			// 어딘가로 에러 보내야함
			// TODO
		}

		switch g.status {
		case FreeFlop:
			if isBetEnd {
				g.status = Flop
			}
			// 어딘가로 HandleBet의 결과값을 보내야함
			// (그 어딘가에서 프론트로 전송해줌)
		case Flop:
			if isBetEnd {
				g.status = Turn
			}
		case Turn:
			if isBetEnd {
				g.status = River
			}
		case River:
			if isBetEnd {
				// 게임 종료
				// 나간 플레이어들 빼줘야함
				// 나간 플레이어들 고려해서 인덱스 변경해야함
				// firstPlayer값 1증가 시켜줘야함, smallBlind, bigBlind 값들도 증가시켜줘야함
				// 인덱스 값들, 베팅 값들 모두 초기화
				// status도 변경시켜줘야함
			}
		}
	}

}

// HandleBet 모든 플레이어들의 베팅이 종료되는 경우면 true를 리턴함
// 다음 플레이어가 누군지도 리턴
func (g *Game) HandleBet(req Request) (string, bool, error) {
	p, err := g.findPlayer(req.PlayerName)
	if err != nil {
		return "", false, err
	}

	expectedPlayer := g.players[g.currentPlayerIdx]
	if p.Nickname != expectedPlayer.Nickname {
		return "", false, InvalidPlayerTurn
	}
	if p.IsDead {
		return "", false, DeadPlayer
	}
	if p.IsLeft {
		return "", false, PlayerLeft
	}

	// 플레이어가 베팅 대신 죽은 경우
	if req.IsDead {
		p.IsDead = true
		nextPlayerIdx, err := g.getNextPlayerIdx()
		if err != nil {
			return "", false, err
		}
		nextPlayer := g.players[nextPlayerIdx].Nickname
		return nextPlayer, false, nil
	}

	betType, err := g.isValidBet(p, req.BetAmount)
	if err != nil {
		return "", false, err
	}
	if betType == AllIn {
		p.IsAllIn = true
	}

	p.CurrentBet += req.BetAmount
	p.TotalBet += req.BetAmount

	nextPlayerIdx, err := g.getNextPlayerIdx()
	if err != nil {
		return "", false, err
	}

	nextPlayerName := g.players[nextPlayerIdx].Nickname

	// 현재 베팅한 플레이어가 베팅한 금액에 따라 베팅리더인지 체크 후에 현재 베팅 턴을 종료할지 결정
	// (현재 플레이어가 베팅리더가 아니고, betLeader 이전 플레이어면 플레이어들의 베팅이 종료됨)
	// 베팅이 종료되면 다음 베팅을 위해서 player들의 currentBet을 초기화시켜주어야함
	if req.BetAmount > g.currentBet { // 현재 플레이어가 베팅 리더가 되는 경우
		g.currentBet = p.CurrentBet
		g.betLeaderIdx = p.Turn
		g.currentPlayerIdx = nextPlayerIdx
		return nextPlayerName, false, nil
	} else {
		// 베팅 종료 조건 달성한 경우
		if g.getNextIdx(p.Turn) == g.betLeaderIdx {
			g.currentPlayerIdx = g.firstPlayerIdx // 다음 베팅을 위해서 초기화
			g.clearPlayersCurrentBet()
			return nextPlayerName, true, nil
		}

		// 베팅은 종료되지 않고 다음 플레이어가 베팅해야함
		g.currentPlayerIdx = nextPlayerIdx
		return nextPlayerName, false, nil
	}
}

func (g Game) isValidBet(p *player.Player, betAmount uint64) (BetType, error) {
	if p.GameBalance == p.TotalBet+betAmount {
		return AllIn, nil
	}
	if p.GameBalance < p.TotalBet+betAmount {
		return -1, OverBalance
	}
	if g.currentBet > p.CurrentBet+betAmount {
		return -1, LowBetting
	}
	if g.currentBet < p.CurrentBet+betAmount {
		return Raise, nil
	}

	return Check, nil
}

func (g *Game) giveCardsToPlayers() {
	for i := 0; i < len(g.players); i++ {
		g.players[i].Hands = append(g.players[i].Hands, g.deck.GetCard(), g.deck.GetCard())
	}
}

//
//func (g *Game) GetAllPlayers() []*player.Player {
//	return g.players
//}

func (g *Game) AddPlayer(player *player.Player) {
	g.players = append(g.players, player)
}

func (g *Game) RemovePlayer(player player.Player) (*player.Player, error) {
	removeIndex := -1
	for i := 0; i < len(g.players); i++ {
		if player.Nickname == g.players[i].Nickname {
			removeIndex = i
			break
		}
	}

	if removeIndex == -1 {
		return nil, NoPlayerExists
	}

	removedPlayer := g.players[removeIndex]
	g.players = append(g.players[:removeIndex], g.players[removeIndex+1:]...)

	return removedPlayer, nil
}

func (g Game) findPlayer(nickname string) (*player.Player, error) {
	for _, p := range g.players {
		if p.Nickname == nickname {
			return p, nil
		}
	}
	return nil, NoPlayerExists
}

// 죽거나 나가는 사람이 있기 때문에 단순히 currentPlayerIdx를 1씩 증가하면 오류가 생김
func (g Game) getNextPlayerIdx() (uint, error) {
	idx := g.currentPlayerIdx
	for i := 0; i < len(g.players); i++ {
		nextPlayer := g.players[g.getNextIdx(idx)]
		idx = g.getNextIdx(idx)

		if !nextPlayer.IsDead && !nextPlayer.IsLeft {
			break
		}
	}

	// 인덱스가 다시 돌아왔단 것은 모든 플레이어가 죽었거나 나갔다는 것
	if idx == g.currentPlayerIdx {
		return 0, NoPlayersLeft
	}

	return idx, nil
}

// 현재 인덱스에 +1한 인덱스 값을 리턴
// 원형 리스트처럼 적용시켜야하므로 mod 연산 이용
func (g Game) getNextIdx(idx uint) uint {
	nextIdx := (idx + 1) % uint(len(g.players))
	return nextIdx
}

func (g *Game) clearPlayersCurrentBet() {
	for _, p := range g.players {
		p.CurrentBet = 0
	}
}
