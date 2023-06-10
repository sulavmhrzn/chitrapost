package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Chitra struct {
	ID     int    `json:"id"`
	URL    string `json:"chitra"`
	UserID int    `json:"userid"`
}

type ChitraModel struct {
	db *pgx.Conn
}

func (m ChitraModel) Insert(c *Chitra) (*Chitra, error) {
	query := `INSERT INTO chitras (chitra_url, user_id) VALUES ($1, $2) RETURNING id`
	args := []any{c.URL, c.UserID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.db.QueryRow(ctx, query, args...).Scan(&c.ID)
	if err != nil {
		return nil, err
	}
	return c, nil
}
