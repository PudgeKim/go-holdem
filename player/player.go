package player

import "github.com/PudgeKim/card"

type Player struct {
	nickname string
	isDead bool
	totalBalance uint64
	gameBalance uint64 // 게임 참가시에 들고갈 돈
	cards []card.Card // 처음 받는 2장의 카드
	handsRank card.HandsRank
}

func New(nickname string, totalBalance, gameBalance uint64) Player {
	return Player{
		nickname: nickname,
		isDead: false,
		totalBalance: totalBalance,
		gameBalance: gameBalance,
	}
}