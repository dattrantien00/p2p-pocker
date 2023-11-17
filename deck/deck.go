package deck

import (
	"math/rand"
	"fmt"
	"strconv"
)

type Suit int

func (s Suit) String() string {
	return []string{"Spades", "Hearts", "Diamonds", "Club"}[s]
}

const (
	Spades   Suit = iota //bích
	Hearts               // cơ
	Diamonds             // rô
	Clubs                // chuồn
)

type Card struct {
	Suit  Suit
	Value int
}

func (c Card) String() string {
	value := strconv.Itoa(c.Value)
	if value == "1" {
		value = "ACE"
	}
	return fmt.Sprintf("%s of %s %s", value, c.Suit, suitToUnicode(c.Suit))
}

func NewCard(s Suit, value int) Card {
	if value > 13 {
		panic("the value of card cannot be higher than 13")
	}
	return Card{s, value}
}

type Deck [52]Card

func New() Deck {
	var (
		nSuits = 4
		nCard  = 13

		d = [52]Card{}
		x = 0
	)

	for i := 0; i < nSuits; i++ {
		for j := 1; j <= nCard; j++ {
			fmt.Println(j)
			d[x] = NewCard(Suit(i), j)
			x++
		}
	}
	d = shuffle(d)
	return d
}

func shuffle(d Deck) Deck{
	for i:=0;i<len(d);i++{
		r := rand.Intn(i+1)
		if i != r{
			d[i],d[r]=d[r],d[i]
		}
	}
	return d
}
func suitToUnicode(s Suit) string {
	return []string{"♠", "♥", "♦", "♣"}[s]
}
