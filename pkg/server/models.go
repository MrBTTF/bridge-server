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
	RoomId string `json:"room_id" example:"string"`
	AuthRequest
}

type sessionGetByUserRequest struct {
	UserId string `json:"user_id" example:"string"`
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

type sessionCloseRequest struct {
	SessionId string `json:"session_id" example:"string"`
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

type roomGetRequest struct {
	RoomId string `json:"room_id" example:"string"`
	AuthRequest
}

type roomCreateRequest struct {
	HostId string `json:"host_id" example:"string"`
	AuthRequest
}

type roomJoinRequest struct {
	RoomId string `json:"room_id" example:"string"`
	UserId string `json:"user_id" example:"string"`
	AuthRequest
}

type roomListRequest struct {
	Open bool `json:"open" example:true`
	AuthRequest
}

type roomDeleteRequest struct {
	RoomId string `json:"room_id" example:"string"`
	AuthRequest
}

type PlayerResponse struct {
	Id        string   `json:"id" example:"string"`
	Name      string   `json:"name" example:"string"`
	Cards     []string `json:"players" example:"string"`
	State     string   `json:"state" example:"string"`
	SessionId string
}

func NewPlayerResponse(player *core.Player) *PlayerResponse {
	return &PlayerResponse{
		Id:        player.Id,
		Name:      player.Nickname,
		Cards:     repositories.DeckToString(player.Cards),
		State:     player.State.String(),
		SessionId: player.SessionId,
	}
}

type SessionResponse struct {
	Id            string   `json:"id" example:"string"`
	Players       []string `json:"players" example:"string"`
	Deck          []string `json:"deck" example:"string"`
	Table         []string `json:"table" example:"string"`
	CurrentPlayer PlayerResponse
}

func NewSessionResponse(session *core.Session, player *core.Player) *SessionResponse {
	return &SessionResponse{
		Id:            session.Id,
		Deck:          repositories.DeckToString(session.Deck),
		Table:         repositories.DeckToString(session.Table),
		Players:       session.Players,
		CurrentPlayer: *NewPlayerResponse(player),
	}
}

type sessionGetResponse struct {
	Session SessionResponse `json:"session" example:SessionResponse`
	DefaultResponse
}

type sessionGetByUserResponse struct {
	SessionId string `json:"session_id" example:"string"`
	DefaultResponse
}

type sessionCreateResponse struct {
	SessionID string `json:"session_id" example:"string"`
	DefaultResponse
}

type sessionNoSessionErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Msg     string `json:"message" example:"User has no session"`
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

type UserResponseSecure struct {
	Id       string `json:"id" example:"string"`
	Nickname string `json:"nickname" example:"string"`
}

func NewUserResponseSecure(user *core.User) *UserResponseSecure {
	return &UserResponseSecure{
		Id:       user.Id,
		Nickname: user.Nickname,
	}
}

type authLoginResponse struct {
	User UserResponse `json:"user" example:UserResponse`
	DefaultResponse
}

type authLogoutResponse struct {
	DefaultResponse
}

type RoomResponse struct {
	Id    string               `json:"id" example:"string"`
	Host  UserResponseSecure   `json:"host" example:UserResponseSecure`
	Users []UserResponseSecure `json:"users" example:UserResponseSecure`
	Open  bool                 `json:"open" example:"true"`
}

func NewRoomResponse(room *core.Room, users []core.User) *RoomResponse {
	users_response := make([]UserResponseSecure, 0, len(users))
	var host *core.User
	for _, user := range users {
		users_response = append(users_response, *NewUserResponseSecure(&user))
		if user.Id == room.Host {
			u := user
			host = &u
		}
	}

	return &RoomResponse{
		Id:    room.Id,
		Host:  *NewUserResponseSecure(host),
		Users: users_response,
		Open:  room.Open,
	}
}

type roomGetResponse struct {
	Room RoomResponse `json:"room" example:RoomResponse`
	DefaultResponse
}

type roomCreateResponse struct {
	RoomId string `json:"room_id" example:"string"`
	DefaultResponse
}

type RoomUser struct {
	room  core.Room
	users []core.User
}

type roomListResponse struct {
	Rooms []RoomResponse `json:"rooms" example:RoomResponse`
	DefaultResponse
}

func NewRoomListResponse(roomUsers []RoomUser) *roomListResponse {
	rooms_response := make([]RoomResponse, 0, len(roomUsers))
	for _, ru := range roomUsers {
		rooms_response = append(rooms_response, *NewRoomResponse(&ru.room, ru.users))
	}
	return &roomListResponse{
		Rooms: rooms_response,
	}
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
