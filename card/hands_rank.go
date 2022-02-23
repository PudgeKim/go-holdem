package card

import (
	"sort"
)

type HandsRank int

const (
	HighCard HandsRank = iota
	OnePair
	TwoPair
	Triple
	Straight
	Flush
	FullHouse
	FourCard
	StraightFlush
	RoyalStraightFlush
)

// 7장의 카드로부터 만들 수 있는 조합 중
// 가장 좋은 5장의 조합을 리턴함
func getBestHandsRank(cards []Card) ([]Card, HandsRank, Rank) {
	if len(cards) != 7 {
		panic("cards' length must be 7")
	}

	var allCombs [][]Card
	makeAllCombinations(cards, []Card{}, &allCombs, 0, 0)

	// 첫번째 족보를 얻어낸 후에 계속 비교해나가서 가장 높은 핸드를 얻어냄
	bestCards := allCombs[0]
	sortCards(bestCards)
	bestHandsRank, bestHighCard := checkHandsRank(bestCards)

	for i := 1; i < len(allCombs); i++ {
		curCards := allCombs[i]
		sortCards(curCards)
		curHandsRank, curHighCard := checkHandsRank(curCards)

		if bestHandsRank < curHandsRank {
			bestCards = curCards
			bestHandsRank = curHandsRank
			bestHighCard = curHighCard
			continue
		}

		if bestHandsRank == curHandsRank {
			if bestHighCard < curHighCard {
				bestCards = curCards
				bestHighCard = curHighCard
				continue
			}

			// 족보도 같고 highCard도 같으면 더 세부적으로 비교를 해봐야함
			// 예를 들어 둘 다 3풀하우스라면 33322보다 333QQ가 더 높음
			// 원페어라면 원페어를 제외한 나머지 숫자들의 크기를 따져봐야함
			// 오름차순 정렬되어있으므로 위에서부터 비교해보면됨
			if bestHighCard == curHighCard {
				for j := 4; j >= 0; j-- {
					if bestCards[j].Rank < curCards[j].Rank {
						bestCards = curCards
						bestHandsRank = curHandsRank
						bestHighCard = curHighCard
						break
					}
				}
			}
		}
	}

	return bestCards, bestHandsRank, bestHighCard
}

func sortCards(cards []Card) {
	sort.Slice(cards, func(i, j int) bool { return cards[i].Rank < cards[j].Rank })
}

// 첫번째 리턴 값은 족보를, 두번째 리턴 값은 하이카드를 나타냄
// 예를 들어 5스트레이트라면 첫번째 값은 스트레이트를, 두번째 값은 5
func checkHandsRank(cards []Card) (HandsRank, Rank) {
	sortCards(cards)

	var res bool
	var highCard Rank

	res, highCard = isRoyalStraightFlush(cards)
	if res == true {
		return RoyalStraightFlush, highCard
	}

	res, highCard = isStraightFlush(cards)
	if res == true {
		return StraightFlush, highCard
	}

	res, highCard = isFourCard(cards)
	if res == true {
		return FourCard, highCard
	}

	res, highCard = isFullHouse(cards)
	if res == true {
		return FullHouse, highCard
	}

	res, highCard = isFlush(cards)
	if res == true {
		return Flush, highCard
	}

	res, highCard = isStraight(cards)
	if res == true {
		return Straight, highCard
	}

	res, highCard = isTriple(cards)
	if res == true {
		return Triple, highCard
	}

	res, highCard = isTwoPair(cards)
	if res == true {
		return TwoPair, highCard
	}

	res, highCard = isOnePair(cards)
	if res == true {
		return OnePair, highCard
	}

	// 정렬되어있으므로 마지막 카드가 가장 높은 카드
	return HighCard, cards[4].Rank
}

// 7장의 카드로 만들 수 있는 5장의 모든 조합
// cards는 7장의 카드
// tmpCards는 allComb에 저장하기 전 임시 배열
// allComb는 처음에는 빈 배열로 함수가 끝나면 모든 조합이 들어가게됨
func makeAllCombinations(cards []Card, tmpCards []Card, allComb *[][]Card, lv int, startIdx int) {
	if lv == 5 {
		copied := make([]Card, 5)
		copy(copied, tmpCards)
		*allComb = append(*allComb, copied)
		return
	}

	for i := startIdx; i < len(cards); i++ {
		tmpCards = append(tmpCards, cards[i])
		makeAllCombinations(cards, tmpCards, allComb, lv+1, i+1)
		tmpCards = tmpCards[:len(tmpCards)-1]
	}
}

// ******
// 아래 함수들은 인자로 들어오는 cards가 정렬이 되어있다고 가정함
// 리턴 값에서 Rank는 예를 들어 4, 2 투페어인 경우 리턴되는 Rank 값은 4임
// ******

func isRoyalStraightFlush(cards []Card) (bool, Rank) {
	if cards[0].Rank != Ten {
		return false, None
	}

	for i := 1; i < 5; i++ {
		if cards[i].Symbol != cards[i-1].Symbol {
			return false, None
		}
		if cards[i].Rank != cards[i-1].Rank+1 {
			return false, None
		}
	}

	return true, Ace
}

func isStraightFlush(cards []Card) (bool, Rank) {
	// A가 14로 취급되기 때문에 A로 시작하는 경우와 아닌 경우로 나누어야함
	// A는 14이고 카드는 정렬되어있으므로 마지막 카드가 A인지 검사해야함(A로 시작한다면)
	if cards[4].Rank == Ace {
		// 4번째 카드는 A와 문양이 같고 5여야함 (StraightFlush를 만족하려면)
		if cards[4].Symbol != cards[3].Symbol {
			return false, None
		}
		if cards[3].Rank != Five {
			return false, None
		}

		// 마지막 카드는 A인걸 확인했으니 첫번째부터 네번째까지 검사해봄
		for i := 1; i < 4; i++ {
			if cards[i].Symbol != cards[i-1].Symbol {
				return false, None
			}
			if cards[i].Rank != cards[i-1].Rank+1 {
				return false, None
			}
		}

		// 5 스트레이트이므로 마지막 카드는 A이므로 4번째 카드인 5를 리턴
		return true, Five
	}
	// A로 시작하지 않는 경우
	for i := 1; i < 5; i++ {
		if cards[i].Symbol != cards[i-1].Symbol {
			return false, None
		}
		if cards[i].Rank != cards[i-1].Rank+1 {
			return false, None
		}
	}

	return true, cards[4].Rank
}

func isFourCard(cards []Card) (bool, Rank) {
	// 1번째부터 4번째까지 포카드인 경우
	if (cards[0].Rank == cards[1].Rank) &&
		(cards[1].Rank == cards[2].Rank) &&
		(cards[2].Rank == cards[3].Rank) {
		return true, cards[3].Rank
	}

	// 2번째부터 5번째까지 포카드인 경우
	if (cards[1].Rank == cards[2].Rank) &&
		(cards[2].Rank == cards[3].Rank) &&
		(cards[3].Rank == cards[4].Rank) {
		return true, cards[4].Rank
	}

	return false, None
}

func isFullHouse(cards []Card) (bool, Rank) {
	// 1번째부터 3번째가 트리플인 경우
	if (cards[0].Rank == cards[1].Rank) &&
		(cards[1].Rank == cards[2].Rank) &&
		(cards[3].Rank == cards[4].Rank) { // 이건 원페어 체크
		return true, cards[2].Rank
	}

	// 3번째부터 5번째까지 트리플인 경우
	if (cards[0].Rank == cards[1].Rank) && // 원페어 체크
		(cards[2].Rank == cards[3].Rank) &&
		(cards[3].Rank == cards[4].Rank) {
		return true, cards[4].Rank
	}

	return false, None
}

func isFlush(cards []Card) (bool, Rank) {
	for i := 1; i < 5; i++ {
		if cards[i].Symbol != cards[i-1].Symbol {
			return false, None
		}
	}
	return true, cards[4].Rank
}

func isStraight(cards []Card) (bool, Rank) {
	// A가 첫번째인 스트레이트인 경우 (A, 2, 3, 4, 5)
	// (정렬되어있고 A는 14므로 마지막 카드가 A)
	if (cards[4].Rank == Ace) &&
		(cards[0].Rank == Two) &&
		(cards[1].Rank == Three) &&
		(cards[2].Rank == Four) &&
		(cards[3].Rank == Five) {
		return true, Five
	}

	// A가 포함되지 않은 경우
	for i := 1; i < 5; i++ {
		if cards[i].Rank != cards[i-1].Rank+1 {
			return false, None
		}
	}
	return true, cards[4].Rank
}

func isTriple(cards []Card) (bool, Rank) {
	// 1번째부터 3번째까지 트리플인 경우
	if (cards[0].Rank == cards[1].Rank) &&
		(cards[1].Rank == cards[2].Rank) {
		return true, cards[2].Rank
	}

	// 2번째부터 4번째까지 트리플인 경우
	if (cards[1].Rank == cards[2].Rank) &&
		(cards[2].Rank == cards[3].Rank) {
		return true, cards[3].Rank
	}

	// 3번째부터 5번째까지 트리플인 경우
	if (cards[2].Rank == cards[3].Rank) &&
		(cards[3].Rank == cards[4].Rank) {
		return true, cards[4].Rank
	}

	return false, None
}

func isTwoPair(cards []Card) (bool, Rank) {
	// 1번째부터 4번째까지 투페어인 경우 (1번이랑 2번이 같고 3번이랑 4번이 같은 경우)
	if (cards[0].Rank == cards[1].Rank) &&
		(cards[2].Rank == cards[3].Rank) {
		return true, cards[3].Rank
	}

	// 2번째부터 5번째까지 투페어인 경우
	if (cards[1].Rank == cards[2].Rank) &&
		(cards[3].Rank == cards[4].Rank) {
		return true, cards[4].Rank
	}

	return false, None
}

// 이 함수의 경우 twoPair나 triple도 onePair로 취급되지만
// 나중에 검사를 할 때 royalStraightFlush부터
// 쭉 아래로 검사를 하기때문에 상관없음
func isOnePair(cards []Card) (bool, Rank) {
	for i := 1; i < 5; i++ {
		if cards[i].Rank == cards[i-1].Rank {
			return true, cards[i].Rank
		}
	}
	return false, None
}
