package accounts

import (
	"fmt"
	"github.com/google/uuid"
)

// удаление счета
func (db *Repo) DeleteAccount(accountID uuid.UUID, userID uuid.UUID) error {
	query := `
				DELETE FROM bank_app.accounts
				WHERE id = $1 AND user_id = $2
			`
	_, err := db.DB.Exec(query, accountID, userID)
	if err != nil {
		return fmt.Errorf("error in DeleteAccount query: %w", err)
	}

	return nil
}
