package transactionsservice

import (
	"bank_app/internal/api/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// получение всех транзакций пользователя
func (t *TransactionsService) AllTransactionsGet(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	transactionsRepo, err := t.transactionsRepo.GetAllUsersTransactions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error in AllTransactionsGet: %w", err)
	}

	transactionsApi := make([]models.Transaction, 0, len(transactionsRepo))

	for i := range transactionsRepo {
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
func (t *TransactionsService) AccountTransactionsGet(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) ([]models.Transaction, error) {
	transactionsRepo, err := t.transactionsRepo.GetTransactionsByAccount(ctx, userID, accountID)
	if err != nil {
		return nil, fmt.Errorf("error in AccountTransactionsGet: %w", err)
	}

	transactionsApi := make([]models.Transaction, 0, len(transactionsRepo))

	for i := range transactionsRepo {
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
func (t *TransactionsService) TransactionByIdGet(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) (models.Transaction, error) {
	transacionRepo, err := t.transactionsRepo.GetTransactionByID(ctx, transactionID, userID)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("error in TransactionByIdGet: %w", err)
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
func (t *TransactionsService) TransactionIncoming(ctx context.Context, transaction models.Transaction) (uuid.UUID, error) {
	// находим счет, на который пополнение
	account, err := t.accountsRepo.GetAccountById(ctx, transaction.AccountTo, transaction.UserTo)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionIncoming: %w", err)
	}

	// проверяем, чтобы счет был в той же валюте, что указана в транзакции
	if models.Currency(account.Currency) != transaction.Currency {
		return uuid.Nil, fmt.Errorf("error in TransactionIncoming: wrong currency")
	}

	// проверяем, что переводится какая-то сумма
	if transaction.Amount <= 0 {
		return uuid.Nil, fmt.Errorf("error in TransactionIncoming: zero amount")
	}

	// находим cчет нулевого админа, с которого будет списано
	adminAccount, err := t.accountsRepo.GetAdminAccountByCurrency(ctx, string(transaction.Currency))
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionIncoming: %w", err)
	}

	// открываем ТХ
	tx, err := t.accountsRepo.DB.Begin()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionIncoming: %w", err)
	}

	// отложенно откатываем ТХ
	defer tx.Rollback()

	// выполняем списание с админского аккаунта и пополнение пользовательского
	err = t.accountsRepo.BalanceOutcoming(ctx, adminAccount.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionIncoming: %w", err)
	}

	err = t.accountsRepo.BalanceIncoming(ctx, account.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionIncoming: %w", err)
	}

	// делаем запись в транзакции
	transactionID, err := t.transactionsRepo.CreateTransaction(
		ctx,
		adminAccount.UserID,
		adminAccount.ID,
		account.UserID,
		account.ID,
		transaction.Amount,
		string(transaction.Currency),
		tx,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionIncoming: %w", err)
	}

	// если все ок - подтверждаем транзакцию
	tx.Commit()
	return transactionID, nil
}

// выполнение исходящей транзакции
func (t *TransactionsService) TransactionOutcoming(ctx context.Context, transaction models.Transaction) (uuid.UUID, error) {
	// находим счет, с которого списываем
	account, err := t.accountsRepo.GetAccountById(ctx, transaction.AccountFrom, transaction.UserFrom)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionOutcoming: %w", err)
	}

	// проверяем, чтобы счет был в той же валюте, что указана в транзакции
	if models.Currency(account.Currency)  != transaction.Currency {
		return uuid.Nil, fmt.Errorf("error in TransactionOutcoming: wrong currency")
	}

	// проверяем, что переводится какая-то сумма
	if transaction.Amount <= 0 {
		return uuid.Nil, fmt.Errorf("error in TransactionOutcoming: zero amount")
	}

	// находим cчет нулевого админа, которому будет зачислено
	adminAccount, err := t.accountsRepo.GetAdminAccountByCurrency(ctx, string(transaction.Currency))
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionOutcoming: %w", err)
	}

	// открываем ТХ
	tx, err := t.accountsRepo.DB.Begin()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionOutcoming: %w", err)
	}

	// отложенно откатываем ТХ
	defer tx.Rollback()

	// проверка, хватает ли на пользовательском аккаунте средств
	err = t.accountsRepo.BalanceCheck(ctx, account.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionOutcoming: %w", err)
	}

	// выполняем списание с пользовательского аккаунта и пополнение админского
	err = t.accountsRepo.BalanceOutcoming(ctx, account.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionOutcoming: %w", err)
	}

	err = t.accountsRepo.BalanceIncoming(ctx, adminAccount.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionOutcoming: %w", err)
	}

	// делаем запись в транзакции
	transactionID, err := t.transactionsRepo.CreateTransaction(
		ctx,
		adminAccount.UserID,
		adminAccount.ID,
		account.UserID,
		account.ID,
		transaction.Amount,
		string(transaction.Currency),
		tx,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionOutcoming: %w", err)
	}

	// если все ок - подтверждаем транзакцию
	tx.Commit()
	return transactionID, nil
}

// выполнение трансфера
func (t *TransactionsService) TransactionTransfer(ctx context.Context, transaction models.Transaction) (uuid.UUID, error) {
	// находим счет, с которого списываем
	accountFrom, err := t.accountsRepo.GetAccountById(ctx, transaction.AccountFrom, transaction.UserFrom)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionTransfer: %w", err)
	}

	// находим счет, на который пополнение
	accountTo, err := t.accountsRepo.GetAccountById(ctx, transaction.AccountTo, transaction.UserTo)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionTransfer: %w", err)
	}

	// проверяем, что одинаковые валюты
	if accountFrom.Currency != accountTo.Currency {
		return uuid.Nil, fmt.Errorf("error in TransactionTransfer: different currencies")
	}

	// проверяем, что переводится какая-то сумма
	if transaction.Amount <= 0 {
		return uuid.Nil, fmt.Errorf("error in TransactionTransfer: zero amount")
	}

	// проверяем, чтобы средства не переводились в рамках одного счета
	if accountFrom.ID == accountTo.ID {
		return uuid.Nil, fmt.Errorf("error in TransactionTransfer: transaction within one account")
	}

	// открываем ТХ
	tx, err := t.accountsRepo.DB.Begin()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionTransfer: %w", err)
	}

	// отложенно откатываем ТХ
	defer tx.Rollback()

	// проверка, хватает ли на исходящем аккаунте средств
	err = t.accountsRepo.BalanceCheck(ctx, accountFrom.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionTransfer: %w", err)
	}

	// выполняем списание с одного аккаунта и пополнение другого
	err = t.accountsRepo.BalanceOutcoming(ctx, accountFrom.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionTransfer: %w", err)
	}

	err = t.accountsRepo.BalanceIncoming(ctx, accountTo.ID, transaction.Amount, tx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionTransfer: %w", err)
	}

	// делаем запись в транзакции
	transactionID, err := t.transactionsRepo.CreateTransaction(
		ctx,
		accountFrom.UserID,
		accountFrom.ID,
		accountTo.UserID,
		accountTo.ID,
		transaction.Amount,
		string(transaction.Currency),
		tx,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in TransactionTransfer: %w", err)
	}

	// если все ок - подтверждаем транзакцию
	tx.Commit()
	return transactionID, nil
}
