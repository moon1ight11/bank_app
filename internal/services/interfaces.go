package services

import (
	"bank_app/internal/api/models"
	"context"
	"github.com/google/uuid"
)

// интерфейс для сервисов юзеров
type UsersServiceInterface interface {
	UserCheck(ctx context.Context, phoneNumber string, userEmail string) (bool, error)
	UserVerification(ctx context.Context, User models.UserAutorization) (models.UserGet, error)
	UserAdd(ctx context.Context, User models.UserRegister) (uuid.UUID, error)
	AdminAdd(ctx context.Context, admin models.UserRegister) (uuid.UUID, error)
	UserGet(ctx context.Context, UserID uuid.UUID) (models.UserGet, error)
	UsersByRoleGet(ctx context.Context, role models.Role) ([]models.UserGet, error)
	UserUpdate(ctx context.Context, name *string, surname *string, password *string, email *string, phone *string, tz *string, ID uuid.UUID) error
	UserDelete(ctx context.Context, userID uuid.UUID) error
	RoleChange(ctx context.Context, userID uuid.UUID, role models.Role) error
}

// интерфейс для сервисов транзакций
type TransactionsServiceInterface interface {
	AllTransactionsGet(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error)
	AccountTransactionsGet(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) ([]models.Transaction, error)
	TransactionByIdGet(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) (models.Transaction, error)
	TransactionIncoming(ctx context.Context, transaction models.Transaction) (uuid.UUID, error)
	TransactionOutcoming(ctx context.Context, transaction models.Transaction) (uuid.UUID, error)
	TransactionTransfer(ctx context.Context, transaction models.Transaction) (uuid.UUID, error)
}

// интерфейс для сервисов счетов
type AccountsServiceInterface interface {
	AccountAdd(ctx context.Context, newAccount models.AccountCreate) (uuid.UUID, error)
	AllAccountsGet(ctx context.Context, userID uuid.UUID) ([]models.AccountsGet, error)
	AccountGet(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (models.AccountsGet, error)
	AccountDelete(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error
}
