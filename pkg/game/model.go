package game

import (
	"fmt"

	"github.com/MrBTTF/gophercises/deck"
	"github.com/google/uuid"
)

const (
	hostHandSize   = 4
	playerHandSize = 5
)

var (
	bridgeDeck []deck.Card
	strToDeck  map[string]deck.Card
)

func init() {
	bridgeDeck = deck.New(func(cards []deck.Card) []deck.Card {
		result := []deck.Card{}
		for _, card := range cards {
			if card.Rank >= deck.Six || card.Rank == deck.Ace {
				result = append(result, card)
			}
		}
		return result
	})
	ranks := []string{"A", "6", "7", "8", "9", "T", "J", "Q", "K"}
	suits := []string{"S", "D", "C", "H"}
	strToDeck = make(map[string]deck.Card)

	for i, suit := range suits {
		for j, rank := range ranks {
			strToDeck[suit+rank] = bridgeDeck[i*(len(ranks))+j]
		}
	}
}

type Session struct {
	ID         string                 `json:id`
	Name       string                 `json:name`
	HostPlayer string                 `json:hostPlayers`
	Players    map[string][]deck.Card `json:players`
	Laid       []deck.Card            `json:laid`
	Deck       []deck.Card            `json:deck`
}

func (s Session) Copy() *Session{
	newDeck := make([]deck.Card, len(s.Deck))
	copy(newDeck, s.Deck)

	laid := make([]deck.Card, len(s.Laid))
	copy(laid, s.Laid)

	players := make(map[string][]deck.Card)
	for k, v := range s.Players {
		players[k] = make([]deck.Card, len(v))
		copy(players[k], v)
	}
	return &Session{
		ID:         s.ID,
		Name:       s.Name,
		HostPlayer: s.HostPlayer,
		Players:    players,
		Deck:       newDeck,
		Laid:       laid,
	}
}

// type Player struct {
// 	Name string      `json:name`
// 	Hand []deck.Card `json:hand`
// }

func New(name string, hostPlayer string) *Session {
	return &Session{
		ID:         uuid.New().String(),
		Name:       name,
		HostPlayer: hostPlayer,
		Players: map[string][]deck.Card{
			hostPlayer: nil,
		},
	}
}

func InitSession(session *Session) *Session {
	newSession := session.Copy()
	
	_deck := createDeck()

	players := newSession.Players
	for player := range session.Players {
		if player == session.HostPlayer {
			players[player], _deck = _deck[len(_deck)-hostHandSize:], _deck[:len(_deck)-hostHandSize]
			continue
		}
		players[player], _deck = _deck[len(_deck)-playerHandSize:], _deck[:len(_deck)-playerHandSize]

	}
	return newSession
}

func createDeck() []deck.Card {
	return deck.New(func(cards []deck.Card) []deck.Card {
		result := []deck.Card{}
		for _, card := range bridgeDeck {
			result = append(result, card)
		}
		return result
	}, deck.Shuffle)
}

func LayCard(session *Session, playerName, cardStr string) (*Session, error) {
	card, exists := strToDeck[cardStr]
	if !exists {
		return nil, fmt.Errorf("Invalid card: %s", cardStr)
	}

	newSession := session.Copy()

	newPlayerCards := newSession.Players[playerName]
	for i, c := range session.Players[playerName] {
		if c == card {
			newPlayerCards = append(newPlayerCards[:i], newPlayerCards[i+1:]...)
			break
		}
	}
	if len(newPlayerCards) == len(session.Players[playerName]) {
		return nil, fmt.Errorf("Player %s doesn't have card %s", playerName, cardStr)
	}

	newSession.Laid = append(newSession.Laid, card)
	return newSession, nil
}
