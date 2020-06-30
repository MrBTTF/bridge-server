package game

import (
	"fmt"
	"strings"

	"github.com/MrBTTF/gophercises/deck"
	"github.com/google/uuid"
)

const (
	hostHandSize   = 4
	playerHandSize = 5
)

type Session struct {
	ID            string             `json:id`
	Name          string             `json:name`
	HostPlayer    string             `json:hostPlayers`
	Players       map[string]*Player `json:players`
	PlayersOrders []string           `json:playersOrders`
	Laid          []deck.Card        `json:laid`
	Deck          []deck.Card        `json:deck`
}

func New(name string, hostPlayer string) *Session {
	return &Session{
		ID:         uuid.New().String(),
		Name:       name,
		HostPlayer: hostPlayer,
		Players: map[string]*Player{
			hostPlayer: {
				Name: hostPlayer,
			},
		},
		PlayersOrders: []string{hostPlayer},
	}
}

func (s Session) NextPlayer() *Player {
	return s.Players[s.PlayersOrders[1]]
}

func (s Session) String() string {
	var result strings.Builder
	fmt.Fprintf(&result, "ID: %s\n", s.ID)
	fmt.Fprintf(&result, "Name: %s\n", s.Name)
	fmt.Fprintf(&result, "Host Player: %s\n", s.HostPlayer)
	fmt.Fprintf(&result, "Laid: %s\n", s.Laid)
	fmt.Fprintf(&result, "Players Orders: %s\n", s.PlayersOrders)
	for _, player := range s.Players {
		fmt.Fprintf(&result, "Player: %s\n", player.Name)
		fmt.Fprint(&result, "\tHand: ")
		for _, card := range player.Hand {
			fmt.Fprintf(&result, "%s, ", card)
		}
		fmt.Fprintln(&result)
		fmt.Fprint(&result, "\tLaid: ")
		for _, card := range player.Laid {
			fmt.Fprintf(&result, "%s, ", card)
		}
		fmt.Fprintln(&result)
		fmt.Fprintf(&result, "\tState: %s\n", player.State)
	}
	return result.String()
}

func (s Session) Copy() *Session {
	newDeck := make([]deck.Card, len(s.Deck))
	copy(newDeck, s.Deck)

	laid := make([]deck.Card, len(s.Laid))
	copy(laid, s.Laid)

	playersOrders := make([]string, len(s.PlayersOrders))
	copy(playersOrders, s.PlayersOrders)

	players := make(map[string]*Player)
	fmt.Printf("%+v\n", s.Players)
	for name, player := range s.Players {
		players[name] = player.Copy()
	}
	return &Session{
		ID:            s.ID,
		Name:          s.Name,
		HostPlayer:    s.HostPlayer,
		Players:       players,
		PlayersOrders: playersOrders,
		Deck:          newDeck,
		Laid:          laid,
	}
}

//go:generate stringer -type=State
type State uint8

const (
	Start State = iota
	Normal
	FirstPull
	ShouldLay
	CanPull
	MustLay
	NextTurn
)

type InvalidActionError struct {
	state  State
	action string
}

func (er InvalidActionError) Error() string {
	return fmt.Sprintf("Invalid action \"%s\" for state %s", er.action, er.state)
}

func (s State) lay(card deck.Card) (State, error) {
	if s != NextTurn {
		switch card.Rank {
		case deck.Ace:
			fallthrough
		case deck.Eight:
			return ShouldLay, nil
		case deck.Six:
			return MustLay, nil
		default:
			return Normal, nil
		}
	}
	return s, &InvalidActionError{s, "lay"}
}

func (s State) pull() (State, error) {
	switch s {
	case Start:
		return FirstPull, nil
	case ShouldLay:
		return CanPull, nil
	case MustLay:
		return MustLay, nil
	}
	return s, &InvalidActionError{s, "pull"}
}

func (s State) end() (State, error) {
	if s == FirstPull || s == Normal || s == CanPull {
		return NextTurn, nil
	}
	return s, &InvalidActionError{s, "end"}
}

type Player struct {
	Name        string      `json:name`
	Hand        []deck.Card `json:hand`
	Laid        []deck.Card `json:laid`
	SuitOrdered *deck.Suit  `json:suitOrdered`
	State       `json:state`
}

func (p Player) Copy() *Player {
	hand := make([]deck.Card, len(p.Hand))
	copy(hand, p.Hand)

	laid := make([]deck.Card, len(p.Laid))
	copy(laid, p.Laid)

	return &Player{
		Name:  p.Name,
		Hand:  hand,
		Laid:  laid,
		State: p.State,
	}
}

func (p Player) Lay(card deck.Card) (*Player, error) {
	player := p.Copy()
	State, err := player.State.lay(card)
	if err != nil {
		return nil, err
	}
	player.State = State
	return player, nil
}

func (p Player) Pull() (*Player, error) {
	player := p.Copy()
	State, err := player.State.pull()
	if err != nil {
		return nil, err
	}
	player.State = State
	return player, nil
}

func (p Player) EndTurn() (*Player, error) {
	player := p.Copy()
	State, err := player.State.end()
	if err != nil {
		return nil, err
	}
	player.State = State
	return player, nil
}
