package transaction

import (
	"bank_app/internal/storage/repos/accounts"
	"bank_app/internal/storage/repos/transactions"
)

type TransactionsService struct {
	transactionsRepo *transactions.Repo
	accountsRepo     *accounts.Repo
}

func NewTransactionsService(transactionsRepo *transactions.Repo, accountsRepo *accounts.Repo) *TransactionsService {
	return &TransactionsService{
		transactionsRepo: transactionsRepo,
		accountsRepo:     accountsRepo,
	}
}