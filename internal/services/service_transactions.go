package services

import (
	"bank_app/internal/storage/repos/accounts"
	"bank_app/internal/storage/repos/transactions"
	"fmt"

	"github.com/google/uuid"
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

// выполнение входящей транзакции
func (t *TransactionsService) TransactionIncoming(transaction transactions.Transaction) (uuid.UUID, error) {
	// находим счет, на который пополнение
	account, err := t.accountsRepo.GetAccountById(transaction.AccountTo, transaction.UserTo)
	if err != nil {
		return uuid.Nil, err
	}

	// проверяем, чтобы счет был в той же валюте, что указана в транзакции
	if account.Currency != transaction.Currency {
		return uuid.Nil, fmt.Errorf("wrong currency")
	}

	// находим cчет нулевого админа, с которого будет списано
	adminAccount, err := t.accountsRepo.GetAdminAccountByCurrency(transaction.Currency)
	if err != nil {
		return uuid.Nil, err
	}

	// открываем ТХ
	tx, err := t.accountsRepo.DB.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	// отложенно откатываем ТХ
	defer tx.Rollback()

	// выполняем списание с админского аккаунта и пополнение пользовательского
	err = t.accountsRepo.BalanceOutcoming(adminAccount.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, err
	}

	err = t.accountsRepo.BalanceIncoming(account.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, err
	}

	// делаем запись в транзакции
	transactionID, err := t.transactionsRepo.CreateTransaction(
		adminAccount.UserID,
		adminAccount.ID,
		account.UserID,
		account.ID,
		transaction.Amount,
		transaction.Currency,
		tx,
	)
	if err != nil {
		return uuid.Nil, err
	}

	// если все ок - подтверждаем транзакцию
	tx.Commit()
	return transactionID, nil
}

// выполнение исходящей транзакции
func (t *TransactionsService) TransactionOutcoming(transaction transactions.Transaction) (uuid.UUID, error) {
	// находим счет, с которого списываем
	account, err := t.accountsRepo.GetAccountById(transaction.AccountFrom, transaction.UserFrom)
	if err != nil {
		return uuid.Nil, err
	}

	// проверяем, чтобы счет был в той же валюте, что указана в транзакции
	if account.Currency != transaction.Currency {
		return uuid.Nil, fmt.Errorf("wrong currency")
	}

	// находим cчет нулевого админа, которому будет зачислено
	adminAccount, err := t.accountsRepo.GetAdminAccountByCurrency(transaction.Currency)
	if err != nil {
		return uuid.Nil, err
	}

	// открываем ТХ
	tx, err := t.accountsRepo.DB.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	// отложенно откатываем ТХ
	defer tx.Rollback()

	// проверка, хватает ли на пользовательском аккаунте средств
	err = t.accountsRepo.BalanceCheck(account.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, err
	}

	// выполняем списание с пользовательского аккаунта и пополнение админского
	err = t.accountsRepo.BalanceOutcoming(account.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, err
	}

	err = t.accountsRepo.BalanceIncoming(adminAccount.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, err
	}

	// делаем запись в транзакции
	transactionID, err := t.transactionsRepo.CreateTransaction(
		adminAccount.UserID,
		adminAccount.ID,
		account.UserID,
		account.ID,
		transaction.Amount,
		transaction.Currency,
		tx,
	)
	if err != nil {
		return uuid.Nil, err
	}

	// если все ок - подтверждаем транзакцию
	tx.Commit()
	return transactionID, nil
}

// выполнение трансфера
func (t *TransactionsService) TransactionTransfer(transaction transactions.Transaction) (uuid.UUID, error) {
	// находим счет, с которого списываем
	accountFrom, err := t.accountsRepo.GetAccountById(transaction.AccountFrom, transaction.UserFrom)
	if err != nil {
		return uuid.Nil, err
	}

	// находим счет, на который пополнение
	accountTo, err := t.accountsRepo.GetAccountById(transaction.AccountTo, transaction.UserTo)
	if err != nil {
		return uuid.Nil, err
	}

	// проверяем, что одинаковые валюты
	if accountFrom.Currency != accountTo.Currency {
		return uuid.Nil, err
	}

	// открываем ТХ
	tx, err := t.accountsRepo.DB.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	// отложенно откатываем ТХ
	defer tx.Rollback()

	// проверка, хватает ли на исходящем аккаунте средств
	err = t.accountsRepo.BalanceCheck(accountFrom.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, err
	}

	// выполняем списание с одного аккаунта и пополнение другого
	err = t.accountsRepo.BalanceOutcoming(accountFrom.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, err
	}

	err = t.accountsRepo.BalanceIncoming(accountTo.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, err
	}

	// делаем запись в транзакции
	transactionID, err := t.transactionsRepo.CreateTransaction(
		accountFrom.UserID,
		accountFrom.ID,
		accountTo.UserID,
		accountTo.ID,
		transaction.Amount,
		transaction.Currency,
		tx,
	)
	if err != nil {
		return uuid.Nil, err
	}

	// если все ок - подтверждаем транзакцию
	tx.Commit()
	return transactionID, nil
}
