package accounts

import (
	"bank_app/internal/storage"
	"github.com/google/uuid"
)

type Repo struct {
	storage.DataBase
}

func NewAccountsRepo(db *storage.DataBase) *Repo {
	return &Repo{DataBase: *db}
}

type Account struct {
	ID       uuid.UUID `json:"account_id"`
	OwnerID  uuid.UUID `json:"owner_id"`
	Balance  float64   `json:"balance"`
	Currency string    `json:"currency"`
}
