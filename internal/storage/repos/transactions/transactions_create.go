package transactions

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// создание транзакции
func (db *Repo) CreateTransaction(
	userFrom uuid.UUID,
	accountFrom uuid.UUID,
	userTo uuid.UUID,
	accountTo uuid.UUID,
	amount int,
	currency string,
	tx *sql.Tx,
) (uuid.UUID, error) {
	query := `
				INSERT INTO bank_app.transactions (user_from, account_from, user_to, account_to, amount, currency)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id
			`
	var transactionID uuid.UUID
	err := tx.QueryRow(
		query, 
		userFrom, 
		accountFrom, 
		userTo,
		accountTo,
		amount,
		currency,
		).Scan(&transactionID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error in CreateTransaction query: %w", err)
	}

	return transactionID, nil
}
