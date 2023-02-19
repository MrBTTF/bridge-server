package core

import (
	"testing"

	"github.com/MrBTTF/gophercises/deck"
	"github.com/stretchr/testify/assert"
)

func TestPullAndLay(t *testing.T) {
	session := &Session{
		deck: deck.New(deck.Filter(func(card deck.Card) bool {
			return card.Rank < deck.Six && card.Rank != deck.Ace
		}), deck.Shuffle),
	}
	player := &Player{}

	err := Pull(session, player)
	card := player.cards[len(player.cards)-1]
	if assert.NoError(t, err) {
		assert.NotContains(t, session.deck, card)
		assert.Contains(t, player.cards, card)
	}

	err = Lay(session, player, card)
	if assert.NoError(t, err){
		assert.Contains(t, session.table, card)
		assert.NotContains(t, player.cards, card)
	}

}
