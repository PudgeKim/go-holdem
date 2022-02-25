package card

import (
	"math/rand"
	"time"
)

type Deck []Card

func NewDeck() *Deck {
	var d Deck
	symbols := []Symbol{Spade, Heart, Diamond, Clover}

	for i := 2; i < 15; i++ {
		for j := 0; j < 4; j++ {
			c := Card{
				Symbol: symbols[j],
				Rank:   Rank(i),
			}
			d = append(d, c)
		}
	}

	d.shuffle()
	return &d
}

func (d *Deck) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(*d), func(i, j int) {
		(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
	})
}

func (d *Deck) GetCard() Card {
	lastIdx := len(*d) - 1
	lastCard := (*d)[lastIdx]
	*d = (*d)[:lastIdx]
	return lastCard
}
