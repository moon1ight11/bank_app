package transactions

import (
	"bank_app/internal/storage"
)

type Repo struct {
	storage.DataBase
}

func NewTransactionsRepo(db *storage.DataBase) *Repo {
	return &Repo{DataBase: *db}
}
