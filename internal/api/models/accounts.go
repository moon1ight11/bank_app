package models

import "github.com/google/uuid"

type Currency string

const (
	CurrencyRUB Currency = "RUB"
	CurrencyEUR Currency = "EUR"
	CurrencyUSD Currency = "USD"
)

// модель данных для создания счета
type AccountCreate struct {
	Currency Currency  `json:"currency"`
	UserID   uuid.UUID `json:"user_id"`
}

// модель данных для получения счета
type AccountsGet struct {
	ID       uuid.UUID `json:"account_id"`
	UserID   uuid.UUID `json:"user_id"`
	Balance  int       `json:"balance"`
	Currency Currency  `json:"currency"`
}
