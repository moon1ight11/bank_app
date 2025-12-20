package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"fmt"
	"regexp"
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

// валидация кастомных полей клеймов
func (c *Claims) CustomFieldsValidate() error {
	// проверяем валидность uuid
	if c.UserId == uuid.Nil {
		return fmt.Errorf("user id is empty")
	}

	// проверяем, что имя пользователя не пустое
	if c.UserName == "" {
		return fmt.Errorf("invalid user name")
	}

	// проверяем что фамилия пользователя не пустая
	if c.UserSurname == "" {
		return fmt.Errorf("invalid user surname")
	}

	// проверяем валидность почты
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, c.UserEmail)
	if err != nil {
		return fmt.Errorf("error in matchString: %w", err)
	}
	if !matched {
		return fmt.Errorf("user email not looks like email")
	}

	return nil
}

// создание токена
func (j *Service) GenerateToken(userID uuid.UUID, userName string, userSurname string, userEmail string) (string, error) {
	claims := &Claims{
		UserId:      userID,
		UserName:    userName,
		UserSurname: userSurname,
		UserEmail:   userEmail,
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
