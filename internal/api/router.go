package api

import (
	"bank_app/internal/api/handlers"
	"bank_app/internal/api/jwt"
	"bank_app/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler         *handlers.AuthHandler
	usersHandler        *handlers.UsersHandler
	accountsHandler     *handlers.AccountsHandler
	transactionsHandler *handlers.TransactionsHandler
	ginEngine           *gin.Engine
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	usersHandler *handlers.UsersHandler,
	accountsHandler *handlers.AccountsHandler,
	transactionsHandler *handlers.TransactionsHandler,
) *Router {
	return &Router{
		authHandler:         authHandler,
		usersHandler:        usersHandler,
		accountsHandler:     accountsHandler,
		transactionsHandler: transactionsHandler,
		ginEngine:           gin.Default(),
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
	// MIDDLEWARE для CORS
	r.ginEngine.Use(middleware.CORS())

	// группировка роутов
	authGroup := r.ginEngine.Group("/api/v1/auth")
	userGroup := r.ginEngine.Group("/api/v1/user")
	verificatorGroup := r.ginEngine.Group("/api/v1/verificator")
	adminGroup := r.ginEngine.Group("/api/v1/admin")

	// MIDDLEWARE для общей аутентификации
	userGroup.Use(middleware.Auth(jwtService))
	verificatorGroup.Use(middleware.Auth(jwtService))
	adminGroup.Use(middleware.Auth(jwtService))

	// MIDDLEWARE для проверки ролей
	userGroup.Use(middleware.AuthUser())
	verificatorGroup.Use(middleware.AuthVerificator())
	adminGroup.Use(middleware.AuthAdmin())

	// АУТЕНТИФИКАЦИЯ //
	// регистрация
	authGroup.POST("/sign-up", r.authHandler.SignUp)
	// авторизация
	authGroup.POST("/sign-in", r.authHandler.SignIn)
	// разлогин
	authGroup.GET("/sign-out", r.authHandler.SignOut)

	// ЮЗЕР // --- любые верифицированные пользователи
	// получение данных пользователя
	userGroup.GET("/profile", r.usersHandler.GetUser)
	// обновление пользователя
	userGroup.PATCH("/profile", r.usersHandler.UpdateUser)
	// удаление пользователя
	userGroup.DELETE("/profile", r.usersHandler.DeleteUser)
	// создание нового счета
	userGroup.POST("/accounts", r.accountsHandler.CreateAccount)
	// список счетов пользвателя
	userGroup.GET("/accounts", r.accountsHandler.GetAllUserAccounts)
	// получение информации о счете
	userGroup.GET("/accounts/:account_id", r.accountsHandler.GetAccountById)
	// удаление счета
	userGroup.DELETE("/accounts/:account_id", r.accountsHandler.DeleteAccount)
	// получение транзакций по счёту
	userGroup.GET("/accounts/:account_id/transactions", r.transactionsHandler.GetAllAccountTransactions)
	// получение всех транзакций пользователя
	userGroup.GET("/transactions", r.transactionsHandler.GetAllUserTransactions)
	// информация о конкретной транзакции
	userGroup.GET("/transactions/:transaction_id", r.transactionsHandler.GetTransactionByID)
	// перевод внутри системы
	userGroup.POST("/transactions/transfer", r.transactionsHandler.CreateTransferTransaction) /// --- для админов и верификаторов недоступно

	// ВЕРИФИКАТОР // --- только верификаторы и админы
	// пополнение счета
	verificatorGroup.POST("/accounts/transactions/income", r.transactionsHandler.CreateIncomingTransaction)
	// списание со счета
	verificatorGroup.POST("/accounts/transactions/outcome", r.transactionsHandler.CreateOutcomingTransaction)

	// АДМИН // --- только админы
	// создание админа или верификатора
	adminGroup.POST("/users", r.usersHandler.CreateAdminOrVerificator)
	// список админов или верификаторов
	adminGroup.GET("/users", r.usersHandler.GetAllUsers)
}
