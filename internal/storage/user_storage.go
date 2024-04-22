package storage

import (
	"context"

	"github.com/1minepowminx/distributed_calculator/internal/utils/orchestrator/jwts"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64
	Name     string
	Password string
}

// Регистрация нового пользователя
func (s *Storage) RegisterUser(ctx context.Context, uname, pswrd string) error {
	q := `
	INSERT INTO users (name, password) values ($1, $2)
	`

	hashedPswrd, err := bcrypt.GenerateFromPassword([]byte(pswrd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, q, uname, hashedPswrd)
	if err != nil {
		return err
	}

	return nil
}

// Аутентефикация пользователя
func (s *Storage) LoginUser(ctx context.Context, uname, pswrd string) (string, error) {
	q := `
	SELECT id, password FROM users WHERE name=$1
	`
	var user User

	err := s.db.QueryRowContext(ctx, q, uname).Scan(&user.ID, &user.Password)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pswrd))
	if err != nil {
		return "", err
	}

	token, err := jwts.GenerateJWTToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
