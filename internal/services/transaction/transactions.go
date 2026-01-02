package transaction

import (
	"bank_app/internal/api/models"
	"fmt"
	"github.com/google/uuid"
)

// получение всех транзакций пользователя
func (t *TransactionsService) AllTransactionsGet(userID uuid.UUID) ([]models.Transaction, error) {
	transactionsRepo, err := t.transactionsRepo.GetAllUsersTransactions(userID)
	if err != nil {
		return nil, err
	}

	transactionsApi := make([]models.Transaction, 0, len(transactionsRepo))

	for i := range transactionsRepo {
		// Создаем новую транзакцию
		transactionApi := models.Transaction{
			ID:          transactionsRepo[i].ID,
			UserFrom:    transactionsRepo[i].UserFrom,
			AccountFrom: transactionsRepo[i].AccountFrom,
			UserTo:      transactionsRepo[i].UserTo,
			AccountTo:   transactionsRepo[i].AccountTo,
			Amount:      transactionsRepo[i].Amount,
			Timestamp:   transactionsRepo[i].Timestamp,
			Currency:    models.Currency(transactionsRepo[i].Currency),
		}

		transactionsApi = append(transactionsApi, transactionApi)
	}

	return transactionsApi, nil
}

// получение транзакций по счету
func (t *TransactionsService) AccountTransactionsGet(userID uuid.UUID, accountID uuid.UUID) ([]models.Transaction, error) {
	transactionsRepo, err := t.transactionsRepo.GetTransactionsByAccount(userID, accountID)
	if err != nil {
		return nil, err
	}

	transactionsApi := make([]models.Transaction, 0, len(transactionsRepo))

	for i := range transactionsRepo {
		// Создаем новую транзакцию
		transactionApi := models.Transaction{
			ID:          transactionsRepo[i].ID,
			UserFrom:    transactionsRepo[i].UserFrom,
			AccountFrom: transactionsRepo[i].AccountFrom,
			UserTo:      transactionsRepo[i].UserTo,
			AccountTo:   transactionsRepo[i].AccountTo,
			Amount:      transactionsRepo[i].Amount,
			Timestamp:   transactionsRepo[i].Timestamp,
			Currency:    models.Currency(transactionsRepo[i].Currency),
		}

		transactionsApi = append(transactionsApi, transactionApi)
	}

	return transactionsApi, nil
}

// получение одной транзакции
func (t *TransactionsService) TransactionByIdGet(userID uuid.UUID, transactionID uuid.UUID) (models.Transaction, error) {
	transacionRepo, err := t.transactionsRepo.GetTransactionByID(transactionID, userID)
	if err != nil {
		return models.Transaction{}, err
	}

	var transactionApi models.Transaction

	transactionApi.ID = transacionRepo.ID
	transactionApi.UserFrom = transacionRepo.UserFrom
	transactionApi.AccountFrom = transacionRepo.AccountFrom
	transactionApi.UserTo = transacionRepo.UserTo
	transactionApi.AccountTo = transacionRepo.AccountTo
	transactionApi.Amount = transacionRepo.Amount
	transactionApi.Timestamp = transacionRepo.Timestamp
	transactionApi.Currency = models.Currency(transacionRepo.Currency)

	return transactionApi, nil
}

// выполнение входящей транзакции
func (t *TransactionsService) TransactionIncoming(transaction models.Transaction) (uuid.UUID, error) {
	// находим счет, на который пополнение
	account, err := t.accountsRepo.GetAccountById(transaction.AccountTo, transaction.UserTo)
	if err != nil {
		return uuid.Nil, err
	}

	// проверяем, чтобы счет был в той же валюте, что указана в транзакции
	if models.Currency(account.Currency) != transaction.Currency {
		return uuid.Nil, fmt.Errorf("wrong currency")
	}

	// проверяем, что переводится какая-то сумма
	if transaction.Amount <= 0 {
		return uuid.Nil, fmt.Errorf("zero amount")
	}

	// находим cчет нулевого админа, с которого будет списано
	adminAccount, err := t.accountsRepo.GetAdminAccountByCurrency(string(transaction.Currency))
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
		string(transaction.Currency),
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
func (t *TransactionsService) TransactionOutcoming(transaction models.Transaction) (uuid.UUID, error) {
	// находим счет, с которого списываем
	account, err := t.accountsRepo.GetAccountById(transaction.AccountFrom, transaction.UserFrom)
	if err != nil {
		return uuid.Nil, err
	}

	// проверяем, чтобы счет был в той же валюте, что указана в транзакции
	if models.Currency(account.Currency)  != transaction.Currency {
		return uuid.Nil, fmt.Errorf("wrong currency")
	}

	// проверяем, что переводится какая-то сумма
	if transaction.Amount <= 0 {
		return uuid.Nil, fmt.Errorf("zero amount")
	}

	// находим cчет нулевого админа, которому будет зачислено
	adminAccount, err := t.accountsRepo.GetAdminAccountByCurrency(string(transaction.Currency))
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
		string(transaction.Currency),
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
func (t *TransactionsService) TransactionTransfer(transaction models.Transaction) (uuid.UUID, error) {
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

	// проверяем, что переводится какая-то сумма
	if transaction.Amount <= 0 {
		return uuid.Nil, fmt.Errorf("zero amount")
	}

	// проверяем, чтобы средства не переводились в рамках одного счета
	if accountFrom.ID == accountTo.ID {
		return uuid.Nil, fmt.Errorf("transaction within one account")
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
		string(transaction.Currency),
		tx,
	)
	if err != nil {
		return uuid.Nil, err
	}

	// если все ок - подтверждаем транзакцию
	tx.Commit()
	return transactionID, nil
}
