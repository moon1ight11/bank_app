package jwt

import (
	"bank_app/internal/storage/repos/users"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateToken(userId uuid.UUID, userName string, userSurname string, userEmail string, userRole users.Role) (string, error)
	ParseToken(token string, claims *Claims) (*jwt.Token, error)
}