package accounts

import (
	"bank_app/internal/storage/repos/transactions"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// получение счета по id
func (db *Repo) GetAccountById(accountID uuid.UUID, userID uuid.UUID) (Account, error) {
	query := `
				SELECT id, user_id, balance, currency
				FROM bank_app.accounts
				WHERE user_id = $1 AND id = $2
			`
	var account Account

	err := db.DB.QueryRow(query, userID, accountID).Scan(&account.ID, &account.UserID, &account.Balance, &account.Currency)
	if err != nil {
		return Account{}, fmt.Errorf("error in GetAccountById query: %w", err)
	}

	return account, nil
}

// получение счетов по id пользователя
func (db *Repo) GetAccountsByUserId(userID uuid.UUID) ([]Account, error) {
	query := `
				SELECT id, user_id, balance, currency
				FROM bank_app.accounts
				WHERE user_id = $1
			`
	var accounts []Account
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error in GetAccountsByUserId query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var account Account
		err := rows.Scan(
			&account.ID,
			&account.UserID,
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

// получение админского счета с указанной валютой
func (db *Repo) GetAdminAccountByCurrency(currency transactions.Currency) (Account, error) {
	query := `
				SELECT id, user_id, balance, currency
				FROM bank_app.accounts
				WHERE user_id = $1 AND currency = $2
			`
	var account Account
	UserID, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if err != nil {
		return Account{}, fmt.Errorf("error in Parse UUID: %w", err)
	}

	err = db.DB.QueryRow(query, UserID, currency).Scan(&account.ID, &account.UserID, &account.Balance, &account.Currency)
	if err != nil {
		return Account{}, fmt.Errorf("error in GetAccountById query: %w", err)
	}

	return account, nil
}

// проверка на количество средств перед транзакцией и заморозка счета
func (db *Repo) BalanceCheck(accountID uuid.UUID, amount int, tx *sql.Tx) error {
	query := `
				SELECT balance
				FROM bank_app.accounts
				WHERE id = $1
				FOR UPDATE
			`
	var balance int
	err := tx.QueryRow(query, accountID).Scan(&balance)
	if err != nil {
		return fmt.Errorf("error in BalanceCheck query: %w", err)
	}

	if balance < amount {
		return fmt.Errorf("not enough founds")
	}

	return nil
}
