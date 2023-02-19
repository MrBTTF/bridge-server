package core

import (
	"errors"
	"fmt"
)

var (
	CardNotFoundError = errors.New("Card not found")
)

func Pull(session *Session, player *Player) error {
	last_idx := len(session.deck) - 1
	card := session.deck[last_idx]
	session.deck = session.deck[:last_idx]
	player.cards = append(player.cards, card)
	return nil
}

func Lay(session *Session, player *Player, card Card) error {
	cardIdx := -1
	for i, c := range player.cards {
		if c == card {
			cardIdx = i
			break
		}
	}
	if cardIdx == -1 {
		return fmt.Errorf("Couldn't lay card %s: %w", card, CardNotFoundError)

	}

	session.table = append(session.table, card)
	player.cards = append( player.cards[:cardIdx],  player.cards[cardIdx+1:]...)

	return nil
}
