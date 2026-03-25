package authentificationhandlers

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/services"
)

type AuthHandler struct {
	userService services.UsersServiceInterface
	jwtService  jwt.TokenService
}

func NewAuthHandler(userService services.UsersServiceInterface, jwtService jwt.TokenService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}