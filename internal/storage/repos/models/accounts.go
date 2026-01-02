package models

import (
	"github.com/google/uuid"
	"time"
)

type Account struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Balance   int
	Currency  Currency
	CreatedAt time.Time
	UpdatedAt *time.Time
}
