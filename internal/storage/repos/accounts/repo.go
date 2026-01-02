package accounts

import "bank_app/internal/storage"

type Repo struct {
	storage.DataBase
}

func NewAccountsRepo(db *storage.DataBase) *Repo {
	return &Repo{DataBase: *db}
}