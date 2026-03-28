package accounts

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// пополнение на счет
func (db *Repo) BalanceIncoming(ctx context.Context, accountID uuid.UUID, amount int, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.accounts
				SET balance = balance + $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.ExecContext(ctx, query, accountID, amount)
	if err != nil {
		return fmt.Errorf("error in BalanceIncoming query: %w", err)
	}

	return nil
}

// списание со счета
func (db *Repo) BalanceOutcoming(ctx context.Context, accountID uuid.UUID, amount int, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.accounts
				SET balance = balance - $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.ExecContext(ctx, query, accountID, amount)
	if err != nil {
		return fmt.Errorf("error in BalanceOutlay query: %w", err)
	}

	return nil
}
