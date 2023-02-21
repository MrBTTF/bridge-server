package core

import (
	"github.com/MrBTTF/gophercises/deck"
	"github.com/mrbttf/bridge-server/pkg/core/state"
)

type Card = deck.Card

func NewCard(suit deck.Suit, rank deck.Rank) deck.Card {
	return Card{
		Suit: suit,
		Rank: rank,
	}
}

func NewDeck() []deck.Card{
	return deck.New(deck.Filter(func(card deck.Card) bool {
		return card.Rank < deck.Six && card.Rank != deck.Ace
	}), deck.Shuffle)
}

type Player struct {
	Id    string
	Name  string
	Cards []Card
	State state.State
}

type Session struct {
	Id      string
	Players []string
	Deck    []Card
	Table   []Card
	CurrentPlayer string
}

func (s Session) HasPlayer(player_id string) bool {
	for _, p := range s.Players {
		if p == player_id {
			return true
		}
	}
	return false
}
