package main

import (
	"bank_app/internal/api"
	"bank_app/internal/api/handlers"
	"bank_app/internal/config"
	"bank_app/internal/jwt"
	"bank_app/internal/services"
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

	usersService := services.NewUsersService(usersRepo)
	accountsService := services.NewAccountsService(accountsRepo, transactionsRepo)
	transactionsService := services.NewTransactionsService(transactionsRepo, accountsRepo)

	authHandler := handlers.NewAuthHandler(usersService, jwtService)
	usersHandler := handlers.NewUsersHandler(usersService, jwtService)
	accountsHandler := handlers.NewAccountsHandler(accountsService, jwtService)
	transactionsHandler := handlers.NewTransactionsHandler(transactionsService, jwtService)

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
