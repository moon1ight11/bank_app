package transactionshandlers

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/services"
)

type TransactionsHandler struct {
	transactionsService services.TransactionsServiceInterface
	jwtService          jwt.TokenService
}

func NewTransactionsHandler(transactionsService services.TransactionsServiceInterface, jwtService jwt.TokenService) *TransactionsHandler {
	return &TransactionsHandler{
		transactionsService: transactionsService,
		jwtService:          jwtService,
	}
}
