package accountshandlers

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/monitoring"
	"bank_app/internal/services"
	"bank_app/pkg/logger"
)

type AccountsHandler struct {
	accountsService services.AccountsServiceInterface
	jwtService      jwt.TokenService
	logger          logger.Logger
	metrics         *monitoring.Metrics
}

func NewAccountsHandler(
	accountsService services.AccountsServiceInterface,
	jwtService jwt.TokenService,
	logger logger.Logger,
	metrics *monitoring.Metrics,
) *AccountsHandler {
	return &AccountsHandler{
		accountsService: accountsService,
		jwtService:      jwtService,
		logger:          logger,
		metrics:         metrics,
	}
}
