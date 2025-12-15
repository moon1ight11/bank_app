package handlers

import (
	"bank_app/internal/jwt"
	"bank_app/internal/services"
)

type AuthHandler struct {
	userService *services.UsersService
	jwtService jwt.TokenService
}

func NewAuthHandler(userService *services.UsersService, jwtService jwt.TokenService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}
