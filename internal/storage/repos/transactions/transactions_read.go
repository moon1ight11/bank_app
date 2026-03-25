package transactions

import (
	"bank_app/internal/storage/repos/models"
	"context"
	"fmt"
	"github.com/google/uuid"
)

// получение всех транзакций пользователя
func (db *Repo) GetAllUsersTransactions(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	query := `
				SELECT id, user_from, account_from, user_to, account_to, amount, currency, timestamp
				FROM bank_app.transactions
				WHERE user_from = $1 OR user_to = $1
				ORDER BY timestamp DESC
			`
	var transactions []models.Transaction

	rows, err := db.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error in GetAllUsersTransactions query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserFrom,
			&transaction.AccountFrom,
			&transaction.UserTo,
			&transaction.AccountTo,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("error in GetAllUsersTransactions scan: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// получение транзакций конкретного счета
func (db *Repo) GetTransactionsByAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) ([]models.Transaction, error) {
	query := `
				SELECT id, user_from, account_from, user_to, account_to, amount, currency, timestamp
				FROM bank_app.transactions
				WHERE (user_from = $1 AND account_from = $2) OR (user_to = $1 AND account_to = $2)
			`
	var transactions []models.Transaction
	rows, err := db.DB.QueryContext(ctx, query, userID, accountID)
	if err != nil {
		return nil, fmt.Errorf("error in GetTransactionsByAccount query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserFrom,
			&transaction.AccountFrom,
			&transaction.UserTo,
			&transaction.AccountTo,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("error in GetTransactionsByAccount scan: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// получение одной транзакции
func (db *Repo) GetTransactionByID(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) (models.Transaction, error) {
	query := `
				SELECT id, user_from, account_from, user_to, account_to, amount, currency, timestamp
				FROM bank_app.transactions
				WHERE (user_from = $1 OR user_to = $1) AND id = $2
			`
	var transaction models.Transaction

	err := db.DB.QueryRowContext(ctx, query, userID, transactionID).Scan(
		&transaction.ID,
		&transaction.UserFrom,
		&transaction.AccountFrom,
		&transaction.UserTo,
		&transaction.AccountTo,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.Timestamp,
	)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("error in GetTransactionByID query: %w", err)
	}

	return transaction, nil
}
