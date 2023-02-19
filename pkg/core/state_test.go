package core

import (
	"testing"

	"github.com/MrBTTF/gophercises/deck"
	"github.com/stretchr/testify/assert"
)

func TestStateMustLayOrPull(t *testing.T) {
	state := StateWaitForTurn

	card := NewCard(deck.Heart, deck.Ten)
	state, err := state.OnNextTurn(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateMustLayOrPull.String())
	}
	state, err = state.OnPull(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateCanLay.String())
	}

	state, err = state.OnLay(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateCanLay.String())
	}

	state, err = state.OnEndTurn(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateWaitForTurn.String())
	}
}
func TestStateMustLayOrPullFor8Ace(t *testing.T) {
	state := StateWaitForTurn

	card := NewCard(deck.Heart, deck.Eight)
	state, err := state.OnNextTurn(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateMustLayOrPull.String())
	}
	state, err = state.OnLay(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateMustLayOrPull.String())
	}

	card = NewCard(deck.Heart, deck.Ace)
	state, err = state.OnLay(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateMustLayOrPull.String())
	}

	card = NewCard(deck.Heart, deck.King)
	state, err = state.OnPull(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateCanLay.String())
	}

	card = NewCard(deck.Heart, deck.Six)
	state, err = state.OnLay(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateMustLay.String())
	}

	card = NewCard(deck.Heart, deck.King)
	state, err = state.OnPull(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateMustLay.String())
	}

	card = NewCard(deck.Heart, deck.King)
	state, err = state.OnLay(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateCanLay.String())
	}


	state, err = state.OnEndTurn(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateWaitForTurn.String())
	}
}

func TestStateMustLay(t *testing.T) {
	state := StateWaitForTurn

	card := NewCard(deck.Heart, deck.Six)
	state, err := state.OnNextTurn(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateMustLay.String())
	}
	state, err = state.OnPull(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateMustLay.String())
	}
	state, err = state.OnPull(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateMustLay.String())
	}

	card = NewCard(deck.Spade, deck.Ten)
	state, err = state.OnLay(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateCanLay.String())
	}

	state, err = state.OnEndTurn(card)
	if assert.NoError(t, err) {
		assert.Equal(t, state.String(), StateWaitForTurn.String())
	}
}
