package repositories

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/mrbttf/bridge-server/pkg/core"
)

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

const SelectPlayer = `
SELECT user_id, nickname, cards, state, session_id
FROM players
WHERE user_id = $1
`

func (pp *PlayerRepository) Get(player_id string) (core.Player, error) {
	var player core.Player
	var cards []string
	err := pp.db.QueryRow(SelectPlayer, player_id).Scan(
		&player.Id,
		&player.Nickname,
		pq.Array(&cards),
		&player.State,
		&player.SessionId,
	)
	player.Cards = StringToDeck(cards)
	if err != nil {
		return core.Player{}, fmt.Errorf("Unable to get player for id %s: %w", player_id, err)
	}

	return player, nil
}

const UpsertPlayer = `
INSERT INTO players (user_id, nickname, cards, state, state_name, session_id)
VALUES($1, $2, $3, $4, $5, $6) 
ON CONFLICT (user_id, session_id) 
WHERE user_id = $1 AND session_id = $6  
DO UPDATE
SET 
	user_id = EXCLUDED.user_id, 
	nickname = EXCLUDED.nickname, 
	cards = EXCLUDED.cards, 
	state = EXCLUDED.state, 
	state_name = EXCLUDED.state_name, 
	session_id = EXCLUDED.session_id
`

func (pp *PlayerRepository) Store(player *core.Player) error {
	cards := DeckToString(player.Cards)
	_, err := pp.db.Exec(UpsertPlayer,
		player.Id, player.Nickname, pq.Array(cards),
		player.State, player.State.String(), player.SessionId,
	)
	if err != nil {
		return err
	}

	return nil
}
