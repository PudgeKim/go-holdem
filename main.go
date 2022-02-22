package main

import (
	"fmt"
	"github.com/PudgeKim/card"
)

func main() {
	deck := card.NewDeck()
	fmt.Println(deck)
	card := deck.GetCard()
	fmt.Println(card)
	fmt.Println(deck)

}
