package accounts

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

// создание счета
func (db *Repo) CreateAccount(ctx context.Context, userID uuid.UUID, currency string) (uuid.UUID, error) {
	query := `
				INSERT INTO bank_app.accounts (user_id, currency)
				VALUES ($1, $2)
				RETURNING id
			`
	var accountID uuid.UUID
	err := db.DB.QueryRowContext(ctx, query, userID, currency).Scan(&accountID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in CreateAccount query: %w", err)
	}

	return accountID, nil
}
