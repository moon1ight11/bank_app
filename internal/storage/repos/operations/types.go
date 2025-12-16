package operations

import (
	"bank_app/internal/storage"
	"time"

	"github.com/google/uuid"
)

type Repo struct {
	storage.DataBase
}

func NewOperationsRepo(db *storage.DataBase) *Repo {
	return &Repo{DataBase: *db}
}

type Operation struct {
	ID            uuid.UUID `json:"operation_id"`
	OwnerID       uuid.UUID `json:"owner_id"`
	AccountID     uuid.UUID `json:"account_id"`
	OperationType string    `json:"type"`
	Amount        float64       `json:"amount"`
	Timestamp     time.Time `json:"timestamp"`
}
