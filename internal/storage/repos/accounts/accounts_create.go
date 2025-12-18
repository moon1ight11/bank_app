package accounts

import (
	"fmt"
	"github.com/google/uuid"
)

// создание счета
func (db *Repo) CreateAccount(ownerID uuid.UUID, currency string) (uuid.UUID, error) {
	query := `
				INSERT INTO bank_app.accounts (user_id, currency)
				VALUES ($1, $2)
				RETURNING id
			`
	var accountID uuid.UUID
	err := db.DB.QueryRow(query, ownerID, currency).Scan(&accountID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in CreateAccount query: %w", err)
	}

	return accountID, nil
}
