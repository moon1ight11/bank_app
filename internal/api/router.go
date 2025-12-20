package api

import (
	"bank_app/internal/api/handlers"
	"bank_app/internal/jwt"
	"bank_app/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler       *handlers.AuthHandler
	usersHandler      *handlers.UsersHandler
	accountsHandler   *handlers.AccountsHandler
	transactionsHandler *handlers.TransactionsHandler
	ginEngine         *gin.Engine
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	usersHandler *handlers.UsersHandler,
	accountsHandler *handlers.AccountsHandler,
	transactionsHandler *handlers.TransactionsHandler,
) *Router {
	return &Router{
		authHandler:       authHandler,
		usersHandler:      usersHandler,
		accountsHandler:   accountsHandler,
		transactionsHandler: transactionsHandler,
		ginEngine:         gin.Default(),
	}
}

func (r *Router) Run() error {
	err := r.ginEngine.Run(":8080")
	if err != nil {
		return err
	}

	return nil
}

func (r *Router) Init(jwtService jwt.TokenService) {
	// MIDDLEWARE для корсов
	r.ginEngine.Use(middleware.CORS())

	// группировка роутов
	privateGroup := r.ginEngine.Group("/api/v1/private")
	authGroup := r.ginEngine.Group("/api/v1/auth")

	// MIDDLEWARE для аутентификации //
	// privateGroup.Use(middleware.Auth(jwtService))

	// АУТЕНТИФИКАЦИЯ //
	// регистрация
	authGroup.POST("/sign-up", r.authHandler.SignUp)
	// авторизация
	authGroup.POST("/sign-in", r.authHandler.SignIn)
	// разлогин
	authGroup.GET("/sign-out", r.authHandler.SignOut)

	// ЮЗЕРЫ //	
	// получение данных пользователя
	privateGroup.GET("/users", r.usersHandler.GetUser)
	// обновление пользователя
	privateGroup.PATCH("/users", r.usersHandler.UpdateUser)
	// удаление пользователя
	privateGroup.DELETE("/users", r.usersHandler.DeleteUser)

	// СЧЕТА //
	// создание нового счета
	privateGroup.POST("/accounts", r.accountsHandler.CreateAccount)
	// список счетов пользвателя
	privateGroup.GET("/accounts", r.accountsHandler.GetAllUserAccounts)
	// получение информации о счете
	privateGroup.GET("/accounts/:account_id", r.accountsHandler.GetAccountById)
	// удаление счета
	privateGroup.DELETE("/accounts/:account_id", r.accountsHandler.DeleteAccount)
	
	// Админ-группа:
	// api/v1/private/admin
	// /accounts/transactions/ - POST. В теле - структрура с полем типа. 

	// ТРАНЗАКЦИИ //
	// получение всех транзакций пользователя
	privateGroup.GET("/transactions", r.transactionsHandler.GetAllAccountTransactions)
	// получение транзакций по конкретному счёту
	privateGroup.GET("/transactions/:account_id", r.transactionsHandler.GetAllAccountTransactions)
	// информация о конкретной транзакции
	privateGroup.GET("/transactions/info/:transaction_id", r.transactionsHandler.GetTransactionByID)
}