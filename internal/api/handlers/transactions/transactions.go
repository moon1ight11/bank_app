package transactionshandlers

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/monitoring"
	"bank_app/internal/services"
	"bank_app/pkg/logger"
)

type TransactionsHandler struct {
	transactionsService services.TransactionsServiceInterface
	jwtService          jwt.TokenService
	logger              logger.Logger
	metrics             *monitoring.Metrics
}

func NewTransactionsHandler(
	transactionsService services.TransactionsServiceInterface,
	jwtService jwt.TokenService,
	logger logger.Logger,
	metrics *monitoring.Metrics,
) *TransactionsHandler {
	return &TransactionsHandler{
		transactionsService: transactionsService,
		jwtService:          jwtService,
		logger:              logger,
		metrics:             metrics,
	}
}
