package data

import "github.com/jackc/pgx/v5"

type Models struct {
	UserModel UserModel
}

func NewModels(db *pgx.Conn) *Models {
	return &Models{
		UserModel: UserModel{db},
	}
}
