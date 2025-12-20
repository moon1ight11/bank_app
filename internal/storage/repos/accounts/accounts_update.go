package accounts

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// пополнение на счет
func (db *Repo) BalanceIncoming(accountID uuid.UUID, amount float64, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.accounts
				SET balance = balance + $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, accountID, amount)
	if err != nil {
		return fmt.Errorf("error in BalanceIncoming query: %w", err)
	}

	return nil
}

// списание со счета
func (db *Repo) BalanceOutlay(accountID uuid.UUID, amount float64, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.accounts
				SET balance = balance - $2, updated_at = NOW()
				WHERE id = $1
			`
	_, err := tx.Exec(query, accountID, amount)
	if err != nil {
		return fmt.Errorf("error in BalanceOutlay query: %w", err)
	}

	return nil
}
