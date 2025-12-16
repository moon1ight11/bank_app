package services

import (
	"bank_app/internal/storage/repos/accounts"
	"bank_app/internal/storage/repos/operations"
	"fmt"

	"github.com/google/uuid"
)

type AccountsService struct {
	accountsRepo   *accounts.Repo
	operationsRepo *operations.Repo
}

func NewAccountsService(accountsRepo *accounts.Repo) *AccountsService {
	return &AccountsService{accountsRepo: accountsRepo}
}

// создание счета
func (a *AccountsService) AccountAdd(account accounts.Account) (uuid.UUID, error) {
	AccountID, err := a.accountsRepo.CreateAccount(account.OwnerID)
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
	err = a.accountsRepo.DeleteAccount(userID, accountID)
	if err != nil {
		return err
	}

	return nil
}

// пополнение счета
func (a *AccountsService) AccountIncoming(userID uuid.UUID, accountID uuid.UUID, amount float64) error {
	// получаем целевой аккаунт пополнения для проверки
	_, err := a.accountsRepo.GetAccountById(accountID, userID)
	if err != nil {
		return err
	}

	// если все ок - открываем транзакцию
	transaction, err := a.accountsRepo.DB.Begin()
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// пополняем баланс
	err = a.accountsRepo.BalanceIncoming(accountID, amount, transaction)
	if err != nil {
		return err
	}

	// запись в журнал операций
	_, err = a.operationsRepo.CreateOperation(userID, accountID, amount, "incoming", transaction)
	if err != nil {
		return err
	}

	// если все ок - подтверждаем транзакцию
	transaction.Commit()
	return nil
}

// списание со счета
func (a *AccountsService) AccountOutlay(userID uuid.UUID, accountID uuid.UUID, amount float64) error {
	// получаем целевой аккаунт списания для проверки
	account, err := a.accountsRepo.GetAccountById(accountID, userID)
	if err != nil {
		return err
	}

	// проверяем, достаточно ли там денег
	if account.Balance < amount {
		return fmt.Errorf("not enough funds")
	}

	// если все ок - открываем транзакцию
	transaction, err := a.accountsRepo.DB.Begin()
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// списываем деньги
	err = a.accountsRepo.BalanceOutlay(accountID, amount, transaction)
	if err != nil {
		return err
	}

	// запись в журнал операций
	_, err = a.operationsRepo.CreateOperation(userID, accountID, amount, "outlay", transaction)
	if err != nil {
		return err
	}

	// если все ок - подтверждаем транзакцию
	transaction.Commit()
	return nil
}

// перевод
func (a *AccountsService) AccountTransfer(
	userInID uuid.UUID,
	accountInID uuid.UUID,
	userOutID uuid.UUID,
	accountOutID uuid.UUID,
	amount float64,
) error {
	// получаем целевой аккаунт списания для проверки
	accountOut, err := a.accountsRepo.GetAccountById(accountOutID, userOutID)
	if err != nil {
		return err
	}

	// проверяем, достаточно ли там денег
	if accountOut.Balance < amount {
		return fmt.Errorf("not enough funds")
	}

	// получаем целевой аккаунт пополнения для проверки
	_, err = a.accountsRepo.GetAccountById(accountInID, userInID)
	if err != nil {
		return err
	}

	// если все ок - открываем транзакцию
	transaction, err := a.accountsRepo.DB.Begin()
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// списываем деньги
	err = a.accountsRepo.BalanceOutlay(accountOutID, amount, transaction)
	if err != nil {
		return err
	}

	// запись в журнал операций
	_, err = a.operationsRepo.CreateOperation(userOutID, accountOutID, amount, "outlay", transaction)
	if err != nil {
		return err
	}

	// пополняем баланс
	err = a.accountsRepo.BalanceIncoming(accountInID, amount, transaction)
	if err != nil {
		return err
	}

	// запись в журнал операций
	_, err = a.operationsRepo.CreateOperation(userInID, accountInID, amount, "incoming", transaction)
	if err != nil {
		return err
	}

	// если все ок - подтверждаем транзакцию
	transaction.Commit()
	return nil
}
