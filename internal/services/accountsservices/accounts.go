package accountsservices

import (
	"bank_app/internal/storage/repos/accounts"
	"bank_app/internal/storage/repos/transactions"
)

type AccountsService struct {
	accountsRepo     *accounts.Repo
	transactionsRepo *transactions.Repo
}

func NewAccountsService(accountsRepo *accounts.Repo, transactionsRepo *transactions.Repo) *AccountsService {
	return &AccountsService{
		accountsRepo:   accountsRepo,
		transactionsRepo: transactionsRepo,
	}
}
