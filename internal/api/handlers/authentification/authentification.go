package authentificationhandlers

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/services"
	"bank_app/pkg/logger"
)

type AuthHandler struct {
	userService services.UsersServiceInterface
	jwtService  jwt.TokenService
	logger logger.Logger
}

func NewAuthHandler(userService services.UsersServiceInterface, jwtService jwt.TokenService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
		logger: logger,
	}
}