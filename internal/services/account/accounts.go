package account

import (
	"bank_app/internal/api/models"
	"fmt"
	"github.com/google/uuid"
)

// создание счета
func (a *AccountsService) AccountAdd(newAccount models.AccountCreate) (uuid.UUID, error) {
	AccountID, err := a.accountsRepo.CreateAccount(newAccount.UserID, string(newAccount.Currency))
	if err != nil {
		return uuid.Nil, err
	}

	return AccountID, nil
}

// вывод всех счетов пользователя
func (a *AccountsService) AllAccountsGet(userID uuid.UUID) ([]models.AccountsGet, error) {
	Accounts, err := a.accountsRepo.GetAccountsByUserId(userID)
	if err != nil {
		return nil, err
	}

	var accounts []models.AccountsGet

	for i := range(Accounts) {
		var account models.AccountsGet

		account.ID = Accounts[i].ID
		account.UserID = Accounts[i].UserID
		account.Balance = Accounts[i].Balance
		account.Currency = models.Currency(string(Accounts[i].Currency))

		accounts = append(accounts, account)
	}

	return accounts, nil
}

// вывод одного счета пользователя
func (a *AccountsService) AccountGet(userID uuid.UUID, accountID uuid.UUID) (models.AccountsGet, error) {
	Account, err := a.accountsRepo.GetAccountById(accountID, userID)
	if err != nil {
		return models.AccountsGet{}, err
	}

	var account models.AccountsGet
	account.ID = Account.ID
	account.UserID = Account.UserID
	account.Balance = Account.Balance
	account.Currency = models.Currency(Account.Currency)

	return account, nil
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