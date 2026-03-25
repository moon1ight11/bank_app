package usershandlers

import (
	"bank_app/internal/api/helpers"
	"bank_app/internal/api/models"
	"context"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// получение данных пользователя
func (u *UsersHandler) GetUser(c *gin.Context) {
	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		log.Println("Error in EAVCUI", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем пользователя
	user, err := u.userService.UserGet(ctx, userID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// изменение данных пользователя
func (u *UsersHandler) UpdateUser(c *gin.Context) {
	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		log.Println("Error in EAVCUI", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var updatedUser models.UserUpdate

	// устанавливаем user_id
	updatedUser.ID = userID

	// получаем обновленного пользователя с фронта
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// если обновляется имя - чтобы было не пустое
	if updatedUser.Name != nil {
		if strings.TrimSpace(*updatedUser.Name) == "" {
			log.Println("New name is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "New name is empty"})
			return
		}
	}

	// если обновляется пароль - чтобы не был пустым
	if updatedUser.Password != nil {
		if strings.TrimSpace(*updatedUser.Password) == "" {
			log.Println("New pass is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "New pass is empty"})
			return
		}
	}

	// если обновляется почта
	if updatedUser.Email != nil {
		// проверяем, похожа ли новая почта на почту
		pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, err := regexp.MatchString(pattern, *updatedUser.Email)
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
	if updatedUser.PhoneNumber != nil {
		if strings.TrimSpace(*updatedUser.PhoneNumber) == "" {
			log.Println("New phone number is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "New phone number is empty"})
			return
		}
	}

	// если обновляется временная зона
	if updatedUser.Timezone != nil {
		if strings.TrimSpace(*updatedUser.Timezone) == "" {
			log.Println("New timezone is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "New timezone is empty"})
			return
		}
	}

	// обновляем нужные поля
	err = u.userService.UserUpdate(
		ctx,
		updatedUser.Name,
		updatedUser.Surname,
		updatedUser.Password,
		updatedUser.Email,
		updatedUser.PhoneNumber,
		updatedUser.Timezone,
		updatedUser.ID,
	)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// получение структуры обновленного пользователя
	foundUser, err := u.userService.UserGet(ctx, updatedUser.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"UpdatedUser": foundUser})
}

// удаление пользователя
func (u *UsersHandler) DeleteUser(c *gin.Context) {
	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		log.Println("Error in EAVCUI", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// удаляем пользователя
	err = u.userService.UserDelete(ctx, userID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// обнуляем куки
	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"message": "delete is complete"})
}

// создание админа или верификатора
func (u *UsersHandler) CreateAdminOrVerificator(c *gin.Context) {
	// получаем с фронта пользователя
	var user models.UserRegister
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверка заполненности всех обязательных полей
	if strings.TrimSpace(user.Name) == "" ||
		strings.TrimSpace(user.Surname) == "" ||
		strings.TrimSpace(user.Email) == "" ||
		strings.TrimSpace(user.Password) == "" ||
		strings.TrimSpace(user.PhoneNumber) == "" ||
		user.Role == "" {

		log.Println("One or more required fields are empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields (Name, Surname, Email, Password, PhoneNumber, Role) are required"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userId, err := u.userService.AdminAdd(ctx, user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"newUserId": userId})
}

// получение списка юзеров с фильтром на роль
func (u *UsersHandler) GetAllUsers(c *gin.Context) {
	// получаем роль из параметров
	param := c.Query("role")
	role := models.Role(param)

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	users, err := u.userService.UsersByRoleGet(ctx, role)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// изменение роли пользователя
func (u *UsersHandler) ChangeRole(c *gin.Context) {
	// получаем с фронта новую роль
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем id пользователя из параметров
	idStr := c.Param("user_id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error in parse uuid"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// меняем ему роль
	err = u.userService.RoleChange(ctx, userID, role)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// получаем обновленного пользователя
	user, err := u.userService.UserGet(ctx, userID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
