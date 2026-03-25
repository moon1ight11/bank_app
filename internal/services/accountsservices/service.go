package accountsservices

import (
	"bank_app/internal/api/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// создание счета
func (a *AccountsService) AccountAdd(ctx context.Context, newAccount models.AccountCreate) (uuid.UUID, error) {
	// создаем счет
	accountID, err := a.accountsRepo.CreateAccount(ctx, newAccount.UserID, string(newAccount.Currency))
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in AccountAdd: %w", err)
	}

	return accountID, nil
}

// вывод всех счетов пользователя
func (a *AccountsService) AllAccountsGet(ctx context.Context, userID uuid.UUID) ([]models.AccountsGet, error) {
	accountsRepo, err := a.accountsRepo.GetAccountsByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error in AllAccountsGet: %w", err)
	}

	var accounts []models.AccountsGet

	for i := range(accountsRepo) {
		var account models.AccountsGet

		account.ID = accountsRepo[i].ID
		account.UserID = accountsRepo[i].UserID
		account.Balance = accountsRepo[i].Balance
		account.Currency = models.Currency(string(accountsRepo[i].Currency))

		accounts = append(accounts, account)
	}

	return accounts, nil
}

// вывод одного счета пользователя
func (a *AccountsService) AccountGet(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (models.AccountsGet, error) {
	accountRepo, err := a.accountsRepo.GetAccountById(ctx, accountID, userID)
	if err != nil {
		return models.AccountsGet{}, fmt.Errorf("error in AccountGet: %w", err)
	}

	var account models.AccountsGet
	account.ID = accountRepo.ID
	account.UserID = accountRepo.UserID
	account.Balance = accountRepo.Balance
	account.Currency = models.Currency(accountRepo.Currency)

	return account, nil
}

// удаление счета
func (a *AccountsService) AccountDelete(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	// нужно сделать атомарно !!!
	
	
	// получаем счет из БД
	account, err := a.accountsRepo.GetAccountById(ctx, accountID, userID)
	if err != nil {
		return fmt.Errorf("error in AccountDelete: %w", err)
	}

	// проверяем, чтобы на счету не осталось денег
	if account.Balance != 0 {
		return fmt.Errorf("cannot delete account with non-zero balance")
	}

	// если денег нет - удаляем счет
	err = a.accountsRepo.DeleteAccount(ctx, accountID, userID)
	if err != nil {
		return fmt.Errorf("error in AccountDelete: %w", err)
	}

	return nil
}