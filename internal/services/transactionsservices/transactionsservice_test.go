package transactionsservice

import (
	"bank_app/internal/api/models"
	"bank_app/internal/storage/repos/accounts"
	"bank_app/internal/storage/repos/transactions"
	"bank_app/internal/storage/repos/users"
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTransactionsService() *TransactionsService {
	accountsRepo := accounts.NewAccountsRepo(testDB)
	transactionsRepo := transactions.NewTransactionsRepo(testDB)
	return NewTransactionsService(transactionsRepo, accountsRepo)
}

func createTestUserWithAccount(t *testing.T, currency models.Currency) (uuid.UUID, uuid.UUID) {
	t.Helper()
	usersRepo := users.NewUsersRepo(testDB)
	accountsRepo := accounts.NewAccountsRepo(testDB)

	userID, err := usersRepo.CreateUser(context.Background(),
		"Test", "User",
		fmt.Sprintf("test%s@example.com", uuid.New().String()[:8]),
		"+"+fmt.Sprintf("7999%07d", 1000000+len(make([]struct{}, 0)))[:11],
		"hash", "Basic")
	require.NoError(t, err)

	accountID, err := accountsRepo.CreateAccount(context.Background(), userID, string(currency))
	require.NoError(t, err)

	return userID, accountID
}

func TestIntegration_TransactionIncoming_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newTransactionsService()
	userID, accountID := createTestUserWithAccount(t, models.CurrencyRUB)

	transaction := models.Transaction{
		UserTo:    userID,
		AccountTo: accountID,
		Amount:    5000,
		Currency:  models.CurrencyRUB,
	}

	id, err := svc.TransactionIncoming(context.Background(), transaction)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	accountsRepo := accounts.NewAccountsRepo(testDB)
	acc, err := accountsRepo.GetAccountById(context.Background(), accountID, userID)
	require.NoError(t, err)
	assert.Equal(t, 5000, acc.Balance)
}

func TestIntegration_TransactionIncoming_WrongCurrency(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newTransactionsService()
	userID, accountID := createTestUserWithAccount(t, models.CurrencyRUB)

	transaction := models.Transaction{
		UserTo:    userID,
		AccountTo: accountID,
		Amount:    5000,
		Currency:  models.CurrencyUSD,
	}

	_, err := svc.TransactionIncoming(context.Background(), transaction)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "wrong currency")
}

func TestIntegration_TransactionIncoming_ZeroAmount(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newTransactionsService()
	userID, accountID := createTestUserWithAccount(t, models.CurrencyRUB)

	transaction := models.Transaction{
		UserTo:    userID,
		AccountTo: accountID,
		Amount:    0,
		Currency:  models.CurrencyRUB,
	}

	_, err := svc.TransactionIncoming(context.Background(), transaction)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "zero amount")
}

func TestIntegration_TransactionOutcoming_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newTransactionsService()
	userID, accountID := createTestUserWithAccount(t, models.CurrencyRUB)

	_, err := svc.TransactionIncoming(context.Background(), models.Transaction{
		UserTo:    userID,
		AccountTo: accountID,
		Amount:    10000,
		Currency:  models.CurrencyRUB,
	})
	require.NoError(t, err)

	id, err := svc.TransactionOutcoming(context.Background(), models.Transaction{
		UserFrom:    userID,
		AccountFrom: accountID,
		Amount:      3000,
		Currency:    models.CurrencyRUB,
	})
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	accountsRepo := accounts.NewAccountsRepo(testDB)
	acc, err := accountsRepo.GetAccountById(context.Background(), accountID, userID)
	require.NoError(t, err)
	assert.Equal(t, 7000, acc.Balance)
}

func TestIntegration_TransactionOutcoming_InsufficientFunds(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newTransactionsService()
	userID, accountID := createTestUserWithAccount(t, models.CurrencyRUB)

	_, err := svc.TransactionOutcoming(context.Background(), models.Transaction{
		UserFrom:    userID,
		AccountFrom: accountID,
		Amount:      1000,
		Currency:    models.CurrencyRUB,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not enough founds")
}

func TestIntegration_TransactionTransfer_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newTransactionsService()

	userFromID, accountFromID := createTestUserWithAccount(t, models.CurrencyRUB)
	userToID, accountToID := createTestUserWithAccount(t, models.CurrencyRUB)

	_, err := svc.TransactionIncoming(context.Background(), models.Transaction{
		UserTo:    userFromID,
		AccountTo: accountFromID,
		Amount:    10000,
		Currency:  models.CurrencyRUB,
	})
	require.NoError(t, err)

	id, err := svc.TransactionTransfer(context.Background(), models.Transaction{
		UserFrom:    userFromID,
		AccountFrom: accountFromID,
		UserTo:      userToID,
		AccountTo:   accountToID,
		Amount:      4000,
		Currency:    models.CurrencyRUB,
	})
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	accountsRepo := accounts.NewAccountsRepo(testDB)
	from, _ := accountsRepo.GetAccountById(context.Background(), accountFromID, userFromID)
	to, _ := accountsRepo.GetAccountById(context.Background(), accountToID, userToID)
	assert.Equal(t, 6000, from.Balance)
	assert.Equal(t, 4000, to.Balance)
}

func TestIntegration_TransactionTransfer_DifferentCurrencies(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newTransactionsService()

	userFromID, accountFromID := createTestUserWithAccount(t, models.CurrencyRUB)
	userToID, accountToID := createTestUserWithAccount(t, models.CurrencyUSD) 

	_, err := svc.TransactionIncoming(context.Background(), models.Transaction{
		UserTo:    userFromID,
		AccountTo: accountFromID,
		Amount:    10000,
		Currency:  models.CurrencyRUB,
	})
	require.NoError(t, err)

	_, err = svc.TransactionTransfer(context.Background(), models.Transaction{
		UserFrom:    userFromID,
		AccountFrom: accountFromID,
		UserTo:      userToID,
		AccountTo:   accountToID,
		Amount:      1000,
		Currency:    models.CurrencyRUB,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "different currencies")
}

func TestIntegration_TransactionTransfer_SameAccount(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newTransactionsService()
	userID, accountID := createTestUserWithAccount(t, models.CurrencyRUB)

	_, err := svc.TransactionIncoming(context.Background(), models.Transaction{
		UserTo:    userID,
		AccountTo: accountID,
		Amount:    10000,
		Currency:  models.CurrencyRUB,
	})
	require.NoError(t, err)

	_, err = svc.TransactionTransfer(context.Background(), models.Transaction{
		UserFrom:    userID,
		AccountFrom: accountID,
		UserTo:      userID,
		AccountTo:   accountID,
		Amount:      1000,
		Currency:    models.CurrencyRUB,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "transaction within one account")
}

func TestIntegration_TransactionTransfer_InsufficientFunds(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newTransactionsService()

	userFromID, accountFromID := createTestUserWithAccount(t, models.CurrencyRUB)
	userToID, accountToID := createTestUserWithAccount(t, models.CurrencyRUB)

	_, err := svc.TransactionTransfer(context.Background(), models.Transaction{
		UserFrom:    userFromID,
		AccountFrom: accountFromID,
		UserTo:      userToID,
		AccountTo:   accountToID,
		Amount:      1000,
		Currency:    models.CurrencyRUB,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not enough founds")
}

func TestIntegration_Transactions_GetAllUserTransactions(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newTransactionsService()
	userID, accountID := createTestUserWithAccount(t, models.CurrencyRUB)

	_, err := svc.TransactionIncoming(context.Background(), models.Transaction{
		UserTo:    userID,
		AccountTo: accountID,
		Amount:    5000,
		Currency:  models.CurrencyRUB,
	})
	require.NoError(t, err)

	transactions, err := svc.AllTransactionsGet(context.Background(), userID)
	require.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, 5000, transactions[0].Amount)
}

func TestIntegration_Transactions_GetByAccount(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newTransactionsService()
	userID, accountID := createTestUserWithAccount(t, models.CurrencyRUB)

	_, err := svc.TransactionIncoming(context.Background(), models.Transaction{
		UserTo:    userID,
		AccountTo: accountID,
		Amount:    3000,
		Currency:  models.CurrencyRUB,
	})
	require.NoError(t, err)

	transactions, err := svc.AccountTransactionsGet(context.Background(), userID, accountID)
	require.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, 3000, transactions[0].Amount)
}
