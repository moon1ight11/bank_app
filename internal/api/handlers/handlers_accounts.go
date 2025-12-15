package handlers

import (
	"bank_app/internal/jwt"
	"bank_app/internal/services"

	"github.com/gin-gonic/gin"
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

// создание нового счёта
func (a *AccountsHandler) CreateAccount(c *gin.Context) {
	
}