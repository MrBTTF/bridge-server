package session

import (
	"errors"
	"fmt"

	"github.com/MrBTTF/gophercises/deck"
	"github.com/google/uuid"
	"github.com/mrbttf/bridge-server/pkg/core"
)

var (
	CardNotFoundError            = errors.New("Card not found")
	PlayerInSessionNotFoundError = errors.New("Player not found in session")
)

type SessionService struct {
	sessions core.SessionRepository
	players  core.PlayerRepository
}

func New(sessions core.SessionRepository, players core.PlayerRepository) *SessionService {
	return &SessionService{
		sessions: sessions,
		players:  players,
	}
}

func (s *SessionService) Create(players []string, _deck []deck.Card) (string, error) {
	last_idx := len(_deck) - 1
	card := _deck[last_idx]
	_deck = _deck[:last_idx]
	table := []core.Card{card}

	session_id := uuid.New().String()
	session := &core.Session{
		Id:            session_id,
		Players:       players,
		Deck:          _deck,
		Table:         table,
		CurrentPlayer: players[0],
	}
	err := s.sessions.Store(session)
	if err != nil {
		return "", fmt.Errorf("Unable to create session: %w", err)
	}
	return session_id, nil
}

func (s *SessionService) Pull(session_id, player_id string) error {
	session, err := s.sessions.Get(session_id)
	if err != nil {
		return fmt.Errorf("Unable to pull for session %s, player %s: %w", session_id, player_id, err)
	}
	if !session.HasPlayer(player_id) {
		return fmt.Errorf("Unable to pull for session %s, player %s: %w", session_id, player_id, PlayerInSessionNotFoundError)
	}

	last_idx := len(session.Deck) - 1
	card := session.Deck[last_idx]
	session.Deck = session.Deck[:last_idx]
	player, err := s.players.Get(player_id)
	if err != nil {
		return fmt.Errorf("Unable to pull for session %s, player %s: %w", session_id, player_id, err)
	}

	player.State, err = player.State.OnPull(card)
	if err != nil {
		return fmt.Errorf("Unable to pull for session %s, player %s, card %s: %w", session_id, player_id, card, err)
	}

	player.Cards = append(player.Cards, card)
	err = s.sessions.Store(&session)
	if err != nil {
		return fmt.Errorf("Unable to pull for session %s, player %s: %w", session_id, player_id, err)
	}
	err = s.players.Store(&player)
	if err != nil {
		return fmt.Errorf("Unable to pull for session %s, player %s: %w", session_id, player_id, err)
	}
	return nil
}

func (s *SessionService) Lay(session_id, player_id string, card core.Card) error {
	session, err := s.sessions.Get(session_id)
	if err != nil {
		return fmt.Errorf("Unable to lay for session %s, player %s, card %s: %w", session_id, player_id, card, err)
	}
	if !session.HasPlayer(player_id) {
		return fmt.Errorf("Unable to lay for session %s, player %s, card %s: %w", session_id, player_id, card, PlayerInSessionNotFoundError)
	}
	player, err := s.players.Get(player_id)
	if err != nil {
		return fmt.Errorf("Unable to lay for session %s, player %s, card %s: %w", session_id, player_id, card, err)
	}

	cardIdx := -1
	for i, c := range player.Cards {
		if c == card {
			cardIdx = i
			break
		}
	}
	if cardIdx == -1 {
		return fmt.Errorf("Unable to lay for session %s, player %s, card %s: %w", session_id, player_id, card, CardNotFoundError)
	}
	err = layCardOnTable(session.Table, card)
	if err != nil {
		return fmt.Errorf("Unable to lay for session %s, player %s, card %s: %w", session_id, player_id, card, err)
	}

	session.Table = append(session.Table, card)
	player.Cards = append(player.Cards[:cardIdx], player.Cards[cardIdx+1:]...)

	player.State, err = player.State.OnLay(card)
	if err != nil {
		return fmt.Errorf("Unable to lay for session %s, player %s, card %s: %w", session_id, player_id, card, err)
	}

	err = s.sessions.Store(&session)
	if err != nil {
		return fmt.Errorf("Unable to lay for session %s, player %s, card %s: %w", session_id, player_id, card, err)
	}
	err = s.players.Store(&player)
	if err != nil {
		return fmt.Errorf("Unable to pull for session %s, player %s: %w", session_id, player_id, err)
	}
	return nil
}

func (s *SessionService) NextTurn(session_id, player_id string) error {
	session, err := s.sessions.Get(session_id)
	if err != nil {
		return fmt.Errorf("Unable to change turn for session %s, player %s: %w", session_id, player_id, err)
	}
	if !session.HasPlayer(player_id) {
		return fmt.Errorf("Unable to change turn for session %s, player %s: %w", session_id, player_id, PlayerInSessionNotFoundError)
	}
	session.CurrentPlayer = player_id

	player, err := s.players.Get(player_id)
	if err != nil {
		return fmt.Errorf("Unable to change turn for session %s, player %s: %w", session_id, player_id, PlayerInSessionNotFoundError)
	}

	topCard := session.Table[len(session.Table)-1]
	player.State, err = player.State.OnNextTurn(topCard)
	if err != nil {
		return fmt.Errorf("Unable to change turn for session %s, player %s: %w", session_id, player_id, err)
	}

	
	err = s.sessions.Store(&session)
	if err != nil {
		return fmt.Errorf("Unable to change turn for session %s, player %s: %w", session_id, player_id, err)
	}
	err = s.players.Store(&player)
	if err != nil {
		return fmt.Errorf("Unable to change turn for session %s, player %s: %w", session_id, player_id, err)
	}
	return nil
}

func layCardOnTable(table []deck.Card, card deck.Card) error {
	if len(table) == 0 {
		return nil
	}
	topCard := table[len(table)-1]
	if topCard.Rank == deck.Jack || card.Rank == deck.Jack {
		return nil
	} else if topCard.Rank == card.Rank || topCard.Suit == card.Suit {
		return nil
	}
	return fmt.Errorf("Cannot lay %s on %s", card, topCard)
}
