package services

import "bank_app/internal/storage/repos/accounts"

type AccountsService struct {
	accountsRepo *accounts.Repo
}

func NewAccountsService(accountsRepo *accounts.Repo) *AccountsService {
	return &AccountsService{accountsRepo: accountsRepo}
}
