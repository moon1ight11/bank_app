package main

import (
	"bank_app/internal/api"
	accountshandlers "bank_app/internal/api/handlers/accounts"
	"bank_app/internal/api/handlers/authentification"
	transactionshandlers "bank_app/internal/api/handlers/transactions"
	usershandlers "bank_app/internal/api/handlers/users"
	"bank_app/internal/api/jwt"
	"bank_app/internal/config"
	"bank_app/internal/monitoring"
	"bank_app/internal/services/accountsservices"
	transactionsservice "bank_app/internal/services/transactionsservices"
	usersservice "bank_app/internal/services/usersservices"
	"bank_app/internal/storage"
	"bank_app/internal/storage/cache"
	"bank_app/internal/storage/repos/accounts"
	"bank_app/internal/storage/repos/transactions"
	"bank_app/internal/storage/repos/users"
	"bank_app/pkg/logger"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// обработка паники
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic recovered: %v\n", r)
			os.Exit(1)
		}
	}()

	// инициализация конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: %w", err)
	}

	// инициализация логгера
	logger, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		panic(err)
	}

	// инициализация метрик
	metrics := monitoring.NewMetrics()

	// соединение с БД
	db, err := storage.NewStorage(cfg)
	if err != nil {
		logger.Error("Failed to connect DB", "error:", err)
		os.Exit(1)
	}

	// применение миграций
	err = db.UpMigrations()
	if err != nil {
		logger.Error("Failed to upping migrations", "error:", err)
		os.Exit(1)
	}

	// подключаемся к редис
	redisClient, err := storage.NewRedisClient(cfg)
	if err != nil {
		logger.Error("Failed to connect to Redis", "error:", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	// инициализируем кэш
	var cacheService *cache.CacheService
	if redisClient != nil {
		cacheService = cache.NewCacheService(redisClient.Client)
	}

	// инициализация зависимостей
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)

	usersRepo := users.NewUsersRepo(db)
	accountsRepo := accounts.NewAccountsRepo(db)
	transactionsRepo := transactions.NewTransactionsRepo(db)

	usersService := usersservice.NewUsersService(usersRepo, cacheService)
	accountsService := accountsservices.NewAccountsService(accountsRepo, transactionsRepo, cacheService)
	transactionsService := transactionsservice.NewTransactionsService(transactionsRepo, accountsRepo)

	authHandler := authentificationhandlers.NewAuthHandler(usersService, jwtService, logger, metrics)
	usersHandler := usershandlers.NewUsersHandler(usersService, jwtService, logger, metrics)
	accountsHandler := accountshandlers.NewAccountsHandler(accountsService, jwtService, logger, metrics)
	transactionsHandler := transactionshandlers.NewTransactionsHandler(transactionsService, jwtService, logger, metrics)

	// инициализация роутера
	router := api.NewRouter(authHandler, usersHandler, accountsHandler, transactionsHandler)

	// инициализация роутов
	router.Init(jwtService, logger, metrics)

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router.GetEngine(),
	}

	// создаем каналы для сигналов завершения и ошибки сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	serverError := make(chan error, 1)

	// запуск сервера
	go func() {
		logger.Info("Server is starting on port :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to run server:", "error", err)
			serverError <- err
		}
	}()

	// ждем сигналы
	select {
	case <-quit:
		// если поступил сигнал завершения - делаем шатдаун с таймаутом
		logger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Server forced to shutdown:", "error", err)
		}
	case err := <-serverError:
		// еслм пришла ошибка от сервера - фаталим
		logger.Error("Server error:", "error", err)
	}

	logger.Info("Server exited")
}
