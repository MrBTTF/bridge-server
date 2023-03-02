package auth

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mrbttf/bridge-server/pkg/core"
)

const tokenLength = 16

var (
	LoginInvalidError = errors.New("Invalid email or password")
	TokenInvalidError = errors.New("Invalid token")
)

type AuthService struct {
	user core.UserRepository
}

func New(user core.UserRepository) *AuthService {
	return &AuthService{
		user: user,
	}
}

func (as *AuthService) Login(email string, password string) (core.User, error) {
	user, err := as.fetchUser(email, password)
	if err != nil {
		return core.User{}, fmt.Errorf("Unable to login: %w", err)
	}
	token := generateSecureToken(tokenLength)
	user.Token = getHash(token)
	err = as.user.Store(&user)
	if err != nil {
		return core.User{}, fmt.Errorf("Unable to login: %w", err)
	}
	user.Token = token
	return user, nil
}

func (as *AuthService) Register(email string, password string, nickname string) error {
	passwordHash := getHash(password)
	user := core.User{
		Id:       uuid.New().String(),
		Email:    email,
		Password: passwordHash,
		Nickname: nickname,
	}
	err := as.user.Store(&user)
	if err != nil {
		return fmt.Errorf("Unable to register: %w", err)
	}
	return nil
}

func (as *AuthService) Logout(email, token string) error {
	user, err := as.user.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("Unable to logout for email: %s: %w", err)
	}
	if user.Token != getHash(token) {
		return fmt.Errorf("Unable to logout for email: %s: %w", email, TokenInvalidError)
	}
	user.Token = ""
	err = as.user.Store(&user)
	if err != nil {
		return fmt.Errorf("Unable to logout for email: %s: %w", err)
	}
	return nil
}

func (as *AuthService) ValidateToken(user_id, token string) error {
	user, err := as.user.Get(user_id)
	if err != nil {
		return fmt.Errorf("Unable to validate token for user_id: %s: %w", user_id, err)
	}
	if user.Token != getHash(token) {
		return fmt.Errorf("Unable to validate token for user_id: %s: %w", user_id, TokenInvalidError)
	}
	return nil
}

func (as *AuthService) fetchUser(email, password string) (core.User, error) {
	user, err := as.user.GetByEmail(email)
	if err != nil {
		return core.User{}, fmt.Errorf("Unable to fetch user for email: %s: %w", email, err)
	}
	passwordHash := getHash(password)
	if user.Password != passwordHash {
		return core.User{}, fmt.Errorf("Unable to fetch user for email: %s: %w", email, LoginInvalidError)
	}
	return user, nil
}

func getHash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
