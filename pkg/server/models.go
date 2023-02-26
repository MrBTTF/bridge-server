package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/mrbttf/bridge-server/pkg/core"
	"github.com/mrbttf/bridge-server/pkg/repositories"
)

type sessionGetRequest struct {
	SessionId string `json:"session_id" example:"string"`
	DefaultRequest
}

type sessionCreateRequest struct {
	PlayerIds []string `json:"player_ids" example:"string"`
	DefaultRequest
}

type sessionLayRequest struct {
	SessionId string `json:"session_id" example:"string"`
	PlayerId string `json:"player_id" example:"string"`
	Card     string `json:"card" example:"string"`
	DefaultRequest
}

type sessionPullRequest struct {
	SessionId string `json:"session_id" example:"string"`
	PlayerId string `json:"player_id" example:"string"`
	DefaultRequest
}

type sessionNextTurnRequest struct {
	SessionId string `json:"session_id" example:"string"`
	PlayerId string `json:"player_id" example:"string"`
	DefaultRequest
}

type PlayerResponse struct {
	Id        string
	Name      string
	Cards     []string
	State     string
	SessionId string
}

type SessionResponse struct {
	Id            string
	Players       []string
	Deck          []string
	Table         []string
	CurrentPlayer PlayerResponse
}

func NewSessionResponse(session *core.Session, player *core.Player) *SessionResponse {
	return &SessionResponse{
		Id:      session.Id,
		Deck:    repositories.DeckToString(session.Deck),
		Table:   repositories.DeckToString(session.Table),
		Players: session.Players,
		CurrentPlayer: PlayerResponse{
			Id:        player.Id,
			Name:      player.Name,
			Cards:     repositories.DeckToString(player.Cards),
			State:     player.State.String(),
			SessionId: player.SessionId,
		},
	}
}

type sessionGetResponse struct {
	Session SessionResponse `json:"session" example:SessionResponse`
	DefaultResponse
}

type sessionCreateResponse struct {
	SessionID string `json:"session_id" example:"string"`
	DefaultResponse
}

type ErrResponse struct {
	Code int `json:"-"`

	Success bool   `json:"success" example:"false"`
	Message string `json:"message,omitempty"`
}

func (er ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	er.Success = false
	render.Status(r, er.Code)
	return nil
}

type DefaultRequest struct {
}

func (req *DefaultRequest) Bind(r *http.Request) error {
	return nil
}

type DefaultResponse struct {
	Success bool `json:"success" example:"true"`
}

func (resp *DefaultResponse) Render(w http.ResponseWriter, r *http.Request) error {
	resp.Success = true
	return nil
}
