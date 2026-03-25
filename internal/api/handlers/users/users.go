package usershandlers

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/services"
)

type UsersHandler struct {
	userService services.UsersServiceInterface
	jwtService  jwt.TokenService
}

func NewUsersHandler(userService services.UsersServiceInterface, jwtService jwt.TokenService) *UsersHandler {
	return &UsersHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}
