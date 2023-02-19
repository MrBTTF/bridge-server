package game

import (
	"encoding/json"
	"fmt"

	"github.com/MrBTTF/gophercises/deck"
	"github.com/google/uuid"
)

const (
	hostHandSize   = 4
	playerHandSize = 5
)

type Card deck.Card

func (c Card) MarshalJSON() ([]byte, error) {
	return []byte(`"` + cardToStr[c] + `"`), nil
}

func (c *Card) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	*c = strToCard[s]
	return err
}

type Session struct {
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	HostPlayer    string             `json:"hostPlayers"`
	Players       map[string]*Player `json:"players"`
	PlayersOrders []string           `json:"playersOrders"`
	Started       bool               `json:"started"`
	Laid          []Card             `json:"laid"`
	Deck          []Card             `json:"deck"`
}

func New(name string, hostPlayer string) *Session {
	sessionId := uuid.New().String()
	return &Session{
		ID:         sessionId,
		Name:       name,
		HostPlayer: hostPlayer,
		Players: map[string]*Player{
			hostPlayer: NewHostPlayer(hostPlayer),
		},
		PlayersOrders: []string{hostPlayer},
		Started:       false,
	}
}

func (s Session) NextPlayer() *Player {
	return s.Players[s.PlayersOrders[1]]
}

func (s Session) String() string {
	result, _ := json.MarshalIndent(s, "", "\t")
	return string(result)
}

func (s Session) Copy() *Session {
	newDeck := make([]Card, len(s.Deck))
	copy(newDeck, s.Deck)

	laid := make([]Card, len(s.Laid))
	copy(laid, s.Laid)

	playersOrders := make([]string, len(s.PlayersOrders))
	copy(playersOrders, s.PlayersOrders)

	players := make(map[string]*Player)
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
		Started:       s.Started,
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

func (s State) lay(card Card) (State, error) {
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
	Name        string     `json:"name"`
	Hand        []Card     `json:"hand"`
	Laid        []Card     `json:"laid"`
	SuitOrdered *deck.Suit `json:"suitOrdered"`
	HasTurn     bool       `json:"hasTurn"`
	State       `json:"state"`
}

func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
	}
}

func NewHostPlayer(name string) *Player {
	return &Player{
		Name:    name,
		HasTurn: true,
	}
}

func (p Player) String() string {
	result, _ := json.MarshalIndent(p, "", "\t")
	return string(result)
}

func (p Player) Copy() *Player {
	hand := make([]Card, len(p.Hand))
	copy(hand, p.Hand)

	laid := make([]Card, len(p.Laid))
	copy(laid, p.Laid)

	return &Player{
		Name:    p.Name,
		Hand:    hand,
		Laid:    laid,
		State:   p.State,
		HasTurn: p.HasTurn,
	}
}

func (p Player) Lay(card Card) (*Player, error) {
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
