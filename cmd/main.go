package main

import (
	"bank_app/internal/api"
	accountshandlers "bank_app/internal/api/handlers/accounts"
	"bank_app/internal/api/handlers/authentification"
	transactionshandlers "bank_app/internal/api/handlers/transactions"
	usershandlers "bank_app/internal/api/handlers/users"
	"bank_app/internal/api/jwt"
	"bank_app/internal/config"
	"bank_app/internal/services/accountsservices"
	transactionsservice "bank_app/internal/services/transactionsservices"
	usersservice "bank_app/internal/services/usersservices"
	"bank_app/internal/storage"
	"bank_app/internal/storage/repos/accounts"
	"bank_app/internal/storage/repos/transactions"
	"bank_app/internal/storage/repos/users"
	"log"
)

func main() {
	// инициализация конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	// соединение с БД
	db, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatal("Failed to connect DB")
	}

	// применение миграций
	err = db.UpMigrations()
	if err != nil {
		log.Fatal("Failed to upping migrations", err)
	}

	// инициализация зависимостей
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)

	usersRepo := users.NewUsersRepo(db)
	accountsRepo := accounts.NewAccountsRepo(db)
	transactionsRepo := transactions.NewTransactionsRepo(db)

	usersService := usersservice.NewUsersService(usersRepo)
	accountsService := accountsservices.NewAccountsService(accountsRepo, transactionsRepo)
	transactionsService := transactionsservice.NewTransactionsService(transactionsRepo, accountsRepo)

	authHandler := authentificationhandlers.NewAuthHandler(usersService, jwtService)
	usersHandler := usershandlers.NewUsersHandler(usersService, jwtService)
	accountsHandler := accountshandlers.NewAccountsHandler(accountsService, jwtService)
	transactionsHandler := transactionshandlers.NewTransactionsHandler(transactionsService, jwtService)

	// инициализация роутера
	router := api.NewRouter(authHandler, usersHandler, accountsHandler, transactionsHandler)

	// инициализация роутов
	router.Init(jwtService)

	// запуск роутера
	err = router.Run()
	if err != nil {
		log.Fatal("Failed to run Gin router", err)
	}
}
