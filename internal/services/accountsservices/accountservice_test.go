package accountsservices

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

func newAccountsService() *AccountsService {
	accountsRepo := accounts.NewAccountsRepo(testDB)
	transactionsRepo := transactions.NewTransactionsRepo(testDB)
	return NewAccountsService(accountsRepo, transactionsRepo, nil)
}

func createTestUser(t *testing.T) uuid.UUID {
	t.Helper()
	usersRepo := users.NewUsersRepo(testDB)
	id, err := usersRepo.CreateUser(context.Background(), "Test", "User", fmt.Sprintf("test%s@example.com", uuid.New().String()[:8]), "+79990000000", "hash", "Basic")
	require.NoError(t, err)
	return id
}

func TestIntegration_AccountAdd_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newAccountsService()
	userID := createTestUser(t)

	account := models.AccountCreate{
		UserID:   userID,
		Currency: models.CurrencyRUB,
	}

	accountID, err := svc.AccountAdd(context.Background(), account)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, accountID)
}

func TestIntegration_AccountGet_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newAccountsService()
	userID := createTestUser(t)

	account := models.AccountCreate{
		UserID:   userID,
		Currency: models.CurrencyUSD,
	}

	accountID, err := svc.AccountAdd(context.Background(), account)
	require.NoError(t, err)

	found, err := svc.AccountGet(context.Background(), userID, accountID)
	require.NoError(t, err)
	assert.Equal(t, accountID, found.ID)
	assert.Equal(t, models.CurrencyUSD, found.Currency)
}

func TestIntegration_AllAccountsGet_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newAccountsService()
	userID := createTestUser(t)

	_, err := svc.AccountAdd(context.Background(), models.AccountCreate{UserID: userID, Currency: models.CurrencyRUB})
	require.NoError(t, err)
	_, err = svc.AccountAdd(context.Background(), models.AccountCreate{UserID: userID, Currency: models.CurrencyEUR})
	require.NoError(t, err)

	accounts, err := svc.AllAccountsGet(context.Background(), userID)
	require.NoError(t, err)
	assert.Len(t, accounts, 2)
}

func TestIntegration_AccountDelete_ZeroBalance(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newAccountsService()
	userID := createTestUser(t)

	account := models.AccountCreate{UserID: userID, Currency: models.CurrencyRUB}
	accountID, err := svc.AccountAdd(context.Background(), account)
	require.NoError(t, err)

	err = svc.AccountDelete(context.Background(), userID, accountID)
	require.NoError(t, err)

	_, err = svc.AccountGet(context.Background(), userID, accountID)
	require.Error(t, err)
}

func TestIntegration_AccountDelete_NonZeroBalance(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newAccountsService()
	userID := createTestUser(t)

	account := models.AccountCreate{UserID: userID, Currency: models.CurrencyRUB}
	accountID, err := svc.AccountAdd(context.Background(), account)
	require.NoError(t, err)

	_, err = testDB.DB.ExecContext(context.Background(),
		"UPDATE bank_app.accounts SET balance = 1000 WHERE id = $1", accountID)
	require.NoError(t, err)

	err = svc.AccountDelete(context.Background(), userID, accountID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "non-zero balance")
}
