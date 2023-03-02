package repositories

import (
	"database/sql"
	"fmt"

	"github.com/mrbttf/bridge-server/pkg/core"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

const SelectUser = `
SELECT user_id, email, password, nickname, token
FROM users
WHERE user_id = $1
`

func (ur *UserRepository) Get(user_id string) (core.User, error) {
	var user core.User
	err := ur.db.QueryRow(SelectUser, user_id).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Nickname,
		&user.Token,
	)
	if err != nil {
		return core.User{}, fmt.Errorf("Unable to get user for id %s: %w", user_id, err)
	}

	return user, nil
}

const SelectUserByEmail = `
SELECT user_id, email, password, nickname, token
FROM users
WHERE email = $1
`

func (ur *UserRepository) GetByEmail(user_id string) (core.User, error) {
	var user core.User
	err := ur.db.QueryRow(SelectUserByEmail, user_id).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Nickname,
		&user.Token,
	)
	if err != nil {
		return core.User{}, fmt.Errorf("Unable to get user ny email for id %s: %w", user_id, err)
	}
	return user, nil
}

const UpsertUser = `
INSERT INTO users (user_id, email, password, nickname, token)
VALUES($1, $2, $3, $4, $5) 
ON CONFLICT (user_id) 
WHERE user_id = $1 
DO UPDATE
SET 
	email = EXCLUDED.email, 
	password = EXCLUDED.password, 
	nickname = EXCLUDED.nickname, 
	token = EXCLUDED.token
`

func (ur *UserRepository) Store(user *core.User) error {
	println(user.Id,
		user.Email,
		user.Password,
		user.Nickname,
		user.Token,
	)
	_, err := ur.db.Exec(UpsertUser,
		user.Id,
		user.Email,
		user.Password,
		user.Nickname,
		user.Token,
	)
	if err != nil {
		return err
	}

	return nil
}
