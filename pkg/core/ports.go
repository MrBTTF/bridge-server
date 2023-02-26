package core

////go:generate mockgen -source=ports.go  -destination=port_mocks.go -package=core

type SessionRepository interface {
	Get(string) (Session, error)
	Store(*Session) error
}

type PlayerRepository interface {
	Get(string) (Player, error)
	Store(*Player) error
}

type SessionServicePort interface {
	Create([]int) error
	Pull(int, int) error
	Lay(int, int, Card) error
	NextTurn(int) error
}
