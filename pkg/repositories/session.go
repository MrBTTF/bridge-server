package repositories

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/mrbttf/bridge-server/pkg/core"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

const SelectSession = `
SELECT session_id, players, deck, session_table, current_player
FROM sessions
WHERE session_id = $1
`

func (sp *SessionRepository) Get(session_id string) (core.Session, error) {
	var session core.Session
	var _deck []string
	var table []string
	err := sp.db.QueryRow(SelectSession, session_id).Scan(
		&session.Id,
		pq.Array(&session.Players),
		pq.Array(&_deck),
		pq.Array(&table),
		&session.CurrentPlayer,
	)
	session.Deck = StringToDeck(_deck)
	session.Table = StringToDeck(table)
	if err != nil {
		return core.Session{}, fmt.Errorf("Unable to get session for id %s: %w", session_id, err)
	}

	return session, nil
}

const UpsertSession = `
INSERT INTO sessions (session_id, players, deck, session_table, current_player)
VALUES($1, $2, $3, $4, $5) 
ON CONFLICT (session_id) 
WHERE session_id = $1 
DO UPDATE
SET 
players = EXCLUDED.players, 
deck = EXCLUDED.deck, 
session_table = EXCLUDED.session_table, 
current_player = EXCLUDED.current_player
`

func (sp *SessionRepository) Store(session *core.Session) error {
	_deck := DeckToString(session.Deck)
	table := DeckToString(session.Table)
	_, err := sp.db.Exec(UpsertSession,
		session.Id, pq.Array(session.Players),
		pq.Array(_deck), pq.Array(table), session.CurrentPlayer,
	)
	if err != nil {
		return fmt.Errorf("Unable to store session for id %s: %w", session.Id, err)

	}

	return nil
}

const DeleteSession = `
DELETE FROM sessions
WHERE session_id = $1
`

func (sp *SessionRepository) Delete(session_id string) error {
	_, err := sp.db.Exec(DeleteSession, session_id)
	if err != nil {
		return fmt.Errorf("Unable to delete session for id %s: %w", session_id, err)
	}

	return nil
}
