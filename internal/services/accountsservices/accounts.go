package accountsservices

import (
	"bank_app/internal/storage/cache"
	"bank_app/internal/storage/repos/accounts"
	"bank_app/internal/storage/repos/transactions"
)

type AccountsService struct {
	accountsRepo     *accounts.Repo
	transactionsRepo *transactions.Repo
	cacheService     cache.CacheInterface
}

func NewAccountsService(
	accountsRepo *accounts.Repo,
	transactionsRepo *transactions.Repo,
	cacheService cache.CacheInterface,
) *AccountsService {
	return &AccountsService{
		accountsRepo:     accountsRepo,
		transactionsRepo: transactionsRepo,
		cacheService:     cacheService,
	}
}
