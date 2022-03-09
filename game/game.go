package game

import (
	"github.com/PudgeKim/go-holdem/card"
	"github.com/PudgeKim/go-holdem/channels"
	"github.com/PudgeKim/go-holdem/gameerror"
	"github.com/PudgeKim/go-holdem/player"
)

// 게임이 처음 시작되는 경우
// smallBlind  = 0번째 인덱스에 해당하는 플레이어 (준비를 한 경우)
// bigBlind    = 1번째 인덱스에 해당하는 플레이어 (준비를 한 경우)
// firstPlayer = bigBlind 다음 인덱스에 해당하는 플레이어 (만약 플레이어가 2명이라면 smallBlind에 해당되는 플레이어)

type Game struct {
	Players    []*player.Player // 게임에 참가하고 있는 플레이어들
	totalBet   uint64           // 해당 게임에서 모든 플레이어들의 베팅액 합산 (새로운 게임이 시작되면 초기화됨)
	currentBet uint64           // 현재 턴에서 최고 베팅액 (player1이 20을 걸었고 player2가 30을 걸었으면 currentBet을 30으로 변경해줘야함)
	IsStarted  bool             // 게임이 시작됬는지

	deck            *card.Deck
	status          string                    // FreeFlop인지 Turn인지 등
	BetChan         chan channels.BetInfo     // 베팅 등을 포함해서 게임에 필요한 요청들을 받고 처리함
	BetResponseChan chan channels.BetResponse // 베팅을 처리하고 프론트에 응답을 주기 위해서 이 채널을 통해 전달함

	smallBlindIdx uint
	bigBlindIdx   uint

	// 처음 베팅하는 플레이어의 인덱스 (Players 배열에서 인덱싱을 하기 위함)
	// 새 게임마다 1씩 증가함
	firstPlayerIdx   uint
	isFirstPlayerBet bool // 첫 플레이어가 베팅을 했는지를 체크
	currentPlayerIdx uint
	betLeaderIdx     uint // 누군가 베팅을 추가로하면 해당 플레이어 이전까지 다시 베팅을 돌아야하기 때문에 저장해둠
}

func New() *Game {
	return &Game{
		Players:    make([]*player.Player, 2),
		totalBet:   0,
		currentBet: 0,
		IsStarted:  false,
		deck:       card.NewDeck(),
		status:     FreeFlop,
		BetChan:    make(chan channels.BetInfo),
	}
}

// SetPlayers 준비된 플레이어들의 순서와 smallBlind, bigBlind를 지정해줌
// 게임이 처음 시작됬는지 아닌지에 따라 구별함
func (g *Game) SetPlayers() error {
	if len(g.Players) < 2 {
		return gameerror.LackOfPlayers
	}

	readyCnt := 0
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].IsReady {
			readyCnt += 1
		}
	}
	if readyCnt < 2 {
		return gameerror.NotEnoughPlayersReady
	}

	// 게임이 처음 시작되는 경우와 이미 몇판 진행되고 있는지에 따라 나눔
	// (이미 진행되고 있었다면 이전 게임의 플레이어들의 순서를 기준으로 세팅해야되기 때문)
	var isFirstGame bool

	if g.smallBlindIdx == 0 && g.bigBlindIdx == 0 {
		isFirstGame = true
	} else {
		isFirstGame = false
	}

	if isFirstGame {
		if len(g.Players) > 2 {
			g.smallBlindIdx = g.getReadyPlayerIdx(0)
			g.bigBlindIdx = g.getReadyPlayerIdx(g.smallBlindIdx + 1)
			g.firstPlayerIdx = g.getReadyPlayerIdx(g.bigBlindIdx + 1)
		} else { // 플레이어가 2명 밖에 없는 경우 (이 함수의 첫 부분 검사에 의해 최소 2명은 Ready 상태이므로 0번과 1번 인덱스가 모두 Ready 상태임)
			g.smallBlindIdx = 0
			g.bigBlindIdx = 1
			g.firstPlayerIdx = g.smallBlindIdx
		}
	} else { // 기존에 진행되던 순서가 있는 경우
		if len(g.Players) > 2 {
			g.smallBlindIdx = g.getReadyPlayerIdx(g.smallBlindIdx + 1)
			g.bigBlindIdx = g.getReadyPlayerIdx(g.smallBlindIdx + 1)
			g.firstPlayerIdx = g.getReadyPlayerIdx(g.bigBlindIdx + 1)
		} else { // 플레이어가 2명만 있는 경우
			g.smallBlindIdx, g.bigBlindIdx = g.bigBlindIdx, g.smallBlindIdx
			g.firstPlayerIdx = g.smallBlindIdx
		}
	}

	g.currentPlayerIdx = g.firstPlayerIdx
	g.betLeaderIdx = g.firstPlayerIdx
	g.isFirstPlayerBet = false

	return nil
}

// Start 게임 시작 요청이 들어오면 해당 함수 실행
func (g *Game) Start() {
	for {
		req := <-g.BetChan

		// 게임 준비 고려해서 코드 짜야됨

		var betResponse channels.BetResponse

		nextPlayerName, isPlayerDead, playerCurBet, playerTotBet, gameCurBet, gameTotBet, isBetEnd, err := g.HandleBet(req)
		if err != nil {
			betResponse.Error = err
			g.BetResponseChan <- betResponse // 에러 전파 (handler에서 해당 에러를 받고 프론트에 에러 넘겨줌)
			continue
		}

		betResponse.IsBetEnd = isBetEnd
		betResponse.IsPlayerDead = isPlayerDead
		betResponse.PlayerCurrentBet = playerCurBet
		betResponse.PlayerTotalBet = playerTotBet
		betResponse.GameCurrentBet = gameCurBet
		betResponse.GameTotalBet = gameTotBet
		betResponse.NextPlayerName = nextPlayerName

		switch g.status {
		case FreeFlop:
			if isBetEnd {
				g.status = Flop
			}
			betResponse.GameStatus = FreeFlop
			g.BetResponseChan <- betResponse
			// 어딘가로 HandleBet의 결과값을 보내야함
			// (그 어딘가에서 프론트로 전송해줌)
		case Flop:
			if isBetEnd {
				g.status = Turn
			}
			betResponse.GameStatus = Flop
			g.BetResponseChan <- betResponse
		case Turn:
			if isBetEnd {
				g.status = River
			}
			betResponse.GameStatus = Turn
			g.BetResponseChan <- betResponse
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
// 다음 플레이어, 현재플레이어 isDead, 현재 플레이어의 currentBet, totalBet, 현재 게임의 currentBet, totalBet, 베팅종료 리턴
func (g *Game) HandleBet(betInfo channels.BetInfo) (string, bool, uint64, uint64, uint64, uint64, bool, error) {
	p, err := g.findPlayer(betInfo.PlayerName)
	if err != nil {
		return "", false, 0, 0, 0, 0, false, err
	}

	expectedPlayer := g.Players[g.currentPlayerIdx]
	if p.Nickname != expectedPlayer.Nickname {
		return "", false, 0, 0, 0, 0, false, gameerror.InvalidPlayerTurn
	}
	if !p.IsReady {
		return "", false, 0, 0, 0, 0, false, gameerror.PlayerNotReady
	}
	if p.IsDead {
		return "", false, 0, 0, 0, 0, false, gameerror.DeadPlayer
	}
	if p.IsLeft {
		return "", false, 0, 0, 0, 0, false, gameerror.PlayerLeft
	}

	// 플레이어가 베팅하는 대신 죽은 경우
	if betInfo.IsDead {
		p.IsDead = true
		nextPlayerIdx, err := g.getNextPlayerIdx()
		if err != nil {
			return "", false, 0, 0, 0, 0, false, err
		}
		nextPlayer := g.Players[nextPlayerIdx].Nickname
		return nextPlayer, true, p.CurrentBet, p.TotalBet, g.currentBet, g.totalBet, false, nil
	}

	betType, err := g.isValidBet(p, betInfo.BetAmount)
	if err != nil {
		return "", false, 0, 0, 0, 0, false, err
	}
	if betType == AllIn {
		p.IsAllIn = true
	}

	p.CurrentBet += betInfo.BetAmount
	p.TotalBet += betInfo.BetAmount

	nextPlayerIdx, err := g.getNextPlayerIdx()
	if err != nil {
		return "", false, 0, 0, 0, 0, false, err
	}

	nextPlayerName := g.Players[nextPlayerIdx].Nickname
	currentPlayerIdx, err := g.getPlayerIdx(p.Nickname)
	if err != nil {
		return "", false, 0, 0, 0, 0, false, err
	}

	// 현재 베팅한 플레이어가 베팅한 금액에 따라 베팅리더인지 체크 후에 현재 베팅 턴을 종료할지 결정
	// (현재 플레이어가 베팅리더가 아니고, betLeader 이전 플레이어면 플레이어들의 베팅이 종료됨)
	// 베팅이 종료되면 다음 베팅을 위해서 player들의 currentBet을 초기화시켜주어야함
	if betInfo.BetAmount > g.currentBet { // 현재 플레이어가 베팅 리더가 되는 경우
		g.currentBet = p.CurrentBet
		g.betLeaderIdx = currentPlayerIdx
		g.currentPlayerIdx = nextPlayerIdx
		return nextPlayerName, false, p.CurrentBet, p.TotalBet, g.currentBet, g.totalBet, false, nil
	} else {
		// 베팅 종료 조건 달성한 경우
		if g.getReadyPlayerIdx(currentPlayerIdx+1) == g.betLeaderIdx {
			g.currentPlayerIdx = g.firstPlayerIdx // 다음 베팅을 위해서 초기화
			g.clearPlayersCurrentBet()
			return nextPlayerName, false, p.CurrentBet, p.TotalBet, g.currentBet, g.totalBet, true, nil
		}

		// 베팅은 종료되지 않고 다음 플레이어가 베팅해야함
		g.currentPlayerIdx = nextPlayerIdx
		return nextPlayerName, false, p.CurrentBet, p.TotalBet, g.currentBet, g.totalBet, false, nil
	}
}

func (g Game) isValidBet(p *player.Player, betAmount uint64) (BetType, error) {
	if p.GameBalance == p.TotalBet+betAmount {
		return AllIn, nil
	}
	if p.GameBalance < p.TotalBet+betAmount {
		return -1, gameerror.OverBalance
	}
	if g.currentBet > p.CurrentBet+betAmount {
		return -1, gameerror.LowBetting
	}
	if g.currentBet < p.CurrentBet+betAmount {
		return Raise, nil
	}

	return Check, nil
}

func (g *Game) giveCardsToPlayers() {
	for i := 0; i < len(g.Players); i++ {
		g.Players[i].Hands = append(g.Players[i].Hands, g.deck.GetCard(), g.deck.GetCard())
	}
}

func (g Game) findPlayer(nickname string) (*player.Player, error) {
	for _, p := range g.Players {
		if p.Nickname == nickname {
			return p, nil
		}
	}
	return nil, gameerror.NoPlayerExists
}

func (g Game) getPlayerIdx(nickname string) (uint, error) {
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].Nickname == nickname {
			return uint(i), nil
		}
	}
	return 0, gameerror.NoPlayerExists
}

// 준비를 안해서 게임을 진행중이지 않거나 죽거나 나가는 사람이 있기 때문에 단순히 currentPlayerIdx를 1씩 증가하면 오류가 생김
func (g Game) getNextPlayerIdx() (uint, error) {
	idx := g.currentPlayerIdx
	for i := 0; i < len(g.Players); i++ {
		nextPlayer := g.Players[g.getReadyPlayerIdx(idx+1)]
		idx = g.getReadyPlayerIdx(idx + 1)

		if nextPlayer.IsReady && !nextPlayer.IsDead && !nextPlayer.IsLeft {
			break
		}
	}

	// 인덱스가 다시 돌아왔단 것은 모든 플레이어가 죽었거나 나갔다는 것
	if idx == g.currentPlayerIdx {
		return 0, gameerror.NoPlayersLeft
	}

	return idx, nil
}

// 함수 인자로 들어온 인덱스에 해당하는 플레이어가 Ready 상태면 해당 인덱스를 리턴하고
// 아니라면 다음 플레이어들 중에서 가장 빠른 순서인 Ready 상태인 플레이어에 해당 인덱스를 리턴
func (g Game) getReadyPlayerIdx(idx uint) uint {
	readyIdx := idx
	for i := 0; i < len(g.Players); i++ {
		if g.Players[readyIdx].IsReady {
			break
		}
		readyIdx = g.getNextIdx(readyIdx)
	}
	return readyIdx
}

// 현재 인덱스에 +1한 인덱스 값을 리턴
// 원형 리스트처럼 적용시켜야하므로 mod 연산 이용
func (g Game) getNextIdx(idx uint) uint {
	nextIdx := (idx + 1) % uint(len(g.Players))
	return nextIdx
}

func (g *Game) clearPlayersCurrentBet() {
	for _, p := range g.Players {
		p.CurrentBet = 0
	}
}

func (g *Game) SetBetResponseChan(betResponseChan chan channels.BetResponse) {
	g.BetResponseChan = betResponseChan
}
