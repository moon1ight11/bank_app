package accounts

import (
	"bank_app/internal/storage/repos/models"
	"github.com/google/uuid"
)

type GetAccount struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Balance  int
	Currency models.Currency
}
