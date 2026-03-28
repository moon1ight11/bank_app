package accounts

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// удаление счета
func (db *Repo) DeleteAccount(ctx context.Context, accountID uuid.UUID, userID uuid.UUID, tx *sql.Tx) error {
	query := `
				DELETE FROM bank_app.accounts
				WHERE id = $1 AND user_id = $2
			`
	_, err := tx.ExecContext(ctx, query, accountID, userID)
	if err != nil {
		return fmt.Errorf("error in DeleteAccount query: %w", err)
	}

	return nil
}
