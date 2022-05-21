package service

import (
	"github.com/PudgeKim/go-holdem/domain/entity"
	"github.com/PudgeKim/go-holdem/gameerror"
)

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

// 현재 인덱스에 +1한 인덱스 값을 리턴
// 원형 리스트처럼 적용시켜야하므로 mod 연산 이용
func getNextIdx(players []*entity.Player, idx uint) uint {
	nextIdx := (idx + 1) % uint(len(players))
	return nextIdx
}

func getPlayerIdx(players []*entity.Player, nickname string) (uint, error) {
	for i := 0; i < len(players); i++ {
		if players[i].Nickname == nickname {
			return uint(i), nil
		}
	}
	return 0, gameerror.NoPlayerExists
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