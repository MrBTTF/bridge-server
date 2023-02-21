package session

import (
	"testing"

	"github.com/mrbttf/bridge-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

const (
	player_id = "test_player"
	room_id   = "test_room"
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
	return m.sessions[session_id], nil
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
	return m.players[player_id], nil
}

func (m *MockPlayerRepository) Store(player *core.Player) error {
	m.players[player.Id] = *player
	return nil
}

func TestSession(t *testing.T) {
	sessions := NewMockSessionRepository()
	players := NewMockPlayerRepository()
	err := players.Store(&core.Player{
		Id: player_id,
	})
	if err != nil {
		panic(err)
	}

	_deck := core.NewDeck()
	tableCard := _deck[len(_deck)-1]
	playerCard := _deck[len(_deck)-2]
	session_service := New(sessions, players)
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
