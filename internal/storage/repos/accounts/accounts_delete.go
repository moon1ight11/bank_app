package accounts

import (
	"fmt"
	"github.com/google/uuid"
)

// удаление счета
func (db *Repo) DeleteAccount(accountID uuid.UUID) error {
	query := `
				DELETE FROM bank_app.accounts
				WHERE id = $1
			`
	_, err := db.DB.Exec(query, accountID)
	if err != nil {
		return fmt.Errorf("error in DeleteAccount query: %w", err)
	}

	return nil
}
