package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/mrbttf/bridge-server/pkg/core"
	"github.com/mrbttf/bridge-server/pkg/repositories"
)

type sessionGetRequest struct {
	SessionId string `json:"session_id" example:"string"`
	AuthRequest
}

type sessionCreateRequest struct {
	PlayerIds []string `json:"player_ids" example:"string"`
	AuthRequest
}

type sessionLayRequest struct {
	SessionId string `json:"session_id" example:"string"`
	PlayerId  string `json:"player_id" example:"string"`
	Card      string `json:"card" example:"string"`
	AuthRequest
}

type sessionPullRequest struct {
	SessionId string `json:"session_id" example:"string"`
	PlayerId  string `json:"player_id" example:"string"`
	AuthRequest
}

type sessionNextTurnRequest struct {
	SessionId string `json:"session_id" example:"string"`
	PlayerId  string `json:"player_id" example:"string"`
	AuthRequest
}

type authRegisterRequest struct {
	Email    string `json:"email" example:"string"`
	Password string `json:"password" example:"string"`
	Nickname string `json:"nickname" example:"string"`
	DefaultRequest
}

type authLoginRequest struct {
	Email    string `json:"email" example:"string"`
	Password string `json:"password" example:"string"`
	DefaultRequest
}

type authLogoutRequest struct {
	Email string `json:"email" example:"string"`
	Token string `json:"token" example:"string"`
	DefaultRequest
}

type PlayerResponse struct {
	Id        string `json:"id" example:"string"`
	Name      string `json:"name" example:"string"`
	Cards     []string `json:"players" example:"string"`
	State     string  `json:"state" example:"string"`
	SessionId string
}

type SessionResponse struct {
	Id            string    `json:"id" example:"string"`
	Players       []string `json:"players" example:"string"`
	Deck          []string `json:"deck" example:"string"`
	Table         []string `json:"table" example:"string"`
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
			Name:      player.Nickname,
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

type authRegisterResponse struct {
	DefaultResponse
}

type UserResponse struct {
	Id       string `json:"id" example:"string"`
	Nickname string `json:"nickname" example:"string"`
	Token    string `json:"token" example:"string"`
}

func NewUserResponse(user *core.User) *UserResponse {
	return &UserResponse{
		Id:       user.Id,
		Nickname: user.Nickname,
		Token:    user.Token,
	}
}

type authLoginResponse struct {
	User UserResponse `json:"user" example:UserResponse`
	DefaultResponse
}

type authLogoutResponse struct {
	DefaultResponse
}

type ErrResponse struct {
	Code int `json:"-"`

	Success bool   `json:"success" example:"false"`
	Message string `json:"message,omitempty" example:"Error occured"`
}

func (er ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	er.Success = false
	render.Status(r, er.Code)
	return nil
}

type AuthRequest struct {
	UserId string `json:"user_id" example:"string"`
	Token  string `json:"token" example:"string"`
	DefaultRequest
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
