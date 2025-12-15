package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"regexp"
)

type Claims struct {
	UserId      uuid.UUID `json:"user_id"`
	UserName    string    `json:"user_name"`
	UserSurname string    `json:"user_surname"`
	UserEmail   string    `json:"user_email"`
	jwt.RegisteredClaims
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
