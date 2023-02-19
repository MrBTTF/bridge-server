package core

import (
	"fmt"

	"github.com/MrBTTF/gophercises/deck"
)

type Card = deck.Card

func NewCard(suit deck.Suit, rank deck.Rank) Card {
	return Card{
		Suit: suit,
		Rank: rank,
	}
}

var mustLayCards = []deck.Rank{deck.Six, deck.Eight, deck.Ace}

func ifMustLay(card Card) bool {
	return card.Rank == deck.Six
}

func ifMustLayOrPull(card Card) bool {
	return card.Rank == deck.Eight || card.Rank == deck.Ace
}

//go:generate stringer -type=Action
type Action uint8

const (
	ActionNone Action = iota
	ActionLay
	ActionPull
	ActionEndTurn
	ActionNextTurn
)

//go:generate stringer -type=State
type State uint8

const (
	StateWaitForTurn State = iota
	StateMustLayOrPull
	StateMustLay
	StateCanLay
)

// func (s State) Next(action Action, card Card) (State, error) {
// 	switch s {
// 	case StateWaitForTurn:
// 		if action == ActionNextTurn {
// 			if ifMustLay(card) {
// 				return StateMustLay, nil
// 			} else {
// 				return StateMustLayOrPull, nil
// 			}
// 		}
// 	case StateMustLayOrPull:
// 		switch action {
// 		case ActionLay:
// 			if ifMustLay(card) {
// 				return StateMustLay, nil
// 			} else if ifMustLayOrPull(card) {
// 				return StateMustLayOrPull, nil
// 			} else {
// 				return StateCanLay, nil
// 			}
// 		case ActionPull:
// 			return StateCanLay, nil
// 		}
// 	case StateMustLay:
// 		switch action {
// 		case ActionLay:
// 			if ifMustLay(card) {
// 				return StateMustLay, nil
// 			} else if ifMustLayOrPull(card) {
// 				return StateMustLayOrPull, nil
// 			} else {
// 				return StateCanLay, nil
// 			}
// 		case ActionPull:
// 			return StateMustLay, nil
// 		}
// 	case StateCanLay:
// 		switch action {
// 		case ActionLay:
// 			if ifMustLay(card) {
// 				return StateMustLay, nil
// 			} else if ifMustLayOrPull(card) {
// 				return StateMustLayOrPull, nil
// 			} else {
// 				return StateCanLay, nil
// 			}
// 		case ActionEndTurn:
// 			return StateWaitForTurn, nil
// 		}
// 	}
// 	err := fmt.Errorf("Cannot move from state %s: action %s, card %s", s, action, card)
// 	return s, err
// }

func (s State) OnNextTurn(card Card) (State, error) {
	if ifMustLay(card) {
		return StateMustLay, nil
	}
	return StateMustLayOrPull, nil
}

func (s State) OnLay(card Card) (State, error) {
	switch s {
	case StateMustLayOrPull:
		fallthrough
	case StateMustLay:
		fallthrough
	case StateCanLay:
		if ifMustLay(card) {
			return StateMustLay, nil
		} else if ifMustLayOrPull(card) {
			return StateMustLayOrPull, nil
		} else {
			return StateCanLay, nil
		}
	}
	err := fmt.Errorf("Cannot move from state %s: action ActionLay, card %s", s, card)
	return s, err
}

func (s State) OnPull(card Card) (State, error) {
	switch s {
	case StateMustLayOrPull:
		return StateCanLay, nil
	case StateMustLay:
		return StateMustLay, nil
	}
	err := fmt.Errorf("Cannot move from state %s: action ActionPull, card %s", s, card)
	return s, err
}

func (s State) OnEndTurn(card Card) (State, error) {
	if s == StateCanLay {
		return StateWaitForTurn, nil
	}
	err := fmt.Errorf("Cannot move from state %s: action EndTurn, card %s", s, card)
	return s, err
}
