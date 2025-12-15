package operations

import (
	"fmt"
	"github.com/google/uuid"
)

// создание операции
func (db *Repo) CreateOperation(ownerID uuid.UUID, accountID uuid.UUID, amount int, operationType string) (uuid.UUID, error) {
	query := `
				INSERT INTO bank_app.operations (user_id, account_id, operation_type, amount)
				VALUES ($1, $2, $3, $4)
				RETURNING id
			`
	var operationID uuid.UUID
	err := db.DB.QueryRow(query, ownerID, accountID, amount, operationType).Scan(&operationID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error in CreateOperation query: %w", err)
	}

	return operationID, nil
}
