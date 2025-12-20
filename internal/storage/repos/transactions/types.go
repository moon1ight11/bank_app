package transactions

import (
	"bank_app/internal/storage"
	"github.com/google/uuid"
	"time"
)

type Repo struct {
	storage.DataBase
}

func NewTransactionsRepo(db *storage.DataBase) *Repo {
	return &Repo{DataBase: *db}
}

type Transaction struct {
	ID              uuid.UUID       `json:"transaction_id"`
	UserFrom        uuid.UUID       `json:"user_from"`
	AccountFrom     uuid.UUID       `json:"account_from"`
	UserTo          uuid.UUID       `json:"user_to"`
	AccountTo       uuid.UUID       `json:"account_to"`
	TransactionType TransactionType `json:"transaction_type"`
	Amount          int             `json:"amount"`
	Currency        Currency        `json:"currency"`
	Timestamp       time.Time       `json:"timestamp"`
}

type TransactionType string

const (
	TransactionIncome   TransactionType = "Income"
	TransactionOutcome  TransactionType = "Outcome"
	TransactionTransfer TransactionType = "Transfer"
)

type Currency string

const (
	CurrencyRUB Currency = "RUB"
	CurrencyEUR Currency = "EUR"
	CurrencyUSD Currency = "USD"
)
