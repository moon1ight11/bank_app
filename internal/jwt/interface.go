package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateToken(userId uuid.UUID, userName string, userSurname string, userEmail string) (string, error)
	ParseToken(token string, claims *Claims) (*jwt.Token, error)
}