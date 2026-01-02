package models

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID          uuid.UUID
	UserFrom    uuid.UUID
	AccountFrom uuid.UUID
	UserTo      uuid.UUID
	AccountTo   uuid.UUID
	Amount      int
	Currency    Currency
	Timestamp   time.Time
}
