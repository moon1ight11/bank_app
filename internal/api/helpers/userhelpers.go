package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ExtractAndValidateContextUserId(c *gin.Context) (uuid.UUID, error) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		return uuid.Nil, fmt.Errorf("error in ExtractAndValidateUserId: UserId not found in context")
	}

	// приводим значение к uuid
	UserID, ok := userIDValue.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("error in ExtractAndValidateUserId: Invalid user ID type")
	}

	return UserID, nil
}
