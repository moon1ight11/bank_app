package accounts

import (
	"bank_app/internal/storage"
	"bank_app/internal/storage/repos/transactions"

	"github.com/google/uuid"
)

type Repo struct {
	storage.DataBase
}

func NewAccountsRepo(db *storage.DataBase) *Repo {
	return &Repo{DataBase: *db}
}

type Account struct {
	ID       uuid.UUID             `json:"account_id"`
	UserID   uuid.UUID             `json:"owner_id"`
	Balance  int                   `json:"balance"`
	Currency transactions.Currency `json:"currency"`
}
