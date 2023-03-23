package core

import (
	"errors"

	"github.com/MrBTTF/gophercises/deck"
)

////go:generate mockgen -source=ports.go  -destination=port_mocks.go -package=core

var (
	NoRoomForUserError = errors.New("User has no room")
)

type SessionRepository interface {
	Get(string) (Session, error)
	Store(*Session) error
	Delete(string) error
}

type PlayerRepository interface {
	Get(string) (Player, error)
	Store(*Player) error
}

type UserRepository interface {
	Get(string) (User, error)
	GetByEmail(string) (User, error)
	GetForRoom(string) ([]User, error)
	Store(*User) error
}

type RoomRepository interface {
	Get(string) (Room, error)
	GetByUserId(string) (string, error)
	List(bool) ([]Room, error)
	Store(*Room) error
	Delete(string) error
}

type SessionServicePort interface {
	GetSession(string) (Session, error)
	GetPlayer(string) (Player, error)
	Create(string, []deck.Card) (string, error)
	Pull(string, string) error
	Lay(string, string, Card) error
	NextTurn(string, string) error
	DeleteSession(string) error
}

type AuthServicePort interface {
	Login(email, password string) (User, error)
	Register(email, password, nickname string) error
	Logout(email, token string) error
	ValidateToken(user_id, token string) error
}

type RoomServicePort interface {
	Create(host_id string) (string, error)
	Get(room_id string) (Room, error)
	GetUsers(room_id string) ([]User, error)
	Join(room_id, user_id string) error
	List(open bool) ([]Room, error)
	Close(room_id string) error
	Delete(room_id string) error
}
