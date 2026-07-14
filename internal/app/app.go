package app

import (
	"bank_app/internal/api"
	accountshandlers "bank_app/internal/api/handlers/accounts"
	authentificationhandlers "bank_app/internal/api/handlers/authentification"
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
	"fmt"
	"log"
)

type Dependencies struct {
	Router  *api.Router
	DB      *storage.DataBase
	Redis   *storage.RedisClient
	Logger  logger.Logger
	Metrics *monitoring.Metrics
}

func InitDependencies(cfg *config.Config) *Dependencies {
	// логгер
	logger, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	// метрики
	metrics := monitoring.NewMetrics()

	// база данных
	db, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	// миграции
	if err := db.UpMigrations(); err != nil {
		log.Fatalf("Failed to upping migrations: %v", err)
	}

	// редис
	redisClient, err := storage.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// кэш
	var cacheService *cache.CacheService
	if redisClient != nil {
		cacheService = cache.NewCacheService(redisClient.Client)
	}

	// jwt
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)

	// слой репозиториев
	usersRepo := users.NewUsersRepo(db)
	accountsRepo := accounts.NewAccountsRepo(db)
	transactionsRepo := transactions.NewTransactionsRepo(db)

	// сервисный слой
	usersService := usersservice.NewUsersService(usersRepo, cacheService)
	accountsService := accountsservices.NewAccountsService(accountsRepo, transactionsRepo, cacheService)
	transactionsService := transactionsservice.NewTransactionsService(transactionsRepo, accountsRepo)

	// слой хэндлеров
	authHandler := authentificationhandlers.NewAuthHandler(usersService, jwtService, logger, metrics)
	usersHandler := usershandlers.NewUsersHandler(usersService, jwtService, logger, metrics)
	accountsHandler := accountshandlers.NewAccountsHandler(accountsService, jwtService, logger, metrics)
	transactionsHandler := transactionshandlers.NewTransactionsHandler(transactionsService, jwtService, logger, metrics)

	// роутер
	router := api.NewRouter(authHandler, usersHandler, accountsHandler, transactionsHandler)
	router.Init(jwtService, logger, metrics)

	return &Dependencies{
		Router:  router,
		DB:      db,
		Redis:   redisClient,
		Logger:  logger,
		Metrics: metrics,
	}
}

func (d *Dependencies) Close() error {
	var errs []error

	if d.Redis != nil {
		if err := d.Redis.Close(); err != nil {
			errs = append(errs, fmt.Errorf("redis close: %w", err))
		}
	}

	if d.DB != nil {
		if err := d.DB.DB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("db close: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during close: %v", errs)
	}

	return nil
}
