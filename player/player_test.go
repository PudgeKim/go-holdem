package player

import (
	"github.com/PudgeKim/go-holdem/card"
	"reflect"
	"testing"
)

func TestCompare(t *testing.T) {
	var p1 Player
	var p2 Player
	var res CardCompareResult

	p1 = Player{
		bestCards: []card.Card{
			{Symbol: card.Heart, Rank: card.Five},
			{Symbol: card.Diamond, Rank: card.Five},
			{Symbol: card.Diamond, Rank: card.Ten},
			{Symbol: card.Clover, Rank: card.Two},
			{Symbol: card.Clover, Rank: card.Ace},
		},
		handsRank: card.OnePair,
		highCard:  card.Five,
	}
	card.SortCards(p1.bestCards)

	p2 = Player{
		bestCards: []card.Card{
			{Symbol: card.Heart, Rank: card.Two},
			{Symbol: card.Diamond, Rank: card.Three},
			{Symbol: card.Heart, Rank: card.Ten},
			{Symbol: card.Clover, Rank: card.Six},
			{Symbol: card.Clover, Rank: card.Jack},
		},
		handsRank: card.HighCard,
		highCard:  card.Jack,
	}
	card.SortCards(p2.bestCards)

	res = compare(p1, p2)
	if res != Player1Win {
		t.Error("Player1 should win")
	}

	p1 = Player{
		bestCards: []card.Card{
			{Symbol: card.Heart, Rank: card.Five},
			{Symbol: card.Diamond, Rank: card.Five},
			{Symbol: card.Diamond, Rank: card.Ten},
			{Symbol: card.Clover, Rank: card.Two},
			{Symbol: card.Clover, Rank: card.Jack},
		},
		handsRank: card.OnePair,
		highCard:  card.Five,
	}
	card.SortCards(p1.bestCards)

	p2 = Player{
		bestCards: []card.Card{
			{Symbol: card.Heart, Rank: card.Two},
			{Symbol: card.Diamond, Rank: card.Ace},
			{Symbol: card.Spade, Rank: card.Five},
			{Symbol: card.Spade, Rank: card.Six},
			{Symbol: card.Spade, Rank: card.Five},
		},
		handsRank: card.OnePair,
		highCard:  card.Five,
	}
	card.SortCards(p2.bestCards)

	res = compare(p1, p2)
	if res != Player2Win {
		t.Error("Player2 should win")
	}

}

func TestGetWinners(t *testing.T) {
	var p1 Player
	var p2 Player
	var p3 Player
	var winners []Player

	p1 = Player{
		bestCards: []card.Card{
			{Symbol: card.Heart, Rank: card.Two},
			{Symbol: card.Diamond, Rank: card.Ace},
			{Symbol: card.Spade, Rank: card.Five},
			{Symbol: card.Spade, Rank: card.Six},
			{Symbol: card.Spade, Rank: card.Five},
		},
		handsRank: card.OnePair,
		highCard:  card.Five,
	}
	card.SortCards(p1.bestCards)

	p2 = Player{
		bestCards: []card.Card{
			{Symbol: card.Clover, Rank: card.Two},
			{Symbol: card.Diamond, Rank: card.Two},
			{Symbol: card.Diamond, Rank: card.Five},
			{Symbol: card.Clover, Rank: card.Six},
			{Symbol: card.Heart, Rank: card.Six},
		},
		handsRank: card.TwoPair,
		highCard:  card.Six,
	}
	card.SortCards(p2.bestCards)

	p3 = Player{
		bestCards: []card.Card{
			{Symbol: card.Heart, Rank: card.King},
			{Symbol: card.Diamond, Rank: card.King},
			{Symbol: card.Spade, Rank: card.Ace},
			{Symbol: card.Clover, Rank: card.Ace},
			{Symbol: card.Heart, Rank: card.Four},
		},
		handsRank: card.TwoPair,
		highCard:  card.Ace,
	}
	card.SortCards(p3.bestCards)

	winners = GetWinners([]Player{p1, p2, p3})
	if !reflect.DeepEqual(winners, []Player{p3}) {
		t.Error("winner should be p3")
	}

	p1 = Player{
		bestCards: []card.Card{
			{Symbol: card.Heart, Rank: card.Two},
			{Symbol: card.Diamond, Rank: card.Ace},
			{Symbol: card.Spade, Rank: card.Five},
			{Symbol: card.Spade, Rank: card.Three},
			{Symbol: card.Spade, Rank: card.Four},
		},
		handsRank: card.Straight,
		highCard:  card.Five,
	}
	card.SortCards(p1.bestCards)

	p2 = Player{
		bestCards: []card.Card{
			{Symbol: card.Clover, Rank: card.Ace},
			{Symbol: card.Diamond, Rank: card.Two},
			{Symbol: card.Diamond, Rank: card.Three},
			{Symbol: card.Clover, Rank: card.Four},
			{Symbol: card.Heart, Rank: card.Five},
		},
		handsRank: card.Straight,
		highCard:  card.Five,
	}
	card.SortCards(p2.bestCards)

	p3 = Player{
		bestCards: []card.Card{
			{Symbol: card.Heart, Rank: card.Five},
			{Symbol: card.Diamond, Rank: card.Four},
			{Symbol: card.Spade, Rank: card.Ace},
			{Symbol: card.Clover, Rank: card.Two},
			{Symbol: card.Heart, Rank: card.Three},
		},
		handsRank: card.Straight,
		highCard:  card.Five,
	}
	card.SortCards(p3.bestCards)

	winners = GetWinners([]Player{p1, p2, p3})
	if !reflect.DeepEqual(winners, []Player{p1, p2, p3}) {
		t.Error("winner should be p1, p2, p3")
	}
}
