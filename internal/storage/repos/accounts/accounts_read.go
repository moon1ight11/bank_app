package accounts

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// получение счета по id
func (db *Repo) GetAccountById(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) (GetAccount, error) {
	query := `
				SELECT id, user_id, balance, currency
				FROM bank_app.accounts
				WHERE user_id = $1 AND id = $2
			`
	var account GetAccount

	err := db.DB.QueryRowContext(ctx, query, userID, accountID).Scan(&account.ID, &account.UserID, &account.Balance, &account.Currency)
	if err != nil {
		return GetAccount{}, fmt.Errorf("error in GetAccountById query: %w", err)
	}

	return account, nil
}

// получение счета по id и заморозка
func (db *Repo) GetAccountByIdTx(ctx context.Context, accountID uuid.UUID, userID uuid.UUID, tx *sql.Tx) (GetAccount, error) {
	query := `
				SELECT id, user_id, balance, currency
				FROM bank_app.accounts
				WHERE user_id = $1 AND id = $2
				FOR UPDATE
			`
	var account GetAccount

	err := tx.QueryRowContext(ctx, query, userID, accountID).Scan(&account.ID, &account.UserID, &account.Balance, &account.Currency)
	if err != nil {
		return GetAccount{}, fmt.Errorf("error in GetAccountByIdTx query: %w", err)
	}

	return account, nil
}

// получение счетов по id пользователя
func (db *Repo) GetAccountsByUserId(ctx context.Context, userID uuid.UUID) ([]GetAccount, error) {
	query := `
				SELECT id, user_id, balance, currency
				FROM bank_app.accounts
				WHERE user_id = $1
			`
	var accounts []GetAccount
	rows, err := db.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error in GetAccountsByUserId query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var account GetAccount
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
func (db *Repo) GetAdminAccountByCurrency(ctx context.Context, currency string) (GetAccount, error) {
	query := `
				SELECT id, user_id, balance, currency
				FROM bank_app.accounts
				WHERE user_id = $1 AND currency = $2
			`
	var account GetAccount

	UserID, err := uuid.Parse("00000000-0000-0000-0000-000000000001")
	if err != nil {
		return GetAccount{}, fmt.Errorf("error in Parse UUID in GetAccountById: %w", err)
	}

	err = db.DB.QueryRowContext(ctx, query, UserID, currency).Scan(&account.ID, &account.UserID, &account.Balance, &account.Currency)
	if err != nil {
		return GetAccount{}, fmt.Errorf("error in GetAccountById query: %w", err)
	}

	return account, nil
}

// проверка на количество средств перед транзакцией и заморозка счета
func (db *Repo) BalanceCheck(ctx context.Context, accountID uuid.UUID, amount int, tx *sql.Tx) error {
	query := `
				SELECT balance
				FROM bank_app.accounts
				WHERE id = $1
				FOR UPDATE
			`
	var balance int

	err := tx.QueryRowContext(ctx, query, accountID).Scan(&balance)
	if err != nil {
		return fmt.Errorf("error in BalanceCheck query: %w", err)
	}

	if balance < amount {
		return fmt.Errorf("error in BalanceCheck: not enough founds")
	}

	return nil
}
