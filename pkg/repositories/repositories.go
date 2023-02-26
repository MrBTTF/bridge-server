package repositories

import (
	"github.com/MrBTTF/gophercises/deck"
	"github.com/mrbttf/bridge-server/pkg/core"
	"golang.org/x/exp/slices"
)

var ranks = []string{
	"", "A", "2", "3", "4", "5",
	"6", "7", "8", "9", "T", "J", "Q", "K",
}
var suits = []string{"S", "D", "C", "H"}

func CardToString(card deck.Card) string {
	return suits[card.Suit] + ranks[card.Rank]
}

func DeckToString(_deck []deck.Card) []string {
	result := make([]string, 0, len(_deck))
	for _, card := range _deck {
		result = append(result, CardToString(card))
	}
	return result
}

func StringToCard(card string) deck.Card {
	suit := string(card[0])
	rank := string(card[1])
	suit_idx := slices.Index(suits, suit)
	rank_idx := slices.Index(ranks, rank)
	return core.NewCard(deck.Suit(suit_idx), deck.Rank(rank_idx))
}

func StringToDeck(_deck []string) []deck.Card {
	result := make([]deck.Card, 0, len(_deck))
	for _, card := range _deck {
		result = append(result, StringToCard(card))
	}
	return result
}
