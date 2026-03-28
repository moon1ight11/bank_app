package api

import (
	"bank_app/internal/api/handlers"
	"bank_app/internal/api/jwt"
	"bank_app/internal/api/middleware"
	"bank_app/internal/monitoring"
	"bank_app/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler         handlers.AuthHandlersInterface
	usersHandler        handlers.UsersHandlersInterface
	accountsHandler     handlers.AccountsHandlersInterface
	transactionsHandler handlers.TransactionsHandlersInterface
	ginEngine           *gin.Engine
}

func NewRouter(
	authHandler handlers.AuthHandlersInterface,
	usersHandler handlers.UsersHandlersInterface,
	accountsHandler handlers.AccountsHandlersInterface,
	transactionsHandler handlers.TransactionsHandlersInterface,
) *Router {
	return &Router{
		authHandler:         authHandler,
		usersHandler:        usersHandler,
		accountsHandler:     accountsHandler,
		transactionsHandler: transactionsHandler,
		ginEngine:           gin.Default(),
	}
}

func (r *Router) Init(jwtService jwt.TokenService, logger logger.Logger, metrics *monitoring.Metrics) {
	// MIDDLEWARE для CORS
	r.ginEngine.Use(middleware.CORS())

	// регистрация хэндлера для сбора метрик
	metrics.RegisterMetricsHandler(r.ginEngine, "/metrics")

	// MIDDLEWARE для сбора метрик
	r.ginEngine.Use(middleware.MetricsMiddleware(metrics))

	// группировка роутов
	authGroup := r.ginEngine.Group("/api/v1/auth")
	userGroup := r.ginEngine.Group("/api/v1/basic")
	verificatorGroup := r.ginEngine.Group("/api/v1/verificator")
	adminGroup := r.ginEngine.Group("/api/v1/admin")

	// MIDDLEWARE для общей аутентификации
	userGroup.Use(middleware.Auth(jwtService, logger))
	verificatorGroup.Use(middleware.Auth(jwtService, logger))
	adminGroup.Use(middleware.Auth(jwtService, logger))

	// MIDDLEWARE для проверки ролей
	userGroup.Use(middleware.AuthUser(logger))
	verificatorGroup.Use(middleware.AuthVerificator(logger))
	adminGroup.Use(middleware.AuthAdmin(logger))

	// АУТЕНТИФИКАЦИЯ //
	// регистрация
	authGroup.POST("/sign-up", r.authHandler.SignUp)
	// авторизация
	authGroup.POST("/sign-in", r.authHandler.SignIn)
	// разлогин
	authGroup.GET("/sign-out", r.authHandler.SignOut)

	// БАЗОВЫЙ // --- любые верифицированные пользователи
	// получение данных пользователя
	userGroup.GET("/user", r.usersHandler.GetUser)
	// обновление пользователя
	userGroup.PATCH("/user", r.usersHandler.UpdateUser)
	// удаление пользователя
	userGroup.DELETE("/user", r.usersHandler.DeleteUser)
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
	// обновление роли пользователя
	adminGroup.PATCH("/users", r.usersHandler.ChangeRole)
}

func (r *Router) GetEngine() *gin.Engine {
	return r.ginEngine
}
