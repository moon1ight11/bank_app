package handlers

import (
	"bank_app/internal/jwt"
	"bank_app/internal/services"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UsersHandler struct {
	userService *services.UsersService
	jwtService  jwt.TokenService
}

func NewUsersHandler(userService *services.UsersService, jwtService jwt.TokenService) *UsersHandler {
	return &UsersHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

// получение данных пользователя
func (u *UsersHandler) GetUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	UserId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	user, err := u.userService.UserGet(UserId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"user": user})
}

// изменение данных пользователя
func (u *UsersHandler) UpdateUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	UserId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var UpdatedUser struct {
		UserID       uuid.UUID `json:"user_id"`
		UserName     *string   `json:"name"`
		UserSurname  *string   `json:"surname"`
		UserPassword *string   `json:"password"`
		UserNumber   *string   `json:"phone_number"`
		UserEmail    *string   `json:"email"`
		Timezone     *string   `json:"timezone"`
	}

	// устанавливаем user_id
	UpdatedUser.UserID = UserId

	// получаем обновленного пользователя с фронта
	if err := c.ShouldBindJSON(&UpdatedUser); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// если обновляется имя - чтобы было не пустое
	if UpdatedUser.UserName != nil {
		if strings.TrimSpace(*UpdatedUser.UserName) == "" {
			log.Println("New name is empty")
			c.JSON((http.StatusBadRequest), gin.H{"error": "New name is empty"})
			return
		}
	}

	// если обновляется пароль - чтобы не был пустым
	if UpdatedUser.UserPassword != nil {
		if strings.TrimSpace(*UpdatedUser.UserPassword) == "" {
			log.Println("New pass is empty")
			c.JSON((http.StatusBadRequest), gin.H{"error": "New pass is empty"})
			return
		}
	}

	// если обновляется почта
	if UpdatedUser.UserEmail != nil {
		// проверяем, похожа ли новая почта на почту
		pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, err := regexp.MatchString(pattern, *UpdatedUser.UserEmail)
		if err != nil {
			log.Println("Error in MatchString", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// если нет - отклоняем
		if !matched {
			log.Println("New email not looks like email")
			c.JSON(http.StatusBadRequest, gin.H{"error": "New email not looks like email"})
			return
		}
	}

	// если обновляется телефон
	if UpdatedUser.UserNumber != nil {
		if strings.TrimSpace(*UpdatedUser.UserNumber) == "" {
			log.Println("New phone number is empty")
			c.JSON((http.StatusBadRequest), gin.H{"error": "New phone number is empty"})
			return
		}
	}

	// если обновляется временная зона
	if UpdatedUser.Timezone != nil {
		if strings.TrimSpace(*UpdatedUser.Timezone) == "" {
			log.Println("New timezone is empty")
			c.JSON((http.StatusBadRequest), gin.H{"error": "New timezone is empty"})
			return
		}
	}

	// проверка уникальности новых почты и телефона
	if UpdatedUser.UserEmail != nil || UpdatedUser.UserNumber != nil {
		var userNumber, userEmail string 

		if UpdatedUser.UserNumber != nil {
			userNumber = *UpdatedUser.UserNumber
		}

		if UpdatedUser.UserEmail != nil {
			userEmail = *UpdatedUser.UserEmail
		}

		userCheck, err := u.userService.UserCheck(userNumber, userEmail)
		if err != nil {
			log.Println(err)
			c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
			return
		}

		if userCheck {
			log.Println("Phone number or email already exist")
			c.JSON((http.StatusUnauthorized), gin.H{"error": "Phone number or email already exist"})
			return
		}
	}

	// обновляем нужные поля
	err := u.userService.UserUpdate(
		UpdatedUser.UserName,
		UpdatedUser.UserSurname,
		UpdatedUser.UserPassword,
		UpdatedUser.UserEmail,
		UpdatedUser.UserNumber,
		UpdatedUser.Timezone,
		UpdatedUser.UserID,
	)

	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"UpdatedUser": UpdatedUser})
}

// удаление пользователя
func (u *UsersHandler) DeleteUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	UserId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// находим пользователя по id
	_, err := u.userService.UserGet(UserId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// удаляем пользователя
	err = u.userService.UserDelete(UserId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete is complete"})
}
