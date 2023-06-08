package data

import (
	"context"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"-" validate:"required,min=5,max=20"`
}

type UserModel struct {
	db *pgx.Conn
}

func CreateHashPassword(plain string) (string, error) {
	hash, err := argon2id.CreateHash(plain, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func ValidateUser(u *User) error {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return err
	}
	return nil
}

func (m UserModel) Insert(u *User) (*User, error) {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	args := []any{u.Email, u.Password}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.db.QueryRow(ctx, query, args...).Scan(&u.ID)
	if err != nil {
		code := err.(*pgconn.PgError).Code
		if code == "23505" {
			return nil, errors.New("email already exists")
		}
		return nil, err
	}
	return u, nil
}
