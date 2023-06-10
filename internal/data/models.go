package data

import "github.com/jackc/pgx/v5"

type Models struct {
	UserModel   UserModel
	ChitraModel ChitraModel
}

func NewModels(db *pgx.Conn) *Models {
	return &Models{
		UserModel:   UserModel{db},
		ChitraModel: ChitraModel{db},
	}
}
