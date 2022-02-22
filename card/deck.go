package card

import (
	"math/rand"
	"time"
)

type deck []Card

func NewDeck() deck {
	var d deck
	symbols := []Symbol{Spade, Heart, Diamond, Clover}

	for i := 2; i < 15; i++ {
		for j := 0; j < 4; j++ {
			c := Card{
				symbol: symbols[j],
				rank:   Rank(i),
			}
			d = append(d, c)
		}
	}

	d.shuffle()
	return d
}

func (d *deck) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(*d), func(i, j int) {
		(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
	})
}

func (d *deck) GetCard() Card {
	lastIdx := len(*d) - 1
	lastCard := (*d)[lastIdx]
	*d = (*d)[:lastIdx]
	return lastCard
}
