package accountsservices

import (
	"bank_app/internal/api/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// создание счета
func (a *AccountsService) AccountAdd(ctx context.Context, newAccount models.AccountCreate) (uuid.UUID, error) {
	accountID, err := a.accountsRepo.CreateAccount(ctx, newAccount.UserID, string(newAccount.Currency))
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in AccountAdd: %w", err)
	}

	cacheKey := fmt.Sprintf("account_user:%s", newAccount.UserID.String())
	if a.cacheService != nil {
		if err := a.cacheService.Set(ctx, cacheKey, accountID, 10*time.Minute); err != nil {
			return uuid.Nil, fmt.Errorf("error in accountAdd: %w; cachekey %s not set", err, cacheKey)
		}
	}

	return accountID, nil
}

// получение всех счетов пользователя
func (a *AccountsService) AllAccountsGet(ctx context.Context, userID uuid.UUID) ([]models.AccountsGet, error) {
	accountsRepo, err := a.accountsRepo.GetAccountsByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error in AllAccountsGet: %w", err)
	}

	var accounts []models.AccountsGet
	for i := range accountsRepo {
		var account models.AccountsGet

		account.ID = accountsRepo[i].ID
		account.UserID = accountsRepo[i].UserID
		account.Balance = accountsRepo[i].Balance
		account.Currency = models.Currency(string(accountsRepo[i].Currency))

		accounts = append(accounts, account)
	}

	return accounts, nil
}

// получение одного счета пользователя
func (a *AccountsService) AccountGet(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (models.AccountsGet, error) {
	cacheKey := fmt.Sprintf("user_account:%s", userID.String())

	if a.cacheService != nil {
		var cachedAccount models.AccountsGet
		err := a.cacheService.Get(ctx, cacheKey, &cachedAccount)
		if err == nil {
			return cachedAccount, nil
		}
	}

	accountRepo, err := a.accountsRepo.GetAccountById(ctx, accountID, userID)
	if err != nil {
		return models.AccountsGet{}, fmt.Errorf("error in AccountGet: %w", err)
	}

	var account models.AccountsGet
	account.ID = accountRepo.ID
	account.UserID = accountRepo.UserID
	account.Balance = accountRepo.Balance
	account.Currency = models.Currency(accountRepo.Currency)

	if a.cacheService != nil {
		if err := a.cacheService.Set(ctx, cacheKey, accountID, 10*time.Minute); err != nil {
			return models.AccountsGet{}, fmt.Errorf("error in accountGet: %w; cachekey %s not set", err, cacheKey)
		}
	}

	return account, nil
}

// удаление счета
func (a *AccountsService) AccountDelete(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	tx, err := a.accountsRepo.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error in AccountDelete: %w", err)
	}
	defer tx.Rollback()

	account, err := a.accountsRepo.GetAccountByIdTx(ctx, accountID, userID, tx)
	if err != nil {
		return fmt.Errorf("error in AccountDelete: %w", err)
	}

	if account.Balance != 0 {
		return fmt.Errorf("cannot delete account with non-zero balance")
	}

	err = a.accountsRepo.DeleteAccount(ctx, accountID, userID, tx)
	if err != nil {
		return fmt.Errorf("error in AccountDelete: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	cacheKey := fmt.Sprintf("user_account:%s", userID.String())

	if a.cacheService != nil {
		if err := a.cacheService.Delete(ctx, cacheKey); err != nil {
			return fmt.Errorf("error in AccountDelete: %w; cachekey %s not set", err, cacheKey)
		}
	}

	return nil
}
