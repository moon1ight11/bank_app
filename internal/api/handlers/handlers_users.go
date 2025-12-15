package handlers

import (
	"bank_app/internal/jwt"
	"bank_app/internal/services"
)

type UsersHandler struct {
	userService *services.UsersService
	jwtService jwt.TokenService
}

func NewUsersHandler(userService *services.UsersService, jwtService jwt.TokenService) *UsersHandler {
	return &UsersHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}