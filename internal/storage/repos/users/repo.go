package users

import (
	"bank_app/internal/storage"
)

type Repo struct {
	storage.DataBase
}

func NewUsersRepo(db *storage.DataBase) *Repo {
	return &Repo{DataBase: *db}
}