package card

import (
	"testing"
)

var cards []Card
var res bool
var highCard Rank

func TestMakeAllCombinations(t *testing.T) {
	cards = []Card{
		{symbol: Diamond, rank: Three},
		{symbol: Diamond, rank: Four},
		{symbol: Spade, rank: Five},
		{symbol: Heart, rank: Five},
		{symbol: Diamond, rank: Nine},
		{symbol: Spade, rank: Queen},
		{symbol: Spade, rank: Ace},
	}

	var tmpCards []Card
	var allCombs [][]Card
	makeAllCombinations(cards, tmpCards, &allCombs, 0, 0)
	if len(allCombs) != 21 {
		t.Error("All combination's length should be 21 (7C5)")
	}

}

func TestRoyalStraightFlush(t *testing.T) {
	cards = []Card{
		{symbol: Spade, rank: Ten},
		{symbol: Spade, rank: Jack},
		{symbol: Spade, rank: Queen},
		{symbol: Spade, rank: King},
		{symbol: Spade, rank: Ace},
	}
	res, highCard = isRoyalStraightFlush(cards)
	if !res {
		t.Error("It should be RoyalStraightFlush but wrong result")
	}
	if highCard != Ace {
		t.Error("RoyalStraightFlush must be end with Ace")
	}

	cards2 := []Card{
		{symbol: Spade, rank: Nine},
		{symbol: Spade, rank: Ten},
		{symbol: Spade, rank: Jack},
		{symbol: Spade, rank: Queen},
		{symbol: Spade, rank: King},
	}
	res, highCard = isRoyalStraightFlush(cards2)
	if res {
		t.Error("RoyalStraightFlush's highCard should be Ace")
	}
}

func TestStraightFlush(t *testing.T) {
	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Spade, rank: Three},
		{symbol: Spade, rank: Four},
		{symbol: Spade, rank: Five},
		{symbol: Spade, rank: Ace},
	}
	res, highCard = isStraightFlush(cards)
	if !res {
		t.Error("It should be StraightFlush")
	}
	if highCard != Five {
		t.Error("highCard should be Five")
	}

	cards = []Card{
		{symbol: Spade, rank: Four},
		{symbol: Spade, rank: Five},
		{symbol: Spade, rank: Six},
		{symbol: Spade, rank: Seven},
		{symbol: Spade, rank: Eight},
	}
	res, highCard = isStraightFlush(cards)
	if !res {
		t.Error("It should be StraightFlush")
	}
	if highCard != Eight {
		t.Error("highCard should be Eight")
	}

	cards = []Card{
		{symbol: Spade, rank: Four},
		{symbol: Spade, rank: Five},
		{symbol: Spade, rank: Seven},
		{symbol: Spade, rank: Eight},
		{symbol: Spade, rank: Ten},
	}
	res, _ = isStraightFlush(cards)
	if res {
		t.Error("It should not be StraightFlush")
	}

	cards = []Card{
		{symbol: Spade, rank: Four},
		{symbol: Heart, rank: Five},
		{symbol: Spade, rank: Seven},
		{symbol: Spade, rank: Eight},
		{symbol: Spade, rank: Ten},
	}
	res, _ = isStraightFlush(cards)
	if res {
		t.Error("It should not be StraightFlush")
	}
}

func TestFourCard(t *testing.T) {
	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Two},
		{symbol: Spade, rank: Two},
		{symbol: Spade, rank: Two},
		{symbol: Spade, rank: Ten},
	}
	res, highCard = isFourCard(cards)
	if !res {
		t.Error("It should be FourCard")
	}
	if highCard != Two {
		t.Error("highCard should be Two")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Three},
		{symbol: Spade, rank: Three},
		{symbol: Spade, rank: Three},
	}
	res, highCard = isFourCard(cards)
	if !res {
		t.Error("It should be FourCard")
	}
	if highCard != Three {
		t.Error("highCard should be Three")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Four},
		{symbol: Spade, rank: Six},
		{symbol: Spade, rank: Ten},
	}
	res, _ = isFourCard(cards)
	if res {
		t.Error("It should not be FourCard")
	}
}

func TestFullHouse(t *testing.T) {
	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Two},
		{symbol: Spade, rank: Two},
		{symbol: Spade, rank: Three},
		{symbol: Spade, rank: Three},
	}
	res, highCard = isFullHouse(cards)
	if !res {
		t.Error("It should be FullHouse")
	}
	if highCard != Two {
		t.Error("highChard should be Two")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Two},
		{symbol: Spade, rank: Three},
		{symbol: Spade, rank: Three},
		{symbol: Spade, rank: Three},
	}
	res, highCard = isFullHouse(cards)
	if !res {
		t.Error("It should be FullHouse")
	}
	if highCard != Three {
		t.Error("highCard should be Three")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Three},
		{symbol: Spade, rank: Four},
		{symbol: Spade, rank: Five},
	}
	res, _ = isFullHouse(cards)
	if res {
		t.Error("It should not be FullHouse")
	}
}

func TestFlush(t *testing.T) {
	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Spade, rank: Two},
		{symbol: Spade, rank: Two},
		{symbol: Spade, rank: Three},
		{symbol: Spade, rank: Three},
	}
	res, highCard = isFlush(cards)
	if !res {
		t.Error("It should be Flush")
	}
	if highCard != Three {
		t.Error("highCard should be Three")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Two},
		{symbol: Spade, rank: Two},
		{symbol: Spade, rank: Three},
		{symbol: Spade, rank: Three},
	}
	res, _ = isFlush(cards)
	if res {
		t.Error("It should not be Flush")
	}
}

func TestStraight(t *testing.T) {
	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Four},
		{symbol: Heart, rank: Five},
		{symbol: Spade, rank: Six},
	}
	res, highCard = isStraight(cards)
	if !res {
		t.Error("It should be Straight")
	}
	if highCard != Six {
		t.Error("highCard should be Six")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Four},
		{symbol: Spade, rank: Five},
		{symbol: Heart, rank: Ace},
	}
	res, highCard = isStraight(cards)
	if !res {
		t.Error("It should be Straight")
	}
	if highCard != Five {
		t.Error("highCard should be Five")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Four},
		{symbol: Spade, rank: King},
		{symbol: Heart, rank: Ace},
	}
	res, _ = isStraight(cards)
	if res {
		t.Error("It should not be Straight")
	}
}

func TestTriple(t *testing.T) {
	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Two},
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Five},
		{symbol: Spade, rank: Six},
	}
	res, highCard = isTriple(cards)
	if !res {
		t.Error("It should be Triple")
	}
	if highCard != Two {
		t.Error("highCard should be Two")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Five},
		{symbol: Spade, rank: Five},
		{symbol: Heart, rank: Five},
		{symbol: Spade, rank: Six},
	}
	res, highCard = isTriple(cards)
	if !res {
		t.Error("It should be Triple")
	}
	if highCard != Five {
		t.Error("highCard should be Five")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Six},
		{symbol: Heart, rank: Six},
		{symbol: Spade, rank: Six},
	}
	res, highCard = isTriple(cards)
	if !res {
		t.Error("It should be Triple")
	}
	if highCard != Six {
		t.Error("highCard should be Six")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Six},
		{symbol: Heart, rank: Ten},
		{symbol: Spade, rank: Jack},
	}
	res, _ = isTriple(cards)
	if res {
		t.Error("It should not be Triple")
	}
}

func TestTwoPair(t *testing.T) {
	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Two},
		{symbol: Spade, rank: Three},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Six},
	}
	res, highCard = isTwoPair(cards)
	if !res {
		t.Error("It should be TwoPair")
	}
	if highCard != Three {
		t.Error("highCard should be Three")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Three},
		{symbol: Heart, rank: Six},
		{symbol: Spade, rank: Six},
	}
	res, highCard = isTwoPair(cards)
	if !res {
		t.Error("It should be TwoPair")
	}
	if highCard != Six {
		t.Error("highCard should be Six")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Three},
		{symbol: Heart, rank: Five},
		{symbol: Spade, rank: Six},
	}
	res, _ = isTwoPair(cards)
	if res {
		t.Error("It should not be TwoPair")
	}
}

func TestOnePair(t *testing.T) {
	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Two},
		{symbol: Spade, rank: Three},
		{symbol: Heart, rank: Five},
		{symbol: Spade, rank: Six},
	}
	res, highCard = isOnePair(cards)
	if !res {
		t.Error("It should be OnePair")
	}
	if highCard != Two {
		t.Error("highCard should be Two")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Three},
		{symbol: Heart, rank: Five},
		{symbol: Spade, rank: Six},
	}
	res, highCard = isOnePair(cards)
	if !res {
		t.Error("It should be OnePair")
	}
	if highCard != Three {
		t.Error("highCard should be Three")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Five},
		{symbol: Heart, rank: Five},
		{symbol: Spade, rank: Six},
	}
	res, highCard = isOnePair(cards)
	if !res {
		t.Error("It should be OnePair")
	}
	if highCard != Five {
		t.Error("highCard should be Five")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Three},
		{symbol: Spade, rank: Four},
		{symbol: Heart, rank: Five},
		{symbol: Spade, rank: Five},
	}
	res, highCard = isOnePair(cards)
	if !res {
		t.Error("It should be OnePair")
	}
	if highCard != Five {
		t.Error("highCard should be Five")
	}

	cards = []Card{
		{symbol: Spade, rank: Two},
		{symbol: Heart, rank: Five},
		{symbol: Spade, rank: Six},
		{symbol: Heart, rank: Jack},
		{symbol: Spade, rank: King},
	}
	res, _ = isOnePair(cards)
	if res {
		t.Error("It should not be OnePair")
	}
}
