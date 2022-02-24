package player

import "github.com/PudgeKim/card"

type Player struct {
	nickname     string
	isDead       bool
	totalBalance uint64
	gameBalance  uint64         // 게임 참가시에 들고갈 돈
	hands        []card.Card    // 처음 받는 2장의 카드
	handsRank    card.HandsRank // 족보 (fullHouse인지 onePair인지.. 등)
	highCard     card.Rank      // 예를 들어 33322 fullHouse면 highCard는 3
	bestCards    []card.Card    // 필드에 카드가 모두 오픈되었을 때 hands까지 합쳐서 가장 좋은 5장의 카드들
}

func New(nickname string, totalBalance, gameBalance uint64) Player {
	return Player{
		nickname:     nickname,
		isDead:       false,
		totalBalance: totalBalance,
		gameBalance:  gameBalance,
	}
}

func GetWinners(players []Player) []Player {
	if len(players) < 1 {
		panic("zero players")
	}

	// 승리자는 여러 명이 나올 수 있으므로 배열을 이용함
	var winners []Player

	// 첫번째 플레이어를 우선 승리자로 지정해놓고 이후 플레이어들과 비교해나감
	winners = append(winners, players[0])

	for i := 1; i < len(players); i++ {
		winner := winners[0]
		player := players[i]

		switch compare(winner, player) {
		case Player1Win:
			continue
		case Player2Win:
			// 배열에 기존 승리자가 여러명 있을 수도 있으니 아예 비워야함
			winners = []Player{}
			winners = append(winners, player)
		case Draw:
			// 기존 승리자들과 같은 패이므로 승리자에 추가
			winners = append(winners, player)
		}
	}

	return winners
}

// 두 플레이어 간의 bestCards를 비교해서 이긴 플레이어를 리턴
// 둘이 같다면 Draw를 리턴
// ** 각 플레이어들의 bestCards는 정렬되어 있음 (bestCards를 만드는 과정에서 정렬 함수가 쓰임)
func compare(player1 Player, player2 Player) CardCompareResult {
	if player1.handsRank > player2.handsRank {
		return Player1Win
	} else if player1.handsRank < player2.handsRank {
		return Player2Win
	} else {
		if player1.highCard > player2.highCard {
			return Player1Win
		} else if player1.highCard < player2.highCard {
			return Player2Win
		} else {
			for i := len(player1.bestCards) - 1; i >= 0; i-- {
				if player1.bestCards[i].Rank > player2.bestCards[i].Rank {
					return Player1Win
				} else if player1.bestCards[i].Rank < player2.bestCards[i].Rank {
					return Player2Win
				}
			}
			return Draw
		}
	}
}

//type Players []Player
//
//func (p Players) Len() int { return len(p) }
//func (p Players) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
//func (p Players) Less(i, j int) bool {
//	if p[i].handsRank < p[j].handsRank {
//		return p[i].handsRank < p[j].handsRank
//	} else if p[i].handsRank > p[j].handsRank {
//		return p[i].handsRank > p[j].handsRank
//	} else {
//		if p[i].highCard < p[j].highCard {
//			return p[i].highCard < p[j].highCard
//		} else if p[i].highCard > p[j].highCard {
//			return p[i].highCard > p[j].highCard
//		} else {
//			for k:=4; k>=0; k-- {
//				if p[i].bestCards[k].Rank < p[j].bestCards[k].Rank {
//					return p[i].bestCards[k].Rank < p[j].bestCards[k].Rank
//				}
//				if p[i].bestCards[k].Rank > p[j].bestCards[k].Rank {
//					return p[i].bestCards[k].Rank > p[j].bestCards[k].Rank
//				}
//			}
//			return false
//		}
//	}
//}
