package services

import (
	"bank_app/internal/storage/repos/accounts"
	"bank_app/internal/storage/repos/transactions"
	"fmt"
	"github.com/google/uuid"
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

// создание счета
func (a *AccountsService) AccountAdd(newAccount accounts.Account) (uuid.UUID, error) {
	AccountID, err := a.accountsRepo.CreateAccount(newAccount.UserID, newAccount.Currency)
	if err != nil {
		return uuid.Nil, err
	}

	return AccountID, nil
}

// вывод всех счетов пользователя
func (a *AccountsService) AllAccountsGet(userID uuid.UUID) ([]accounts.Account, error) {
	Accounts, err := a.accountsRepo.GetAccountsByUserId(userID)
	if err != nil {
		return nil, err
	}

	return Accounts, nil
}

// вывод одного счета пользователя
func (a *AccountsService) AccountGet(userID uuid.UUID, accountID uuid.UUID) (accounts.Account, error) {
	Account, err := a.accountsRepo.GetAccountById(accountID, userID)
	if err != nil {
		return accounts.Account{}, err
	}

	return Account, nil
}

// удаление счета
func (a *AccountsService) AccountDelete(userID uuid.UUID, accountID uuid.UUID) error {
	// получаем счет из БД
	account, err := a.accountsRepo.GetAccountById(accountID, userID)
	if err != nil {
		return err
	}

	// проверяем, чтобы на счету не осталось денег
	if account.Balance != 0 {
		return fmt.Errorf("cannot delete account with non-zero balance")
	}

	// если денег нет - удаляем счет
	err = a.accountsRepo.DeleteAccount(accountID, userID)
	if err != nil {
		return err
	}

	return nil
}