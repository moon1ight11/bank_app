package accounts

import (
	"bank_app/internal/storage/repos/transactions"
	"fmt"

	"github.com/google/uuid"
)

// создание счета
func (db *Repo) CreateAccount(userID uuid.UUID, currency transactions.Currency) (uuid.UUID, error) {
	query := `
				INSERT INTO bank_app.accounts (user_id, currency)
				VALUES ($1, $2)
				RETURNING id
			`
	var accountID uuid.UUID
	err := db.DB.QueryRow(query, userID, currency).Scan(&accountID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in CreateAccount query: %w", err)
	}

	return accountID, nil
}
