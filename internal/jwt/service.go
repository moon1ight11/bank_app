package jwt

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	secret     []byte
	expiration time.Duration
}

func NewJWTService(secret string, expiration time.Duration) TokenService {
	return &Service{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

// создание токена
func (j *Service) GenerateToken(userID uuid.UUID, userName string, userSurname string, userEmail string) (string, error) {
	claims := &Claims{
		UserId:    userID,
		UserName:  userName,
		UserSurname: userSurname,
		UserEmail: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
		},
	}

	// валидация кастомных полей
	if err := claims.CustomFieldsValidate(); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// декодировка токена
func (j *Service) ParseToken(tokenString string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid method")
		}
		return j.secret, nil
	})
}
