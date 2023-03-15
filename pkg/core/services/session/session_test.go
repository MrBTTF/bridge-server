package session

import (
	"errors"
	"testing"

	"github.com/MrBTTF/gophercises/deck"
	"github.com/mrbttf/bridge-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

const (
	player_id = "test_player"
	room_id   = "test_room"
)

var (
	NotFoundError = errors.New("Not found")
)

type MockSessionRepository struct {
	sessions map[string]core.Session
}

func NewMockSessionRepository() *MockSessionRepository {

	return &MockSessionRepository{
		sessions: map[string]core.Session{},
	}
}

func (m *MockSessionRepository) Get(session_id string) (core.Session, error) {
	v, ok := m.sessions[session_id]
	if !ok {
		return core.Session{}, NotFoundError
	}
	return v, nil
}

func (m *MockSessionRepository) Store(session *core.Session) error {
	m.sessions[session.Id] = *session
	return nil
}

type MockPlayerRepository struct {
	players map[string]core.Player
}

func NewMockPlayerRepository() *MockPlayerRepository {

	return &MockPlayerRepository{
		players: map[string]core.Player{},
	}
}

func (m *MockPlayerRepository) Get(player_id string) (core.Player, error) {
	v, ok := m.players[player_id]
	if !ok {
		return core.Player{}, NotFoundError
	}
	return v, nil
}

func (m *MockPlayerRepository) Store(player *core.Player) error {
	m.players[player.Id] = *player
	return nil
}

type MockUserRepository struct {
	users map[string]core.User
}

func NewMockUserRepository() *MockUserRepository {

	return &MockUserRepository{
		users: map[string]core.User{},
	}
}

func (m *MockUserRepository) Get(user_id string) (core.User, error) {
	v, ok := m.users[user_id]
	if !ok {
		return core.User{}, NotFoundError
	}
	return v, nil
}

func (m *MockUserRepository) GetByEmail(email string) (core.User, error) {
	return m.users[email], nil
}

func (m *MockUserRepository) Store(player *core.User) error {
	m.users[player.Id] = *player
	return nil
}

func TestSession(t *testing.T) {
	sessions := NewMockSessionRepository()
	players := NewMockPlayerRepository()
	users := NewMockUserRepository()
	err := users.Store(&core.User{
		Id: player_id,
	})
	if err != nil {
		panic(err)
	}

	_deck := core.NewDeck()
	tableCard := core.NewCard(deck.Diamond, deck.Queen)
	playerCard := core.NewCard(deck.Heart, deck.Queen)
	setLastCards(_deck, tableCard, playerCard)
	session_service := New(sessions, players, users)
	session_id, err := session_service.Create([]string{player_id}, _deck)
	if err != nil {
		panic(err)
	}
	session_service.NextTurn(session_id, player_id)

	err = session_service.Pull(session_id, player_id)
	if err != nil {
		panic(err)
	}
	session, err := sessions.Get(session_id)
	if err != nil {
		panic(err)
	}
	player, err := players.Get(player_id)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	assert.NotContains(t, session.Deck, tableCard)
	assert.Contains(t, player.Cards, playerCard)

	err = session_service.Lay(session_id, player_id, playerCard)
	if err != nil {
		panic(err)
	}
	session, err = sessions.Get(session_id)
	if err != nil {
		panic(err)
	}
	player, err = players.Get(player_id)
	if err != nil {
		panic(err)
	}
	assert.Contains(t, session.Table, tableCard)
	assert.NotContains(t, player.Cards, playerCard)
}

func setLastCards(_deck []deck.Card, tableCard, playerCard deck.Card) {
	for i, card := range _deck {
		if card == tableCard {
			_deck[i] = _deck[len(_deck)-1]
		} else if card == playerCard {
			_deck[i] = _deck[len(_deck)-2]
		}
	}
	_deck[len(_deck)-1] = tableCard
	_deck[len(_deck)-2] = playerCard
}
