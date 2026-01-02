package jwt

import (
	"bank_app/internal/api/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserId      *uuid.UUID  `json:"user_id"`
	UserName    string      `json:"name"`
	UserSurname string      `json:"surname"`
	UserEmail   string      `json:"email"`
	Role        models.Role `json:"role"`
	jwt.RegisteredClaims
}
