package accountshandlers

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/services"
	"bank_app/pkg/logger"
)

type AccountsHandler struct {
	accountsService services.AccountsServiceInterface
	jwtService      jwt.TokenService
	logger          logger.Logger
}

func NewAccountsHandler(accountsService services.AccountsServiceInterface, jwtService jwt.TokenService, logger logger.Logger) *AccountsHandler {
	return &AccountsHandler{
		accountsService: accountsService,
		jwtService:      jwtService,
		logger:    logger,
	}
}
