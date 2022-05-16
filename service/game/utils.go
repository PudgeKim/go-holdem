package gameservice

import (
	"errors"
	"sort"

	"github.com/PudgeKim/go-holdem/domain/entity"
	"github.com/PudgeKim/go-holdem/gameerror"
)
type CardCompareResult int

const (
	Player1Win CardCompareResult = iota
	Player2Win
	Draw
)

func giveCardsToPlayers(game *entity.Game) error {
	validPlayers := getValidPlayers(game.Players)

	for i := 0; i < len(validPlayers); i++ {
		game.Players[i].Hands = append(game.Players[i].Hands, game.Deck.GetCard(), game.Deck.GetCard())
	}

	return nil 
}

// 현재 플레이어들 중 나가지도 않고 죽지도 않고 준비도 된 플레이어들만 리턴
func getValidPlayers(players []*entity.Player) []*entity.Player {
	var validPlayers []*entity.Player

	for _, p := range players {
		if p.IsReady && !p.IsDead && !p.IsLeft {
			validPlayers = append(validPlayers, p)
		}
	}

	return validPlayers
}

func getReadyPlayers(players []*entity.Player) []*entity.Player {
	var readyPlayers []*entity.Player

	for _, p := range players {
		if p.IsReady {
			readyPlayers = append(readyPlayers, p)
		}
	}

	return readyPlayers
}

// SetPlayers 준비된 플레이어들의 순서와 smallBlind, bigBlind를 지정해줌
// 게임이 처음 시작됬는지 아닌지에 따라 구별함
func setPlayers(game *entity.Game) ([]string, error) {
	
	if len(game.Players) < 2 {
		return nil, gameerror.LackOfPlayers
	}

	var readyPlayersName []string 
	readyCnt := 0
	for i := 0; i < len(game.Players); i++ {
		if game.Players[i].IsReady {
			readyPlayersName = append(readyPlayersName, game.Players[i].Nickname)
			readyCnt += 1
		}
	}
	if readyCnt < 2 {
		return nil, gameerror.NotEnoughPlayersReady
	}

	// 게임이 처음 시작되는 경우와 이미 몇판 진행되고 있는지에 따라 나눔
	// (이미 진행되고 있었다면 이전 게임의 플레이어들의 순서를 기준으로 세팅해야되기 때문)
	var isFirstGame bool

	if game.SmallBlindIdx == 0 && game.BigBlindIdx == 0 {
		isFirstGame = true
	} else {
		isFirstGame = false
	}

	if isFirstGame {
		if len(game.Players) > 2 {
			game.SmallBlindIdx = getReadyPlayerIdx(game.Players, 0)
			game.BigBlindIdx = getReadyPlayerIdx(game.Players, game.SmallBlindIdx + 1)
			game.FirstPlayerIdx = getReadyPlayerIdx(game.Players, game.BigBlindIdx + 1)
		} else { // 플레이어가 2명 밖에 없는 경우 (이 함수의 첫 부분 검사에 의해 최소 2명은 Ready 상태이므로 0번과 1번 인덱스가 모두 Ready 상태임)
			game.SmallBlindIdx = 0
			game.BigBlindIdx = 1
			game.FirstPlayerIdx = game.SmallBlindIdx
		}
	} else { // 기존에 진행되던 순서가 있는 경우
		if len(game.Players) > 2 {
			game.SmallBlindIdx = getReadyPlayerIdx(game.Players, game.SmallBlindIdx + 1)
			game.BigBlindIdx = getReadyPlayerIdx(game.Players, game.SmallBlindIdx + 1)
			game.FirstPlayerIdx = getReadyPlayerIdx(game.Players, game.BigBlindIdx + 1)
		} else { // 플레이어가 2명만 있는 경우
			game.SmallBlindIdx, game.BigBlindIdx = game.BigBlindIdx, game.SmallBlindIdx
			game.FirstPlayerIdx = game.SmallBlindIdx
		}
	}

	game.CurrentPlayerIdx = game.FirstPlayerIdx
	game.BetLeaderIdx = game.FirstPlayerIdx
	game.IsFirstPlayerBet = false

	return readyPlayersName, nil
}

func findPlayer(nickname string, game *entity.Game) (*entity.Player, error) {
	for _, p := range game.Players {
		if p.Nickname == nickname {
			return p, nil 
		}
	}

	return nil, gameerror.NoPlayerExists
}

// 나간 플레이어들 고려해서 인덱스 변경해야함
// smallBlind와 bigBlind가 바뀌었는지 리턴
func removeLeftPlayers(game *entity.Game) {
	var indexesToRemove []int

	for idx, p := range game.Players {
		if p.IsLeft {
			if game.SmallBlindIdx == uint(idx) {
				game.SmallBlindIdx = getNextIdx(game.Players, uint(idx))
			}
			if game.BigBlindIdx == uint(idx) {
				game.BigBlindIdx = getNextIdx(game.Players, uint(idx))
			}

			indexesToRemove = append(indexesToRemove, idx)
		}
	}

	for _, idx := range indexesToRemove {
		game.Players = removePlayerByIndex(game.Players, idx)
	}

}

func getPlayerIdx(players []*entity.Player, nickname string) (uint, error) {
	for i := 0; i < len(players); i++ {
		if players[i].Nickname == nickname {
			return uint(i), nil
		}
	}
	return 0, gameerror.NoPlayerExists
}

// 함수 인자로 들어온 인덱스에 해당하는 플레이어가 Ready 상태면 해당 인덱스를 리턴하고
// 아니라면 다음 플레이어들 중에서 가장 빠른 순서인 Ready 상태인 플레이어에 해당 인덱스를 리턴
func getReadyPlayerIdx(players []*entity.Player, idx uint) uint {
	readyIdx := idx
	for i := 0; i < len(players); i++ {
		if players[readyIdx].IsReady {
			break
		}
		readyIdx = getNextIdx(players, readyIdx)
	}
	return readyIdx
}

func removePlayerByIndex(players []*entity.Player, s int) []*entity.Player {
	return append(players[:s], players[s+1:]...)
}

// 현재 인덱스에 +1한 인덱스 값을 리턴
// 원형 리스트처럼 적용시켜야하므로 mod 연산 이용
func getNextIdx(players []*entity.Player, idx uint) uint {
	nextIdx := (idx + 1) % uint(len(players))
	return nextIdx
}

// 준비를 안해서 게임을 진행중이지 않거나 죽거나 나가는 사람이 있기 때문에 단순히 currentPlayerIdx를 1씩 증가하면 오류가 생김
func getNextPlayerIdx(game *entity.Game) (uint, error) {
	idx := game.CurrentPlayerIdx
	for i := 0; i < len(game.Players); i++ {
		nextPlayer := game.Players[getReadyPlayerIdx(game.Players, idx+1)]
		idx = getReadyPlayerIdx(game.Players, idx + 1)

		if nextPlayer.IsReady && !nextPlayer.IsDead && !nextPlayer.IsLeft {
			break
		}
	}

	// 인덱스가 다시 돌아왔단 것은 모든 플레이어가 죽었거나 나갔다는 것
	if idx == game.CurrentPlayerIdx {
		return 0, gameerror.NoPlayersLeft
	}

	return idx, nil
}

func getBetType(p *entity.Player, betAmount uint64, game *entity.Game) BetType {
	if p.GameBalance == p.TotalBet+betAmount {
		return ALLIN
	}
	if game.CurrentBet < p.CurrentBet + betAmount {
		return RAISE
	}
	return CHECK
}

func clearPlayersCurrentBet(players []*entity.Player) {
	for _, p := range players {
		p.CurrentBet = 0
	}
}

// 리턴되는 배열에는 1명 이상의 승리자가 들어가게 되는데 
// 배열에 포함되는 승리자들이 2명 이상인 경우 단순히 1/n으로 나누면 안됨 
// (돈이 상대적으로 없는 사용자가 올인하고 나머지 플레이어들은 추가 베팅이 가능하기 때문)
func getWinnersAndLosers(game *entity.Game) (winners []*entity.Player, losers []*entity.Player, err error) {
	// 죽거나 나가지 않은 플레이어들로 승리자 계산해야함 
	validPlayers := getValidPlayers(game.Players)

	if len(validPlayers) == 0 {
		return nil, nil, errors.New("zero player")
	}

	if len(validPlayers) == 1 {
		return []*entity.Player{validPlayers[0]}, nil, nil 
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
	for _, p := range getReadyPlayers(game.Players) {
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
func compare(player1 *entity.Player, player2 *entity.Player) CardCompareResult {
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