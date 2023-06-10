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

var (
	ErrDuplicateEmail = errors.New("email already exists")
	ErrNoRows         = errors.New("no rows exists")
)

func CreateHashPassword(plain string) (string, error) {
	hash, err := argon2id.CreateHash(plain, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func ComparePasswordAndHash(plain, hashed string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(plain, hashed)
	if err != nil {
		return false, err
	}
	return match, nil
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
			return nil, ErrDuplicateEmail
		}
		return nil, err
	}
	return u, nil
}

func (m UserModel) GetUser(user *User) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE email = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRow(ctx, query, user.Email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}
	return user, nil
}

func (m UserModel) GetUserFromID(id int) (*User, error) {
	query := `SELECT id, email FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &User{}
	err := m.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}

	return user, nil
}
