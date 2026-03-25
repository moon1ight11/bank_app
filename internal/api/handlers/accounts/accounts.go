package accountshandlers

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/services"
)

type AccountsHandler struct {
	accountsService services.AccountsServiceInterface
	jwtService      jwt.TokenService
}

func NewAccountsHandler(accountsService services.AccountsServiceInterface, jwtService jwt.TokenService) *AccountsHandler {
	return &AccountsHandler{
		accountsService: accountsService,
		jwtService:      jwtService,
	}
}
