package transactions

import (
	"bank_app/internal/storage/repos/models"
	"fmt"
	"github.com/google/uuid"
)

// получение всех транзакций пользователя
func (db *Repo) GetAllUsersTransactions(userID uuid.UUID) ([]models.Transaction, error) {
	query := `
				SELECT id, user_from, account_from, user_to, account_to, amount, currency, timestamp
				FROM bank_app.transactions
				WHERE user_from = $1 OR user_to = $1
				ORDER BY timestamp DESC
			`
	var transactions []models.Transaction
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("Error in GetAllUsersTransactions query: %w", err)
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
			return nil, fmt.Errorf("Error in GetAllUsersTransactions scan: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// получение транзакций конкретного счета
func (db *Repo) GetTransactionsByAccount(userID uuid.UUID, accountID uuid.UUID) ([]models.Transaction, error) {
	query := `
				SELECT id, user_from, account_from, user_to, account_to, amount, currency, timestamp
				FROM bank_app.transactions
				WHERE (user_from = $1 AND account_from = $2) OR (user_to = $1 AND account_to = $2)
			`
	var transactions []models.Transaction
	rows, err := db.DB.Query(query, userID, accountID)
	if err != nil {
		return nil, fmt.Errorf("Error in GetTransactionsByAccount query: %w", err)
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
			return nil, fmt.Errorf("Error in GetTransactionsByAccount scan: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// получение одной транзакции
func (db *Repo) GetTransactionByID(transactionID uuid.UUID, userID uuid.UUID) (models.Transaction, error) {
	query := `
				SELECT id, user_from, account_from, user_to, account_to, amount, currency, timestamp
				FROM bank_app.transactions
				WHERE (user_from = $1 OR user_to = $1) AND id = $2
			`
	var transaction models.Transaction

	err := db.DB.QueryRow(query, userID, transactionID).Scan(
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
