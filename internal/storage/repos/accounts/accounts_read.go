package accounts

import (
	"fmt"
	"github.com/google/uuid"
)

// получение счета по id
func (db *Repo) GetAccountById(accountID uuid.UUID, ownerID uuid.UUID) (Account, error) {
	query := `
				SELECT id, user_id, balance, currency
				FROM bank_app.accounts
				WHERE owner_id = $1 AND id = $2
			`
	var account Account

	err := db.DB.QueryRow(query, ownerID, accountID).Scan(&account.ID, &account.OwnerID, &account.Balance, &account.Currency)
	if err != nil {
		return Account{}, fmt.Errorf("error in GetAccountById query: %w", err)
	}

	return account, nil
}

// получение счетов по id пользователя
func (db *Repo) GetAccountsByUserId(ownerID uuid.UUID) ([]Account, error) {
	query := `
				SELECT id, user_id, balance, currency
				FROM bank_app.accounts
				WHERE owner_id = $1
			`
	var accounts []Account
	rows, err := db.DB.Query(query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("error in GetAccountsByUserId query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var account Account
		err := rows.Scan(
			&account.ID,
			&account.OwnerID,
			&account.Balance,
			&account.Currency,
		)
		if err != nil {
			return nil, fmt.Errorf("error in GetAccountsByUserId scan: %w", err)
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
