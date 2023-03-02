package core

import "github.com/MrBTTF/gophercises/deck"

////go:generate mockgen -source=ports.go  -destination=port_mocks.go -package=core

type SessionRepository interface {
	Get(string) (Session, error)
	Store(*Session) error
}

type PlayerRepository interface {
	Get(string) (Player, error)
	Store(*Player) error
}

type UserRepository interface {
	Get(string) (User, error)
	GetByEmail(string) (User, error)
	Store(*User) error
}

type SessionServicePort interface {
	GetSession(string) (Session, error)
	GetPlayer(string) (Player, error)
	Create([]string, []deck.Card) (string, error)
	Pull(string, string) error
	Lay(string, string, Card) error
	NextTurn(string, string) error
}

type AuthServicePort interface {
	Login(email, password string) (User, error)
	Register(email, password, nickname string) error
	Logout(email, token string) error
	ValidateToken(user_id, token string) error
}
