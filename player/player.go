package player

import (
	"github.com/PudgeKim/go-holdem/card"
)

type Player struct {
	Id 			 int64 // User struct의 id
	Nickname     string
	IsReady      bool // 게임준비
	IsDead       bool
	IsLeft       bool // 게임 중간에 나간 경우 여기에 우선 체크를 해두고 게임이 종료되면 실제로 나가게 처리함 (인덱스가 꼬이는거 방지하기 위해)
	IsAllIn      bool
	PrevTotalBalance uint64     // 롤백+이전 금액 확인용 
	TotalBalance uint64         // 매 게임 또는 플레이어가 죽거나 나가는 경우 갱신
	PrevGameBalance uint64      // 롤백+이전 금액 확인용 
	GameBalance  uint64         // 게임 참가시에 들고갈 돈 (매 게임 또는 플레이어가 죽거나 나가는 경우 갱신)
	TotalBet     uint64         // 해당 게임에서 누적 베팅액
	CurrentBet   uint64         // 현재 턴에서 베팅한 금액
	Hands        []card.Card    // 처음 받는 2장의 카드
	HandsRank    card.HandsRank // 족보 (fullHouse인지 onePair인지.. 등)
	HighCard     card.Rank      // 예를 들어 33322 fullHouse면 highCard는 3
	BestCards    []card.Card    // 필드에 카드가 모두 오픈되었을 때 hands까지 합쳐서 가장 좋은 5장의 카드들
}

func New(id int64, nickname string, totalBalance, gameBalance uint64) Player {
	return Player{
		Id: 			id,
		Nickname:     nickname,
		IsReady:      false,
		IsDead:       false,
		IsLeft:       false,
		IsAllIn:      false,
		PrevTotalBalance: totalBalance,
		TotalBalance: totalBalance,
		PrevGameBalance: gameBalance,
		GameBalance:  gameBalance,
		TotalBet:     0,
		CurrentBet:   0,
	}
}


