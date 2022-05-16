package entity

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/PudgeKim/go-holdem/card"
// 	"github.com/PudgeKim/go-holdem/game"
// )

// func TestCompare(t *testing.T) {
// 	var p1 Player
// 	var p2 Player
// 	var res game.CardCompareResult

// 	p1 = Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Heart, Rank: card.Five},
// 			{Symbol: card.Diamond, Rank: card.Five},
// 			{Symbol: card.Diamond, Rank: card.Ten},
// 			{Symbol: card.Clover, Rank: card.Two},
// 			{Symbol: card.Clover, Rank: card.Ace},
// 		},
// 		HandsRank: card.OnePair,
// 		HighCard:  card.Five,
// 	}
// 	card.SortCards(p1.BestCards)

// 	p2 = Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Heart, Rank: card.Two},
// 			{Symbol: card.Diamond, Rank: card.Three},
// 			{Symbol: card.Heart, Rank: card.Ten},
// 			{Symbol: card.Clover, Rank: card.Six},
// 			{Symbol: card.Clover, Rank: card.Jack},
// 		},
// 		HandsRank: card.HighCard,
// 		HighCard:  card.Jack,
// 	}
// 	card.SortCards(p2.BestCards)

// 	res = compare(&p1, &p2)
// 	if res != Player1Win {
// 		t.Error("Player1 should win")
// 	}

// 	p1 = Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Heart, Rank: card.Five},
// 			{Symbol: card.Diamond, Rank: card.Five},
// 			{Symbol: card.Diamond, Rank: card.Ten},
// 			{Symbol: card.Clover, Rank: card.Two},
// 			{Symbol: card.Clover, Rank: card.Jack},
// 		},
// 		HandsRank: card.OnePair,
// 		HighCard:  card.Five,
// 	}
// 	card.SortCards(p1.BestCards)

// 	p2 = Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Heart, Rank: card.Two},
// 			{Symbol: card.Diamond, Rank: card.Ace},
// 			{Symbol: card.Spade, Rank: card.Five},
// 			{Symbol: card.Spade, Rank: card.Six},
// 			{Symbol: card.Spade, Rank: card.Five},
// 		},
// 		HandsRank: card.OnePair,
// 		HighCard:  card.Five,
// 	}
// 	card.SortCards(p2.BestCards)

// 	res = compare(&p1, &p2)
// 	if res != Player2Win {
// 		t.Error("Player2 should win")
// 	}

// }

// func TestGetWinners(t *testing.T) {
// 	var p1 Player
// 	var p2 Player
// 	var p3 Player
// 	var winners []*Player

// 	p1 = Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Heart, Rank: card.Two},
// 			{Symbol: card.Diamond, Rank: card.Ace},
// 			{Symbol: card.Spade, Rank: card.Five},
// 			{Symbol: card.Spade, Rank: card.Six},
// 			{Symbol: card.Spade, Rank: card.Five},
// 		},
// 		HandsRank: card.OnePair,
// 		HighCard:  card.Five,
// 	}
// 	card.SortCards(p1.BestCards)

// 	p2 = Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Clover, Rank: card.Two},
// 			{Symbol: card.Diamond, Rank: card.Two},
// 			{Symbol: card.Diamond, Rank: card.Five},
// 			{Symbol: card.Clover, Rank: card.Six},
// 			{Symbol: card.Heart, Rank: card.Six},
// 		},
// 		HandsRank: card.TwoPair,
// 		HighCard:  card.Six,
// 	}
// 	card.SortCards(p2.BestCards)

// 	p3 = Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Heart, Rank: card.King},
// 			{Symbol: card.Diamond, Rank: card.King},
// 			{Symbol: card.Spade, Rank: card.Ace},
// 			{Symbol: card.Clover, Rank: card.Ace},
// 			{Symbol: card.Heart, Rank: card.Four},
// 		},
// 		HandsRank: card.TwoPair,
// 		HighCard:  card.Ace,
// 	}
// 	card.SortCards(p3.BestCards)

// 	winners, err := GetWinners([]*Player{&p1, &p2, &p3})
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	if !reflect.DeepEqual(winners, []*Player{&p3}) {
// 		t.Error("winner should be p3")
// 	}

// 	p1 = Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Heart, Rank: card.Two},
// 			{Symbol: card.Diamond, Rank: card.Ace},
// 			{Symbol: card.Spade, Rank: card.Five},
// 			{Symbol: card.Spade, Rank: card.Three},
// 			{Symbol: card.Spade, Rank: card.Four},
// 		},
// 		HandsRank: card.Straight,
// 		HighCard:  card.Five,
// 	}
// 	card.SortCards(p1.BestCards)

// 	p2 = Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Clover, Rank: card.Ace},
// 			{Symbol: card.Diamond, Rank: card.Two},
// 			{Symbol: card.Diamond, Rank: card.Three},
// 			{Symbol: card.Clover, Rank: card.Four},
// 			{Symbol: card.Heart, Rank: card.Five},
// 		},
// 		HandsRank: card.Straight,
// 		HighCard:  card.Five,
// 	}
// 	card.SortCards(p2.BestCards)

// 	p3 = Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Heart, Rank: card.Five},
// 			{Symbol: card.Diamond, Rank: card.Four},
// 			{Symbol: card.Spade, Rank: card.Ace},
// 			{Symbol: card.Clover, Rank: card.Two},
// 			{Symbol: card.Heart, Rank: card.Three},
// 		},
// 		HandsRank: card.Straight,
// 		HighCard:  card.Five,
// 	}
// 	card.SortCards(p3.BestCards)

// 	winners, err = GetWinners([]*Player{&p1, &p2, &p3})
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	if !reflect.DeepEqual(winners, []*Player{&p1, &p2, &p3}) {
// 		t.Error("winner should be p1, p2, p3")
// 	}

// 	p1 = Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Heart, Rank: card.Two},
// 			{Symbol: card.Diamond, Rank: card.Ace},
// 			{Symbol: card.Spade, Rank: card.Five},
// 			{Symbol: card.Spade, Rank: card.Three},
// 			{Symbol: card.Spade, Rank: card.Four},
// 		},
// 		HandsRank: card.Straight,
// 		HighCard:  card.Five,
// 	}
// 	card.SortCards(p1.BestCards)

// 	p2 = Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Clover, Rank: card.Ace},
// 			{Symbol: card.Diamond, Rank: card.Ace},
// 			{Symbol: card.Diamond, Rank: card.Three},
// 			{Symbol: card.Heart, Rank: card.Ace},
// 			{Symbol: card.Heart, Rank: card.Five},
// 		},
// 		HandsRank: card.Triple,
// 		HighCard:  card.Ace,
// 	}
// 	card.SortCards(p2.BestCards)

// 	p3 = Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Heart, Rank: card.Five},
// 			{Symbol: card.Diamond, Rank: card.Five},
// 			{Symbol: card.Spade, Rank: card.Ten},
// 			{Symbol: card.Clover, Rank: card.Nine},
// 			{Symbol: card.Heart, Rank: card.Eight},
// 		},
// 		HandsRank: card.OnePair,
// 		HighCard:  card.Five,
// 	}
// 	card.SortCards(p3.BestCards)

// 	p4 := Player{
// 		BestCards: []card.Card{
// 			{Symbol: card.Heart, Rank: card.Five},
// 			{Symbol: card.Diamond, Rank: card.Five},
// 			{Symbol: card.Spade, Rank: card.Ten},
// 			{Symbol: card.Clover, Rank: card.Nine},
// 			{Symbol: card.Heart, Rank: card.Eight},
// 		},
// 		HandsRank: card.OnePair,
// 		HighCard:  card.Five,
// 	}
// 	card.SortCards(p4.BestCards)

// 	p1.TotalBet = 100
// 	p2.TotalBet = 150
// 	p3.TotalBet = 200
// 	p4.TotalBet = 200

// 	winners, err = GetWinners([]*Player{&p1, &p2, &p3, &p4})
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	if !reflect.DeepEqual(winners, []*Player{&p1, &p2, &p3, &p4}) {
// 		t.Error("winner should be p1, p2, p3, p4")
// 	}

// }
