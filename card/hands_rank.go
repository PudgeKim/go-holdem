package card

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

// 7장의 카드로 만들 수 있는 5장의 모든 조합
// cards는 7장의 카드
// tmpCards는 allComb에 저장하기 전 임시 배열
// allComb는 처음에는 빈 배열로 함수가 끝나면 모든 조합이 들어가게됨
func makeAllCombinations(cards []Card, tmpCards []Card, allComb *[][]Card, lv int, startIdx int) {
	if lv == 5 {
		*allComb = append(*allComb, tmpCards)
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
	if cards[0].rank != Ten {
		return false, None
	}

	for i := 1; i < 5; i++ {
		if cards[i].symbol != cards[i-1].symbol {
			return false, None
		}
		if cards[i].rank != cards[i-1].rank+1 {
			return false, None
		}
	}

	return true, Ace
}

func isStraightFlush(cards []Card) (bool, Rank) {
	// A가 14로 취급되기 때문에 A로 시작하는 경우와 아닌 경우로 나누어야함
	// A는 14이고 카드는 정렬되어있으므로 마지막 카드가 A인지 검사해야함(A로 시작한다면)
	if cards[4].rank == Ace {
		// 4번째 카드는 A와 문양이 같고 5여야함 (StraightFlush를 만족하려면)
		if cards[4].symbol != cards[3].symbol {
			return false, None
		}
		if cards[3].rank != Five {
			return false, None
		}

		// 마지막 카드는 A인걸 확인했으니 첫번째부터 네번째까지 검사해봄
		for i := 1; i < 4; i++ {
			if cards[i].symbol != cards[i-1].symbol {
				return false, None
			}
			if cards[i].rank != cards[i-1].rank+1 {
				return false, None
			}
		}

		// 5 스트레이트이므로 마지막 카드는 A이므로 4번째 카드인 5를 리턴
		return true, Five
	}
	// A로 시작하지 않는 경우
	for i := 1; i < 5; i++ {
		if cards[i].symbol != cards[i-1].symbol {
			return false, None
		}
		if cards[i].rank != cards[i-1].rank+1 {
			return false, None
		}
	}

	return true, cards[4].rank
}

func isFourCard(cards []Card) (bool, Rank) {
	// 1번째부터 4번째까지 포카드인 경우
	if (cards[0].rank == cards[1].rank) &&
		(cards[1].rank == cards[2].rank) &&
		(cards[2].rank == cards[3].rank) {
		return true, cards[3].rank
	}

	// 2번째부터 5번째까지 포카드인 경우
	if (cards[1].rank == cards[2].rank) &&
		(cards[2].rank == cards[3].rank) &&
		(cards[3].rank == cards[4].rank) {
		return true, cards[4].rank
	}

	return false, None
}

func isFullHouse(cards []Card) (bool, Rank) {
	// 1번째부터 3번째가 트리플인 경우
	if (cards[0].rank == cards[1].rank) &&
		(cards[1].rank == cards[2].rank) &&
		(cards[3].rank == cards[4].rank) { // 이건 원페어 체크
		return true, cards[2].rank
	}

	// 3번째부터 5번째까지 트리플인 경우
	if (cards[0].rank == cards[1].rank) && // 원페어 체크
		(cards[2].rank == cards[3].rank) &&
		(cards[3].rank == cards[4].rank) {
		return true, cards[4].rank
	}

	return false, None
}

func isFlush(cards []Card) (bool, Rank) {
	for i := 1; i < 5; i++ {
		if cards[i].symbol != cards[i-1].symbol {
			return false, None
		}
	}
	return true, cards[4].rank
}

func isStraight(cards []Card) (bool, Rank) {
	// A가 첫번째인 스트레이트인 경우 (A, 2, 3, 4, 5)
	// (정렬되어있고 A는 14므로 마지막 카드가 A)
	if (cards[4].rank == Ace) &&
		(cards[0].rank == Two) &&
		(cards[1].rank == Three) &&
		(cards[2].rank == Four) &&
		(cards[3].rank == Five) {
		return true, Five
	}

	// A가 포함되지 않은 경우
	for i := 1; i < 5; i++ {
		if cards[i].rank != cards[i-1].rank+1 {
			return false, None
		}
	}
	return true, cards[4].rank
}

func isTriple(cards []Card) (bool, Rank) {
	// 1번째부터 3번째까지 트리플인 경우
	if (cards[0].rank == cards[1].rank) &&
		(cards[1].rank == cards[2].rank) {
		return true, cards[2].rank
	}

	// 2번째부터 4번째까지 트리플인 경우
	if (cards[1].rank == cards[2].rank) &&
		(cards[2].rank == cards[3].rank) {
		return true, cards[3].rank
	}

	// 3번째부터 5번째까지 트리플인 경우
	if (cards[2].rank == cards[3].rank) &&
		(cards[3].rank == cards[4].rank) {
		return true, cards[4].rank
	}

	return false, None
}

func isTwoPair(cards []Card) (bool, Rank) {
	// 1번째부터 4번째까지 투페어인 경우 (1번이랑 2번이 같고 3번이랑 4번이 같은 경우)
	if (cards[0].rank == cards[1].rank) &&
		(cards[2].rank == cards[3].rank) {
		return true, cards[3].rank
	}

	// 2번째부터 5번째까지 투페어인 경우
	if (cards[1].rank == cards[2].rank) &&
		(cards[3].rank == cards[4].rank) {
		return true, cards[4].rank
	}

	return false, None
}

// 이 함수의 경우 twoPair나 triple도 onePair로 취급되지만
// 나중에 검사를 할 때 royalStraightFlush부터
// 쭉 아래로 검사를 하기때문에 상관없음
func isOnePair(cards []Card) (bool, Rank) {
	for i := 1; i < 5; i++ {
		if cards[i].rank == cards[i-1].rank {
			return true, cards[i].rank
		}
	}
	return false, None
}
