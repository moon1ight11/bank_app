package transactionshandlers

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/services"
	"bank_app/pkg/logger"
)

type TransactionsHandler struct {
	transactionsService services.TransactionsServiceInterface
	jwtService          jwt.TokenService
	logger              logger.Logger
}

func NewTransactionsHandler(transactionsService services.TransactionsServiceInterface, jwtService jwt.TokenService, logger logger.Logger) *TransactionsHandler {
	return &TransactionsHandler{
		transactionsService: transactionsService,
		jwtService:          jwtService,
	}
}
