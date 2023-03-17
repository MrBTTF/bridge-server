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
		return core.User{}, fmt.Errorf("Unable to get user by email for id %s: %w", user_id, err)
	}
	return user, nil
}

const SelectUsersForRoom = `
SELECT user_id, email, password, nickname, token
FROM users
JOIN rooms ON user_id = any(rooms.user_ids)
WHERE rooms.room_id = $1
`

func (ur *UserRepository) GetForRoom(room_id string) ([]core.User, error) {
	rows, err := ur.db.Query(SelectUsersForRoom, room_id)
	if err != nil {
		return nil, fmt.Errorf("Unable to get users for room id %s: %w", room_id, err)
	}
	defer rows.Close()

	var users []core.User
	for rows.Next() {
		var user core.User
		if err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.Password,
			&user.Nickname,
			&user.Token,
		); err != nil {
			return nil, fmt.Errorf("Unable to get users for room id %s: %w", room_id, err)
		}
		users = append(users, user)
	}
	return users, nil
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
