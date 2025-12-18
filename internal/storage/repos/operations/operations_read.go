package operations

import (
	"fmt"
	"github.com/google/uuid"
)

// получение всех операций пользователя
func (db *Repo) GetAllUsersOperations(ownerID uuid.UUID) ([]Operation, error) {
	query := `
				SELECT id, user_id, account_id, amount, currency, operation_type, timestamp
				FROM bank_app.operations
				WHERE user_id = $1 
			`
	var operations []Operation
	rows, err := db.DB.Query(query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("Error in GetAllUsersOperations query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var operation Operation
		err := rows.Scan(
			&operation.ID,
			&operation.OwnerID,
			&operation.AccountID,
			&operation.Amount,
			&operation.Currency,
			&operation.OperationType,
			&operation.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("Error in GetAllUsersOperations scan: %w", err)
		}
		operations = append(operations, operation)
	}
	return operations, nil
}

// получение операций конкретного счета
func (db *Repo) GetOperationsByAccount(ownerID uuid.UUID, accountID uuid.UUID) ([]Operation, error) {
	query := `
				SELECT id, user_id, account_id, amount, currency, operation_type, timestamp
				FROM bank_app.operations
				WHERE user_id = $1 AND account_id = $2
			`
	var operations []Operation
	rows, err := db.DB.Query(query, ownerID, accountID)
	if err != nil {
		return nil, fmt.Errorf("Error in GetOperationsByAccount query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var operation Operation
		err := rows.Scan(
			&operation.ID,
			&operation.OwnerID,
			&operation.AccountID,
			&operation.Amount,
			&operation.Currency,
			&operation.OperationType,
			&operation.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("Error in GetOperationsByAccount scan: %w", err)
		}
		operations = append(operations, operation)
	}
	return operations, nil
}

// получение одной операции
func (db *Repo) GetOperationByID(operationID uuid.UUID, userID uuid.UUID) (Operation, error) {
	query := `
				SELECT id, user_id, account_id, operation_type, amount, currency, timestamp
				FROM bank_app.operations
				WHERE user_id = $1 AND id = $2
			`
	var operation Operation

	err := db.DB.QueryRow(query, userID, operationID).Scan(
		&operation.ID,
		&operation.OwnerID,
		&operation.AccountID,
		&operation.OperationType,
		&operation.Amount,
		&operation.Currency,
		&operation.Timestamp,
	)
	if err != nil {
		return Operation{}, fmt.Errorf("error in GetOperationByID query: %w", err)
	}

	return operation, nil
}
