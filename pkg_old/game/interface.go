package game


type DB interface {
	Close() error
	CreateSession(string, string) (string, error)
	GetSession(string) (*Session, error)
	ListSessions() ([]*Session, error)
	SaveSession(*Session) error
}