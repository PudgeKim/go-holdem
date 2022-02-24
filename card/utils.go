package card

import "sort"

func SortCards(cards []Card) {
	sort.Slice(cards, func(i, j int) bool { return cards[i].Rank < cards[j].Rank })
}
