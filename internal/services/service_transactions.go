package services

import (
	"bank_app/internal/storage/repos/transactions"
	"github.com/google/uuid"
)

type TransactionsService struct {
	transactionsRepo *transactions.Repo
}

func NewTransactionsService(transactionsRepo *transactions.Repo) *TransactionsService {
	return &TransactionsService{transactionsRepo: transactionsRepo}
}

// получение всех транзакций пользователя
func (t *TransactionsService) AllTransactionsGet(userID uuid.UUID) ([]transactions.Transaction, error) {
	transactions, err := t.transactionsRepo.GetAllUsersTransactions(userID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// получение транзакций по счету
func (t *TransactionsService) AccountTransactionsGet(userID uuid.UUID, accountID uuid.UUID) ([]transactions.Transaction, error) {
	transactions, err := t.transactionsRepo.GetTransactionsByAccount(userID, accountID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// получение одной транзакции
func (t *TransactionsService) TransactionByIdGet(userID uuid.UUID, transactionID uuid.UUID) (transactions.Transaction, error) {
	transacion, err := t.transactionsRepo.GetTransactionByID(transactionID, userID)
	if err != nil {
		return transactions.Transaction{}, err
	}

	return transacion, nil
}