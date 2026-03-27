package usershandlers

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/services"
	"bank_app/pkg/logger"
)

type UsersHandler struct {
	userService  services.UsersServiceInterface
	jwtService   jwt.TokenService
	logger logger.Logger
}

func NewUsersHandler(userService services.UsersServiceInterface, jwtService jwt.TokenService, logger logger.Logger) *UsersHandler {
	return &UsersHandler{
		userService:  userService,
		jwtService:   jwtService,
		logger: logger,
	}
}
