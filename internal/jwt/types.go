package jwt

import (
	"bank_app/internal/storage/repos/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserId      uuid.UUID  `json:"user_id"`
	UserName    string     `json:"name"`
	UserSurname string     `json:"surname"`
	UserEmail   string     `json:"email"`
	Roles       users.Role `json:"role"`
	jwt.RegisteredClaims
}
