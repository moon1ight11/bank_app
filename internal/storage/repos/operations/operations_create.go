package operations

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// создание операции
func (db *Repo) CreateOperation(ownerID uuid.UUID, accountID uuid.UUID, amount float64, operationType string, currency string, tx *sql.Tx) (uuid.UUID, error) {
	query := `
				INSERT INTO bank_app.operations (user_id, account_id, operation_type, currency, amount)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id
			`
	var operationID uuid.UUID
	err := tx.QueryRow(query, ownerID, accountID, operationType, currency, amount).Scan(&operationID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error in CreateOperation query: %w", err)
	}

	return operationID, nil
}
