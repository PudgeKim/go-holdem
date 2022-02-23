package card

import (
	"fmt"
	"reflect"
	"testing"
)

var cards []Card
var res bool
var highCard Rank

func TestMakeAllCombinations(t *testing.T) {
	cards = []Card{
		{Symbol: Spade, Rank: Queen},
		{Symbol: Heart, Rank: Six},
		{Symbol: Diamond, Rank: Ace},
		{Symbol: Spade, Rank: King},
		{Symbol: Spade, Rank: Six},
		{Symbol: Clover, Rank: Three},
		{Symbol: Clover, Rank: Two},
	}

	var allCombs [][]Card

	makeAllCombinations(cards, []Card{}, &allCombs, 0, 0)

	if len(allCombs) != 21 {
		t.Error("All combination's length should be 21 (7C5)")
	}
}

func TestGetBestHandsRank(t *testing.T) {
	cards = []Card{
		{Symbol: Spade, Rank: Queen},
		{Symbol: Heart, Rank: Six},
		{Symbol: Diamond, Rank: Ace},
		{Symbol: Spade, Rank: King},
		{Symbol: Spade, Rank: Six},
		{Symbol: Clover, Rank: Three},
		{Symbol: Clover, Rank: Two},
	}
	bestCards, bestHandsRank, bestHighCard := getBestHandsRank(cards)
	fmt.Println(bestCards) // to delete
	expected := []Card{
		{Symbol: Heart, Rank: Six},
		{Symbol: Spade, Rank: Six},
		{Symbol: Spade, Rank: Queen},
		{Symbol: Spade, Rank: King},
		{Symbol: Diamond, Rank: Ace},
	}
	if !reflect.DeepEqual(bestCards, expected) {
		t.Error("BestCard is wrong")
	}
	if bestHandsRank != OnePair {
		t.Error("HandsRank should be OnePair")
	}
	if bestHighCard != Six {
		t.Error("HighCard should be Six")
	}

}

func TestRoyalStraightFlush(t *testing.T) {
	cards = []Card{
		{Symbol: Spade, Rank: Ten},
		{Symbol: Spade, Rank: Jack},
		{Symbol: Spade, Rank: Queen},
		{Symbol: Spade, Rank: King},
		{Symbol: Spade, Rank: Ace},
	}
	res, highCard = isRoyalStraightFlush(cards)
	if !res {
		t.Error("It should be RoyalStraightFlush but wrong result")
	}
	if highCard != Ace {
		t.Error("RoyalStraightFlush must be end with Ace")
	}

	cards2 := []Card{
		{Symbol: Spade, Rank: Nine},
		{Symbol: Spade, Rank: Ten},
		{Symbol: Spade, Rank: Jack},
		{Symbol: Spade, Rank: Queen},
		{Symbol: Spade, Rank: King},
	}
	res, highCard = isRoyalStraightFlush(cards2)
	if res {
		t.Error("RoyalStraightFlush's highCard should be Ace")
	}
}

func TestStraightFlush(t *testing.T) {
	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Spade, Rank: Three},
		{Symbol: Spade, Rank: Four},
		{Symbol: Spade, Rank: Five},
		{Symbol: Spade, Rank: Ace},
	}
	res, highCard = isStraightFlush(cards)
	if !res {
		t.Error("It should be StraightFlush")
	}
	if highCard != Five {
		t.Error("highCard should be Five")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Four},
		{Symbol: Spade, Rank: Five},
		{Symbol: Spade, Rank: Six},
		{Symbol: Spade, Rank: Seven},
		{Symbol: Spade, Rank: Eight},
	}
	res, highCard = isStraightFlush(cards)
	if !res {
		t.Error("It should be StraightFlush")
	}
	if highCard != Eight {
		t.Error("highCard should be Eight")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Four},
		{Symbol: Spade, Rank: Five},
		{Symbol: Spade, Rank: Seven},
		{Symbol: Spade, Rank: Eight},
		{Symbol: Spade, Rank: Ten},
	}
	res, _ = isStraightFlush(cards)
	if res {
		t.Error("It should not be StraightFlush")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Four},
		{Symbol: Heart, Rank: Five},
		{Symbol: Spade, Rank: Seven},
		{Symbol: Spade, Rank: Eight},
		{Symbol: Spade, Rank: Ten},
	}
	res, _ = isStraightFlush(cards)
	if res {
		t.Error("It should not be StraightFlush")
	}
}

func TestFourCard(t *testing.T) {
	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Two},
		{Symbol: Spade, Rank: Two},
		{Symbol: Spade, Rank: Two},
		{Symbol: Spade, Rank: Ten},
	}
	res, highCard = isFourCard(cards)
	if !res {
		t.Error("It should be FourCard")
	}
	if highCard != Two {
		t.Error("highCard should be Two")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Three},
		{Symbol: Spade, Rank: Three},
		{Symbol: Spade, Rank: Three},
	}
	res, highCard = isFourCard(cards)
	if !res {
		t.Error("It should be FourCard")
	}
	if highCard != Three {
		t.Error("highCard should be Three")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Four},
		{Symbol: Spade, Rank: Six},
		{Symbol: Spade, Rank: Ten},
	}
	res, _ = isFourCard(cards)
	if res {
		t.Error("It should not be FourCard")
	}
}

func TestFullHouse(t *testing.T) {
	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Two},
		{Symbol: Spade, Rank: Two},
		{Symbol: Spade, Rank: Three},
		{Symbol: Spade, Rank: Three},
	}
	res, highCard = isFullHouse(cards)
	if !res {
		t.Error("It should be FullHouse")
	}
	if highCard != Two {
		t.Error("highChard should be Two")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Two},
		{Symbol: Spade, Rank: Three},
		{Symbol: Spade, Rank: Three},
		{Symbol: Spade, Rank: Three},
	}
	res, highCard = isFullHouse(cards)
	if !res {
		t.Error("It should be FullHouse")
	}
	if highCard != Three {
		t.Error("highCard should be Three")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Three},
		{Symbol: Spade, Rank: Four},
		{Symbol: Spade, Rank: Five},
	}
	res, _ = isFullHouse(cards)
	if res {
		t.Error("It should not be FullHouse")
	}
}

func TestFlush(t *testing.T) {
	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Spade, Rank: Two},
		{Symbol: Spade, Rank: Two},
		{Symbol: Spade, Rank: Three},
		{Symbol: Spade, Rank: Three},
	}
	res, highCard = isFlush(cards)
	if !res {
		t.Error("It should be Flush")
	}
	if highCard != Three {
		t.Error("highCard should be Three")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Two},
		{Symbol: Spade, Rank: Two},
		{Symbol: Spade, Rank: Three},
		{Symbol: Spade, Rank: Three},
	}
	res, _ = isFlush(cards)
	if res {
		t.Error("It should not be Flush")
	}
}

func TestStraight(t *testing.T) {
	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Four},
		{Symbol: Heart, Rank: Five},
		{Symbol: Spade, Rank: Six},
	}
	res, highCard = isStraight(cards)
	if !res {
		t.Error("It should be Straight")
	}
	if highCard != Six {
		t.Error("highCard should be Six")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Four},
		{Symbol: Spade, Rank: Five},
		{Symbol: Heart, Rank: Ace},
	}
	res, highCard = isStraight(cards)
	if !res {
		t.Error("It should be Straight")
	}
	if highCard != Five {
		t.Error("highCard should be Five")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Four},
		{Symbol: Spade, Rank: King},
		{Symbol: Heart, Rank: Ace},
	}
	res, _ = isStraight(cards)
	if res {
		t.Error("It should not be Straight")
	}
}

func TestTriple(t *testing.T) {
	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Two},
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Five},
		{Symbol: Spade, Rank: Six},
	}
	res, highCard = isTriple(cards)
	if !res {
		t.Error("It should be Triple")
	}
	if highCard != Two {
		t.Error("highCard should be Two")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Five},
		{Symbol: Spade, Rank: Five},
		{Symbol: Heart, Rank: Five},
		{Symbol: Spade, Rank: Six},
	}
	res, highCard = isTriple(cards)
	if !res {
		t.Error("It should be Triple")
	}
	if highCard != Five {
		t.Error("highCard should be Five")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Six},
		{Symbol: Heart, Rank: Six},
		{Symbol: Spade, Rank: Six},
	}
	res, highCard = isTriple(cards)
	if !res {
		t.Error("It should be Triple")
	}
	if highCard != Six {
		t.Error("highCard should be Six")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Six},
		{Symbol: Heart, Rank: Ten},
		{Symbol: Spade, Rank: Jack},
	}
	res, _ = isTriple(cards)
	if res {
		t.Error("It should not be Triple")
	}
}

func TestTwoPair(t *testing.T) {
	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Two},
		{Symbol: Spade, Rank: Three},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Six},
	}
	res, highCard = isTwoPair(cards)
	if !res {
		t.Error("It should be TwoPair")
	}
	if highCard != Three {
		t.Error("highCard should be Three")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Three},
		{Symbol: Heart, Rank: Six},
		{Symbol: Spade, Rank: Six},
	}
	res, highCard = isTwoPair(cards)
	if !res {
		t.Error("It should be TwoPair")
	}
	if highCard != Six {
		t.Error("highCard should be Six")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Three},
		{Symbol: Heart, Rank: Five},
		{Symbol: Spade, Rank: Six},
	}
	res, _ = isTwoPair(cards)
	if res {
		t.Error("It should not be TwoPair")
	}
}

func TestOnePair(t *testing.T) {
	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Two},
		{Symbol: Spade, Rank: Three},
		{Symbol: Heart, Rank: Five},
		{Symbol: Spade, Rank: Six},
	}
	res, highCard = isOnePair(cards)
	if !res {
		t.Error("It should be OnePair")
	}
	if highCard != Two {
		t.Error("highCard should be Two")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Three},
		{Symbol: Heart, Rank: Five},
		{Symbol: Spade, Rank: Six},
	}
	res, highCard = isOnePair(cards)
	if !res {
		t.Error("It should be OnePair")
	}
	if highCard != Three {
		t.Error("highCard should be Three")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Five},
		{Symbol: Heart, Rank: Five},
		{Symbol: Spade, Rank: Six},
	}
	res, highCard = isOnePair(cards)
	if !res {
		t.Error("It should be OnePair")
	}
	if highCard != Five {
		t.Error("highCard should be Five")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Three},
		{Symbol: Spade, Rank: Four},
		{Symbol: Heart, Rank: Five},
		{Symbol: Spade, Rank: Five},
	}
	res, highCard = isOnePair(cards)
	if !res {
		t.Error("It should be OnePair")
	}
	if highCard != Five {
		t.Error("highCard should be Five")
	}

	cards = []Card{
		{Symbol: Spade, Rank: Two},
		{Symbol: Heart, Rank: Five},
		{Symbol: Spade, Rank: Six},
		{Symbol: Heart, Rank: Jack},
		{Symbol: Spade, Rank: King},
	}
	res, _ = isOnePair(cards)
	if res {
		t.Error("It should not be OnePair")
	}
}
