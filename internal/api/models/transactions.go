package models

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID          uuid.UUID `json:"transaction_id"`
	UserFrom    uuid.UUID `json:"user_from"`
	AccountFrom uuid.UUID `json:"account_from"`
	UserTo      uuid.UUID `json:"user_to"`
	AccountTo   uuid.UUID `json:"account_to"`
	Amount      int       `json:"amount"`
	Currency    Currency  `json:"currency"`
	Timestamp   time.Time `json:"timestamp"`
}
