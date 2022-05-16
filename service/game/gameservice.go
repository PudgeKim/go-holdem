package gameservice

import (
	"context"

	"github.com/PudgeKim/go-holdem/domain/entity"
	"github.com/PudgeKim/go-holdem/domain/repository"
	"github.com/PudgeKim/go-holdem/gameerror"
)
type BetType string 
const (
	CHECK = "Check"
	RAISE = "Raise"
	ALLIN = "AllIn"
)

const (
	FreeFlop = "FreeFlop"
	Flop     = "Flop"
	Turn     = "Turn"
	River    = "River"
	GameEnd  = "GameEnd"
)

// 게임이 처음 시작되는 경우
// smallBlind  = 0번째 인덱스에 해당하는 플레이어 (준비를 한 경우)
// bigBlind    = 1번째 인덱스에 해당하는 플레이어 (준비를 한 경우)
// firstPlayer = bigBlind 다음 인덱스에 해당하는 플레이어 (만약 플레이어가 2명이라면 smallBlind에 해당되는 플레이어)

type GameService struct {
	ctx context.Context
	userRepo repository.UserRepository
	gameRepo repository.GameRepository
}

func New(ctx context.Context, userRepo repository.UserRepository, gameRepo repository.GameRepository) *GameService {
	return &GameService{
		ctx: ctx,
		userRepo: userRepo,
		gameRepo: gameRepo,
	}
}

func (g *GameService) GetGame(ctx context.Context, roomId string) (*entity.Game, error) {
	return g.gameRepo.GetGame(ctx, roomId)
}

func (g *GameService) SaveGame(ctx context.Context, roomId string, game *entity.Game) error {
	return g.gameRepo.SaveGame(ctx, roomId, game)
}

func (g *GameService) CreateGame(ctx context.Context) (*entity.Game, error) {
	return g.gameRepo.CreateGame(ctx)
}

func (g *GameService) DeleteGame(ctx context.Context, roomId string) error {
	return g.gameRepo.DeleteGame(ctx, roomId)
}

func (g *GameService) FindPlayer(ctx context.Context, roomId string, nickname string) (*entity.Player, error) {
	return g.gameRepo.FindPlayer(ctx, roomId, nickname)
}

func (g *GameService) StartGame(ctx context.Context) (*GameStartResponse, error) {
	game, err := g.CreateGame(ctx); if err != nil {
		return nil, err 
	}

	// bigBlind, firstPlayer 등 세팅 
	readyPlayers, err := setPlayers(game); if err != nil {
		return nil, err 
	}

	// 카드 분배
	if err := giveCardsToPlayers(game); err != nil {
		return nil, err 
	}

	game.IsStarted = true 
	
	if err := g.SaveGame(ctx, game.RoomId.String(), game); err != nil {
		return nil, err 
	}

	firstPlayer := game.Players[game.FirstPlayerIdx].Nickname
	smallBlind := game.Players[game.SmallBlindIdx].Nickname
	bigBlind := game.Players[game.BigBlindIdx].Nickname
	
	
	gameStartResponse := NewGameStartResponse(readyPlayers, firstPlayer, smallBlind, bigBlind)
	return gameStartResponse, nil 
}

func (g *GameService) HandleReady(ctx context.Context, roomId string, nickname string, isReady bool) error {
	game, err := g.GetGame(ctx, roomId)
	if err != nil {
		return err
	}

	if game.IsStarted {
		return gameerror.GameAlreadyStarted
	}

	p, err := findPlayer(nickname, game); if err != nil {
		return err 
	}

	p.IsReady = isReady

	if err := g.SaveGame(ctx, roomId, game); err != nil {
		return err 
	}

	p.SetMemento()

	return nil
}

func (g *GameService) Bet(ctx context.Context, roomId string, betInfo BetInfo) (*BetResponse, error) {
	game, err := g.GetGame(ctx, roomId); if err != nil {
		return nil, err 
	}
	// 게임 준비 고려해서 코드 짜야됨

	var betResponse BetResponse

	nextPlayerName, isPlayerDead, playerCurBet, playerTotBet, gameCurBet, gameTotBet, isBetEnd, err := g.handleBet(ctx, roomId, game, betInfo)
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

	switch game.Status {
	case FreeFlop:
		if isBetEnd {
			game.Status = Flop
		}
		betResponse.GameStatus = Flop
	case Flop:
		if isBetEnd {
			game.Status = Turn
		}
		betResponse.GameStatus = Turn
	case Turn:
		if isBetEnd {
			game.Status = River
		}
		betResponse.GameStatus = River
	case River:
		if isBetEnd {
			// 게임 종료
			game.Status = GameEnd
			betResponse.GameStatus = GameEnd

			// 승자 계산과 승자와 패자 잔고 업데이트
			winners, losers := g.distributeMoneyToWinners(game)
			winnersAndLosers := append(winners, losers...)
			if err := g.updatePlayersBalance(winnersAndLosers...); err != nil {
				return nil, err 
			}
			
			var winnersName []string

			for _, p := range winners {
				winnersName = append(winnersName, p.Nickname)
			}

			betResponse.Winners = winnersName

			// 나간 플레이어들 빼줘야함
			removeLeftPlayers(game)
			if err != nil {
				return nil, err 
			}
			// 베팅 값들 모두 초기화
			clearPlayersCurrentBet(game.Players)

			game.CurrentBet = 0
			game.TotalBet = 0
			if err := g.SaveGame(ctx, game.RoomId.String(), game); err != nil {
				
				return nil, err 
			}	
		}
	}
	return &betResponse, nil 
}

// handleBet 모든 플레이어들의 베팅이 종료되는 경우면 true를 리턴함
// 다음 플레이어, 현재플레이어 isDead, 현재 플레이어의 currentBet, totalBet, 현재 게임의 currentBet, totalBet, 베팅종료, 에러 리턴
func (g *GameService) handleBet(ctx context.Context, roomId string, game *entity.Game, betInfo BetInfo) (string, bool, uint64, uint64, uint64, uint64, bool, error) {
	p, err := findPlayer(betInfo.PlayerName, game)
	if err != nil {
		return "", false, 0, 0, 0, 0, false, err
	}

	expectedPlayer := game.Players[game.CurrentPlayerIdx]
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
		nextPlayerIdx, err := getNextPlayerIdx(game)
		if err != nil {
			p.Undo()
			return "", false, 0, 0, 0, 0, false, err
		}
		nextPlayer := game.Players[nextPlayerIdx].Nickname
		return nextPlayer, true, p.CurrentBet, p.TotalBet, game.CurrentBet, game.TotalBet, false, nil
	}

	betType := getBetType(p, betInfo.BetAmount, game)
	if err != nil {
		return "", false, 0, 0, 0, 0, false, err
	}
	if betType == ALLIN {
		p.IsAllIn = true
	}

	p.CurrentBet += betInfo.BetAmount
	p.TotalBet += betInfo.BetAmount

	nextPlayerIdx, err := getNextPlayerIdx(game)
	if err != nil {
		p.Undo()
		return "", false, 0, 0, 0, 0, false, err
	}

	nextPlayerName := game.Players[nextPlayerIdx].Nickname
	currentPlayerIdx, err := getPlayerIdx(game.Players, p.Nickname)
	if err != nil {
		p.Undo()
		return "", false, 0, 0, 0, 0, false, err
	}

	// 현재 베팅한 플레이어가 베팅한 금액에 따라 베팅리더인지 체크 후에 현재 베팅 턴을 종료할지 결정
	// (현재 플레이어가 베팅리더가 아니고, betLeader 이전 플레이어면 플레이어들의 베팅이 종료됨)
	// 베팅이 종료되면 다음 베팅을 위해서 player들의 currentBet을 초기화시켜주어야함
	if betInfo.BetAmount > game.CurrentBet { // 현재 플레이어가 베팅 리더가 되는 경우
		game.CurrentBet = p.CurrentBet
		game.BetLeaderIdx = currentPlayerIdx
		game.CurrentPlayerIdx = nextPlayerIdx
		return nextPlayerName, false, p.CurrentBet, p.TotalBet, game.CurrentBet, game.TotalBet, false, nil
	} else {
		// 베팅 종료 조건 달성한 경우
		if getReadyPlayerIdx(game.Players, currentPlayerIdx+1) == game.BetLeaderIdx {
			game.CurrentPlayerIdx = game.FirstPlayerIdx // 다음 베팅을 위해서 초기화
			clearPlayersCurrentBet(game.Players)
			if err := g.SaveGame(ctx, game.RoomId.String(), game); err != nil {
				p.Undo()
				game.Undo()
			}
			return nextPlayerName, false, p.CurrentBet, p.TotalBet, game.CurrentBet, game.TotalBet, true, nil
		}

		// 베팅은 종료되지 않고 다음 플레이어가 베팅해야함
		game.CurrentPlayerIdx = nextPlayerIdx

		if err := g.SaveGame(ctx, game.RoomId.String(), game); err != nil {
			p.Undo()
			game.Undo()
			return "", false, 0, 0, 0, 0, false, err
		}

		return nextPlayerName, false, p.CurrentBet, p.TotalBet, game.CurrentBet, game.TotalBet, false, nil
	}
}

func (g *GameService) distributeMoneyToWinners(game *entity.Game) (winners []*entity.Player, losers []*entity.Player) {
	winners, losers, err := getWinnersAndLosers(game); if err != nil {
		return nil, nil
	}

	if len(winners) == 1 {
		winner := winners[0]
		winner.GameBalance += game.TotalBet - winner.TotalBet
		winner.TotalBalance += game.TotalBet - winner.TotalBet
		
		for _, loser := range losers {
			loser.GameBalance -= loser.TotalBet
			loser.TotalBalance -= loser.TotalBet
		}

		return winners, losers
	}

	// 승자가 여러명인 경우 
	// 단순히 1/N 하면 안되고 중간에 올인여부를 판단해야함 

	var tmpWinners []*entity.Player
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

	return winners, losers
}

func (g *GameService) updatePlayersBalance(players ...*entity.Player) error {
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

	// DB 저장 성공시 메멘토 업데이트
	for _, p := range players {
		p.SetMemento()
	}

	return nil 
}