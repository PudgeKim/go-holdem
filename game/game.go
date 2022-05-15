package game

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/PudgeKim/go-holdem/card"
	"github.com/PudgeKim/go-holdem/domain/repository"
	"github.com/PudgeKim/go-holdem/gameerror"
	"github.com/PudgeKim/go-holdem/player"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)



type CardCompareResult int
const (
	RoomLimit = 7
	RedisTimeDuration = time.Hour * 144
)
const (
	Player1Win CardCompareResult = iota
	Player2Win
	Draw
)

// 게임이 처음 시작되는 경우
// smallBlind  = 0번째 인덱스에 해당하는 플레이어 (준비를 한 경우)
// bigBlind    = 1번째 인덱스에 해당하는 플레이어 (준비를 한 경우)
// firstPlayer = bigBlind 다음 인덱스에 해당하는 플레이어 (만약 플레이어가 2명이라면 smallBlind에 해당되는 플레이어)

type Game struct {
	memento *GameMemento

	userRepo repository.UserRepository
	redisClient *redis.Client // redis같은 캐시서버에 현재 방에 있는 플레이어들을 저장해둠
	ctx context.Context

	// 아래 필드들은 redis에 저장됨 
	RoomId uuid.UUID
	RoomLimit uint 
	Players    []*player.Player // 게임에 참가하고 있는 플레이어들
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

func New(ctx context.Context, userRepo repository.UserRepository, redisClient *redis.Client, roomId uuid.UUID) *Game {
	game := &Game{
		userRepo: userRepo,
		redisClient: redisClient,
		ctx: ctx,
		RoomId: roomId,
		RoomLimit: RoomLimit,
		Players:    make([]*player.Player, 2),
		TotalBet:   0,
		CurrentBet: 0,
		IsStarted:  false,
		Deck:       card.NewDeck(),
		Status:     FreeFlop,
	}
	memento := NewGameMemento(*game)
	game.memento = memento

	return game
}


// Start 게임 시작 요청이 들어오면 해당 함수 실행
func (g *Game) Start() (*GameStartResponse, error) {

	// bigBlind, firstPlayer 등 세팅 
	readyPlayers, err := g.setPlayers(); if err != nil {
		g.Undo()
		return nil, err 
	}

	// 카드 분배
	if err := g.giveCardsToPlayers(); err != nil {
		g.Undo()
		return nil, err 
	}

	g.IsStarted = true 
	
	if err := g.setRedis(); err != nil {
		g.Undo()
		return nil, err 
	}

	firstPlayer := g.Players[g.FirstPlayerIdx].Nickname
	smallBlind := g.Players[g.SmallBlindIdx].Nickname
	bigBlind := g.Players[g.BigBlindIdx].Nickname
	
	
	gameStartResponse := NewGameStartResponse(readyPlayers, firstPlayer, smallBlind, bigBlind)
	return gameStartResponse, nil 
}

// SetPlayers 준비된 플레이어들의 순서와 smallBlind, bigBlind를 지정해줌
// 게임이 처음 시작됬는지 아닌지에 따라 구별함
func (g *Game) setPlayers() ([]string, error) {
	
	if len(g.Players) < 2 {
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
			g.SmallBlindIdx = g.getReadyPlayerIdx(0)
			g.BigBlindIdx = g.getReadyPlayerIdx(g.SmallBlindIdx + 1)
			g.FirstPlayerIdx = g.getReadyPlayerIdx(g.BigBlindIdx + 1)
		} else { // 플레이어가 2명 밖에 없는 경우 (이 함수의 첫 부분 검사에 의해 최소 2명은 Ready 상태이므로 0번과 1번 인덱스가 모두 Ready 상태임)
			g.SmallBlindIdx = 0
			g.BigBlindIdx = 1
			g.FirstPlayerIdx = g.SmallBlindIdx
		}
	} else { // 기존에 진행되던 순서가 있는 경우
		if len(g.Players) > 2 {
			g.SmallBlindIdx = g.getReadyPlayerIdx(g.SmallBlindIdx + 1)
			g.BigBlindIdx = g.getReadyPlayerIdx(g.SmallBlindIdx + 1)
			g.FirstPlayerIdx = g.getReadyPlayerIdx(g.BigBlindIdx + 1)
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

func (g *Game) Bet(betInfo BetInfo) (*BetResponse, error) {
	// 게임 준비 고려해서 코드 짜야됨

	var betResponse BetResponse

	nextPlayerName, isPlayerDead, playerCurBet, playerTotBet, gameCurBet, gameTotBet, isBetEnd, err := g.handleBet(betInfo)
	if err != nil {
		return nil, err 
	}

	betResponse.IsBetEnd = isBetEnd
	betResponse.IsPlayerDead = isPlayerDead
	betResponse.PlayerCurrentBet = playerCurBet
	betResponse.PlayerTotalBet = playerTotBet
	betResponse.GameCurrentBet = gameCurBet
	betResponse.GameTotalBet = gameTotBet
	betResponse.NextPlayerName = nextPlayerName

	switch g.Status {
	case FreeFlop:
		if isBetEnd {
			g.Status = Flop
		}
		betResponse.GameStatus = Flop
	case Flop:
		if isBetEnd {
			g.Status = Turn
		}
		betResponse.GameStatus = Turn
	case Turn:
		if isBetEnd {
			g.Status = River
		}
		betResponse.GameStatus = River
	case River:
		if isBetEnd {
			// 게임 종료
			g.Status = GameEnd
			betResponse.GameStatus = GameEnd

			// 승자 계산과 승자와 패자 잔고 DB에 저장
			winners, err := g.DistributeMoneyToWinners()
			if err != nil {
				betResponse.Error = err
			}
			
			var winnersName []string

			for _, p := range winners {
				winnersName = append(winnersName, p.Nickname)
			}

			betResponse.Winners = winnersName

			// 나간 플레이어들 빼줘야함
			err = g.removeLeftPlayers()
			if err != nil {
				betResponse.Error = err 
			}
			// 베팅 값들 모두 초기화
			err = g.clearPlayersCurrentBet()
			if err != nil {
				betResponse.Error = err 
			}

			g.CurrentBet = 0
			g.TotalBet = 0
			err = g.setRedis()
			if err != nil {
				betResponse.Error = err 
			}			
		}
	}
	return &betResponse, nil 
}

// handleBet 모든 플레이어들의 베팅이 종료되는 경우면 true를 리턴함
// 다음 플레이어, 현재플레이어 isDead, 현재 플레이어의 currentBet, totalBet, 현재 게임의 currentBet, totalBet, 베팅종료 리턴
func (g *Game) handleBet(betInfo BetInfo) (string, bool, uint64, uint64, uint64, uint64, bool, error) {
	p, err := g.FindPlayer(betInfo.PlayerName)
	if err != nil {
		return "", false, 0, 0, 0, 0, false, err
	}

	expectedPlayer := g.Players[g.CurrentPlayerIdx]
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
			p.Undo()
			return "", false, 0, 0, 0, 0, false, err
		}
		nextPlayer := g.Players[nextPlayerIdx].Nickname
		return nextPlayer, true, p.CurrentBet, p.TotalBet, g.CurrentBet, g.TotalBet, false, nil
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
		p.Undo()
		return "", false, 0, 0, 0, 0, false, err
	}

	nextPlayerName := g.Players[nextPlayerIdx].Nickname
	currentPlayerIdx, err := g.getPlayerIdx(p.Nickname)
	if err != nil {
		p.Undo()
		return "", false, 0, 0, 0, 0, false, err
	}

	// 현재 베팅한 플레이어가 베팅한 금액에 따라 베팅리더인지 체크 후에 현재 베팅 턴을 종료할지 결정
	// (현재 플레이어가 베팅리더가 아니고, betLeader 이전 플레이어면 플레이어들의 베팅이 종료됨)
	// 베팅이 종료되면 다음 베팅을 위해서 player들의 currentBet을 초기화시켜주어야함
	if betInfo.BetAmount > g.CurrentBet { // 현재 플레이어가 베팅 리더가 되는 경우
		g.CurrentBet = p.CurrentBet
		g.BetLeaderIdx = currentPlayerIdx
		g.CurrentPlayerIdx = nextPlayerIdx
		return nextPlayerName, false, p.CurrentBet, p.TotalBet, g.CurrentBet, g.TotalBet, false, nil
	} else {
		// 베팅 종료 조건 달성한 경우
		if g.getReadyPlayerIdx(currentPlayerIdx+1) == g.BetLeaderIdx {
			g.CurrentPlayerIdx = g.FirstPlayerIdx // 다음 베팅을 위해서 초기화
			if err := g.clearPlayersCurrentBet(); err != nil {
				p.Undo()
				g.Undo()
			}
			return nextPlayerName, false, p.CurrentBet, p.TotalBet, g.CurrentBet, g.TotalBet, true, nil
		}

		// 베팅은 종료되지 않고 다음 플레이어가 베팅해야함
		g.CurrentPlayerIdx = nextPlayerIdx

		if err := g.setRedis(); err != nil {
			p.Undo()
			g.Undo()
			return "", false, 0, 0, 0, 0, false, err
		}

		return nextPlayerName, false, p.CurrentBet, p.TotalBet, g.CurrentBet, g.TotalBet, false, nil
	}
}

func (g Game) isValidBet(p *player.Player, betAmount uint64) (BetType, error) {
	if p.GameBalance == p.TotalBet+betAmount {
		return AllIn, nil
	}
	if p.GameBalance < p.TotalBet+betAmount {
		return -1, gameerror.OverBalance
	}
	if g.CurrentBet > p.CurrentBet+betAmount {
		return -1, gameerror.LowBetting
	}
	if g.CurrentBet < p.CurrentBet+betAmount {
		return Raise, nil
	}

	return Check, nil
}

func (g *Game) giveCardsToPlayers() error {
	validPlayers, err := g.getValidPlayers()
	if err != nil {
		return err 
	}

	for i := 0; i < len(validPlayers); i++ {
		g.Players[i].Hands = append(g.Players[i].Hands, g.Deck.GetCard(), g.Deck.GetCard())
	}

	if err := g.setRedis(); err != nil {
		// 에러 발생시 롤백 
		for _, p := range validPlayers {
			p.Undo()
		}
		g.Undo()
		return err 
	}

	return nil 
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
	idx := g.CurrentPlayerIdx
	for i := 0; i < len(g.Players); i++ {
		nextPlayer := g.Players[g.getReadyPlayerIdx(idx+1)]
		idx = g.getReadyPlayerIdx(idx + 1)

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

func (g *Game) clearPlayersCurrentBet() error {
	for _, p := range g.Players {
		p.CurrentBet = 0
	}

	if err := g.setRedis(); err != nil {
		return err 
	}
	return nil
}

// 현재 플레이어들 중 나가지도 않고 죽지도 않고 준비도 된 플레이어들만 리턴
func (g *Game) getValidPlayers() ([]*player.Player, error) {
	if err := g.getRedis(); err != nil {
		return nil, err 
	}

	var validPlayers []*player.Player

	for _, p := range g.Players {
		if p.IsReady && !p.IsDead && !p.IsLeft {
			validPlayers = append(validPlayers, p)
		}
	}

	return validPlayers, nil
}

func (g *Game) DistributeMoneyToWinners() ([]*player.Player, error) {
	winners, losers, err := g.getWinnersAndLosers(); if err != nil {
		return nil, err 
	}

	if len(winners) == 1 {
		winner := winners[0]
		winner.GameBalance += g.TotalBet - winner.TotalBet
		winner.TotalBalance += g.TotalBet - winner.TotalBet
		
		for _, loser := range losers {
			loser.GameBalance -= loser.TotalBet
			loser.TotalBalance -= loser.TotalBet
		}

		// DB에 저장 
		if err := g.updatePlayersBalance(); err != nil {
			// 저장 실패시 롤백
			winner.Undo()
			for _, loser := range losers {
				loser.Undo()
			}
			return nil, err 
		}
		
		// DB 저장 성공시 메멘토 업데이트
		// DB에 저장까지 잘 성공했으면 메멘토 업데이트 
		winner.SetMemento()
		for _, loser := range losers {
			loser.SetMemento()
		}

		return winners, nil
	}

	// 승자가 여러명인 경우 
	// 단순히 1/N 하면 안되고 중간에 올인여부를 판단해야함 

	var tmpWinners []*player.Player
	// 첫번째 winner를 우선 추가해둠 
	tmpWinners = append(tmpWinners, winners[0])

	for i:=1; i<len(winners); i++ {
		prevWinner, curWinner := winners[i-1], winners[i]

		// winners는 totalBet을 기준으로 오름차순 정렬되어있으므로
		// prevWinner와 curWinner의 totalBet이 다르다면 prevWinner는 중간에 올인을 한 것이고 curWinner는 더 베팅을 한 플레이어임
		if prevWinner.TotalBet != curWinner.TotalBet {
			winnerReward := uint64((len(losers) * int(prevWinner.TotalBet)) / len(tmpWinners))

			for _, winner := range tmpWinners {
				winner.GameBalance += winnerReward
				winner.TotalBalance += winnerReward
			}

			// 올인한사람에게 돈을 줬으므로 나머지 winner에서 돈 빼줘야함 
			for j:=i; j<len(winners); j++ {
				winners[i].GameBalance -= winnerReward
				winners[i].TotalBalance -= winnerReward
			}

			tmpWinners = nil 
		} else {
			tmpWinners = append(tmpWinners, curWinner)
		}
	}

	if len(tmpWinners) > 0 {
		for _, winner := range tmpWinners {
			winnerReward := uint64((len(losers) * int(winner.TotalBet)) / len(tmpWinners))
			winner.GameBalance += winnerReward
			winner.TotalBalance += winnerReward
		}
	}

	for _, loser := range losers {
		loser.GameBalance -= loser.TotalBet
		loser.TotalBalance -= loser.TotalBet
	}

	// winners와 losers들의 잔고 DB에 업데이트
	winnersAndLosers := append(winners, losers...)
	err = g.updatePlayersBalance(winnersAndLosers...); if err != nil {
		// DB업데이트 실패했으면 롤백 
		for _, winner := range winners {
			winner.Undo()
		}
		for _, loser := range losers {
			loser.Undo()
		}

		return nil, err 
	}
	
	// DB에 저장까지 잘 성공했으면 메멘토 업데이트 
	for _, winner := range winners {
		winner.SetMemento()
	}
	for _, loser := range losers {
		loser.SetMemento()
	}

	return winners, nil 
}

func (g *Game) updatePlayersBalance(players ...*player.Player) error {
	var userIdWithBalances []repository.UserIdWithBalance

	for _, p := range players {
		userIdWithBalances = append(userIdWithBalances, repository.NewUserIdWithBalance(p.Id, p.TotalBalance))
	}

	if err := g.userRepo.UpdateMultipleBalance(g.ctx, userIdWithBalances); err != nil {
		// 실패시 롤백
		for _, p := range players {
			p.Undo()
		}
		return err 
	}

	// 디비 저장 성공시 메멘토 업데이트
	for _, p := range players {
		p.SetMemento()
	}

	return nil 
}

// 리턴되는 배열에는 1명 이상의 승리자가 들어가게 되는데 
// 배열에 포함되는 승리자들이 2명 이상인 경우 단순히 1/n으로 나누면 안됨 
// (돈이 상대적으로 없는 사용자가 올인하고 나머지 플레이어들은 추가 베팅이 가능하기 때문)
func (g *Game) getWinnersAndLosers() (winners []*player.Player, losers []*player.Player, err error) {
	// 죽거나 나가지 않은 플레이어들로 승리자 계산해야함 
	validPlayers, err := g.getValidPlayers()
	if err != nil {
		return nil, nil, err 
	}

	if len(validPlayers) == 0 {
		return nil, nil, errors.New("zero player")
	}

	if len(validPlayers) == 1 {
		return []*player.Player{validPlayers[0]}, nil, nil 
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
	for _, p := range validPlayers {
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

// 두 플레이어 간의 bestCards를 비교해서 이긴 플레이어를 리턴
// 둘이 같다면 Draw를 리턴
// ** 각 플레이어들의 bestCards는 정렬되어 있음 (bestCards를 만드는 과정에서 정렬 함수가 쓰임)
func compare(player1 *player.Player, player2 *player.Player) CardCompareResult {
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

func (g *Game) HandleReady(nickname string, isReady bool) error {
	p, err := g.FindPlayer(nickname)
	if err != nil {
		return err
	}

	if g.IsStarted {
		return gameerror.GameAlreadyStarted
	}

	p.IsReady = isReady

	if err = g.setRedis(); err != nil {
		p.Undo()
		return err 
	}

	p.SetMemento()

	return nil
}

func (g *Game) FindPlayer(nickname string) (*player.Player, error) {
	if err := g.getRedis(); err != nil {
		return nil, err 
	}

	for _, p := range g.Players {
		if p.Nickname == nickname {
			return p, nil
		}
	}
	return nil, gameerror.NoPlayerExists
}

func (g *Game) AddPlayer(player *player.Player) error {
	_, err := g.FindPlayer(player.Nickname)
	// nil이면 플레이어가 이미 존재하는데 또 요청이 온 것
	if err == nil {
		return gameerror.PlayerAlreadyExists
	}

	if len(g.Players) > RoomLimit {
		return gameerror.PlayerLimitationError
	}

	g.Players = append(g.Players, player)

	if err := g.setRedis(); err != nil {
		g.Undo()
		return err 
	}

	g.SetMemento()

	return nil
}

// 나간 플레이어들 고려해서 인덱스 변경해야함
// smallBlind와 bigBlind가 바뀌었는지 리턴
func (g *Game) removeLeftPlayers() error {
	if err := g.getRedis(); err != nil {
		return err 
	}

	var indexesToRemove []int

	for idx, p := range g.Players {
		if p.IsLeft {
			if g.SmallBlindIdx == uint(idx) {
				g.SmallBlindIdx = g.getNextIdx(uint(idx))
			}
			if g.BigBlindIdx == uint(idx) {
				g.BigBlindIdx = g.getNextIdx(uint(idx))
			}

			indexesToRemove = append(indexesToRemove, idx)
		}
	}

	for _, idx := range indexesToRemove {
		g.Players = removePlayerByIndex(g.Players, idx)
	}

	if err := g.setRedis(); err != nil {
		return err 
	}

	return nil 
}

func removePlayerByIndex(players []*player.Player, s int) []*player.Player {
	return append(players[:s], players[s+1:]...)
}


func (g Game) setRedis() error {
	statusCmd := g.redisClient.Set(g.ctx, g.RoomId.String(), g, RedisTimeDuration)
	if statusCmd.Err() != nil {
		// TODO: 여기서 에러 발생시 롤백해야함 
		return statusCmd.Err()
	}
	return nil 
}

func (g *Game) getRedis() error {
	stringCmd := g.redisClient.Get(g.ctx, g.RoomId.String())
	if stringCmd.Err() != nil {
		return stringCmd.Err()
	}

	var gameFromRedis Game 

	if err := stringCmd.Scan(&gameFromRedis); err != nil {
		return err 
	}

	g.RoomLimit = gameFromRedis.RoomLimit
	g.Players = gameFromRedis.Players
	g.TotalBet = gameFromRedis.TotalBet
	g.CurrentBet = gameFromRedis.CurrentBet
	g.IsStarted = gameFromRedis.IsStarted

	g.Deck = gameFromRedis.Deck
	g.Status = gameFromRedis.Status
	
	g.SmallBlindIdx = gameFromRedis.SmallBlindIdx
	g.BigBlindIdx = gameFromRedis.BigBlindIdx

	g.FirstPlayerIdx = gameFromRedis.FirstPlayerIdx
	g.IsFirstPlayerBet = gameFromRedis.IsFirstPlayerBet
	g.CurrentPlayerIdx = gameFromRedis.CurrentPlayerIdx
	g.BetLeaderIdx = gameFromRedis.BetLeaderIdx
	
	return nil 
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

func (g *Game) Undo() {
	memento := g.memento

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
	memento := g.memento

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