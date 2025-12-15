package handlers

import (
	"bank_app/internal/jwt"
	"bank_app/internal/services"
)

type AccountsHandler struct {
	accountsService *services.AccountsService
	jwtService jwt.TokenService
}

func NewAccountsHandler(accountsService *services.AccountsService, jwtService jwt.TokenService) *AccountsHandler {
	return &AccountsHandler{
		accountsService: accountsService,
		jwtService:  jwtService,
	}
}