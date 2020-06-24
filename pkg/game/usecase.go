package game

import (
	"fmt"

	"github.com/MrBTTF/gophercises/deck"
)

var (
	bridgeDeck []deck.Card
	strToCard  map[string]deck.Card
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
	strToCard = make(map[string]deck.Card)

	for i, suit := range suits {
		for j, rank := range ranks {
			strToCard[suit+rank] = bridgeDeck[i*(len(ranks))+j]
		}
	}
}

func InitSession(session *Session) *Session {
	newSession := session.Copy()

	_deck := createDeck()

	players := newSession.Players
	fmt.Println("InitSession")
	fmt.Println(players)
	for player := range newSession.Players {
		players[player].Laid = []deck.Card{}
		if player == session.HostPlayer {
			players[player].Hand, _deck = _deck[len(_deck)-hostHandSize:], _deck[:len(_deck)-hostHandSize]
			players[player].State = Start
			continue
		}
		players[player].Hand, _deck = _deck[len(_deck)-playerHandSize:], _deck[:len(_deck)-playerHandSize]
		players[player].Laid = []deck.Card{}
		players[player].State = NextTurn

	}
	newSession.Deck, newSession.Laid = _deck[:len(_deck)-1], _deck[len(_deck)-1:]
	players[session.HostPlayer], _ = players[session.HostPlayer].Lay(newSession.Laid[0])
	if players[session.HostPlayer].State == Normal {
		players[session.HostPlayer].State = Start
	}
	newSession.Players = players
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

func mustBeCovered(card deck.Card) bool {
	return card.Rank == deck.Six || card.Rank == deck.Ace
}

func EndTurn(session *Session, playerName string) (*Session, error) {
	newSession := session.Copy()

	player := newSession.Players[playerName]

	player, err := player.EndTurn()
	if err != nil {
		return nil, fmt.Errorf("Player %s can't end turn: %s", playerName, err)
	}

	pulledByNextPlayer := 0
	for _, card := range player.Laid {
		if card.Rank == deck.Seven {
			pulledByNextPlayer++
		}
		if card.Rank == deck.Eight {
			pulledByNextPlayer += 2
		}
	}

	newSession.PlayersOrders = append(newSession.PlayersOrders[1:], newSession.PlayersOrders[0])
	nextPlayer := newSession.NextPlayer()
	pullDeck(newSession, nextPlayer, pulledByNextPlayer)

	newSession.Laid = append(newSession.Laid, player.Laid...)
	player.Laid = []deck.Card{}
	newSession.Players[playerName] = player

	// TODO: init next player
	return newSession, nil
}

func LayCard(session *Session, playerName, cardStr, suit string) (*Session, error) {
	card, exists := strToCard[cardStr]
	if !exists {
		return nil, fmt.Errorf("Invalid card: %s", cardStr)
	}

	newSession := session.Copy()

	player := newSession.Players[playerName]
	playerHasCard := false
	for i, c := range player.Hand {
		if c == card {
			player.Hand = append(player.Hand[:i], player.Hand[i+1:]...)
			playerHasCard = true
			break
		}
	}
	if !playerHasCard {
		return nil, fmt.Errorf("Player %s doesn't have card %s", playerName, cardStr)
	}

	topCard := session.Laid[len(session.Laid)-1]

	if len(player.Laid) > 0 {
		topCard = player.Laid[len(player.Laid)-1]
	}


	if player.SuitOrdered != nil && card.Suit != *player.SuitOrdered {
		return nil, fmt.Errorf("Player %s must lay suit %s, not %s", playerName, card.Suit, player.SuitOrdered)
	} else if card.Rank != deck.Jack && (card.Rank != topCard.Rank && card.Suit != topCard.Suit) {
		return nil, fmt.Errorf("Player %s can't lay %s on %s", playerName, card, topCard)
	}

	if card.Rank == deck.Jack {
		err := orderSuit(session, suit)
		if err != nil {
			return nil, fmt.Errorf("Player %s can't lay: %s", playerName, err)
		}
	}	

	player, err := player.Lay(card)
	if err != nil {
		return nil, fmt.Errorf("Player %s can't lay: %s", playerName, err)
	}

	player.Laid = append(player.Laid, card)
	newSession.Players[playerName] = player
	return newSession, nil
}

func UnlayCard(session *Session, playerName string) (*Session, error) {
	newSession := session.Copy()

	player := newSession.Players[playerName]

	if len(player.Laid) == 0 {
		return nil, fmt.Errorf("Player %s didn't lay a card", playerName)
	}

	player, err := player.Lay(player.Laid[len(player.Laid)-1])
	if err != nil {
		return nil, fmt.Errorf("Player %s can't unlay: %s", playerName, err)
	}

	player.Hand = append(player.Hand, player.Laid[len(player.Laid)-1])
	player.Laid = player.Laid[:len(player.Laid)-1]
	newSession.Players[playerName] = player
	return newSession, nil
}

func orderSuit(session *Session, suitStr string) error {
	var suit deck.Suit
	switch suitStr {
	case "S":
		suit = deck.Spade
		break
	case "D":
		suit = deck.Diamond
		break
	case "C":
		suit = deck.Club
		break
	case "H":
		suit = deck.Heart
		break
	default:
		return fmt.Errorf("Invalid suit: %s", suitStr)
	}

	session.NextPlayer().SuitOrdered = &suit
	return nil
}

func PullDeck(session *Session, playerName string) (*Session, error) {
	newSession := session.Copy()

	player := newSession.Players[playerName]

	player, err := player.Pull()
	if err != nil {
		return nil, fmt.Errorf("Player %s can't pull: %s", playerName, err)
	}

	pullDeck(session, player, 1)

	newSession.Players[playerName] = player
	return newSession, nil
}

func pullDeck(session *Session, player *Player, count int) {
	session.Deck, player.Hand = session.Deck[:len(session.Deck)-count], append(player.Hand, session.Deck[len(session.Deck)-count:]...)
}
