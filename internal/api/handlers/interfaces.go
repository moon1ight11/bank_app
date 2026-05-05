package handlers

import "github.com/gin-gonic/gin"

// интерфейс для хэндлеров счетов
type AccountsHandlersInterface interface {
	CreateAccount(c *gin.Context)
	GetAllUserAccounts(c *gin.Context)
	GetAccountById(c *gin.Context)
	DeleteAccount(c *gin.Context)
}

// интерфейс для хэндлеров аутентификации
type AuthHandlersInterface interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
}

// интерфейс для хэндлеров транзакций
type TransactionsHandlersInterface interface {
	GetAllUserTransactions(c *gin.Context)
	GetAllAccountTransactions(c *gin.Context)
	GetTransactionByID(c *gin.Context)
	CreateIncomingTransaction(c *gin.Context)
	CreateOutcomingTransaction(c *gin.Context)
	CreateTransferTransaction(c *gin.Context)
}

// интерфейс для хэндлеров юзеров
type UsersHandlersInterface interface {
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	CreateAdminOrVerificator(c *gin.Context)
	GetAllUsers(c *gin.Context)
	ChangeRole(c *gin.Context)
}
