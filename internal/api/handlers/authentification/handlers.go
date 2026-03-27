package authentificationhandlers

import (
	"bank_app/internal/api/helpers"
	"bank_app/internal/api/models"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// регистрация
func (u *AuthHandler) SignUp(c *gin.Context) {
	// получаем данные пользователя с фронта
	var user models.UserRegister
	if err := c.ShouldBindJSON(&user); err != nil {
		u.logger.Error("Error in SignUp", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверяем, что все нужные поля заполнены
	if strings.TrimSpace(user.Name) == "" ||
		strings.TrimSpace(user.Surname) == "" ||
		strings.TrimSpace(user.Email) == "" ||
		strings.TrimSpace(user.Password) == "" ||
		strings.TrimSpace(user.PhoneNumber) == "" {

		u.logger.Error("Error in SignUp", "error:", "One or more required fields are empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields (Name, Surname, Email, Password, PhoneNumber) are required"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// добавляем пользователя в БД
	userID, err := u.userService.UserAdd(ctx, user)
	if err != nil {
		u.logger.Error("Error in SignUp", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// генерируем токен для нового пользователя
	token, err := u.jwtService.GenerateToken(userID, user.Name, user.Surname, user.Email, user.Role)
	if err != nil {
		u.logger.Error("Error in SignUp", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// устанавливаем куки
	c.SetCookie("cookie", token, 3600, "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{"user_id": userID})
}

// авторизация
func (u *AuthHandler) SignIn(c *gin.Context) {
	// получаем данные пользователя с фронта
	var user models.UserAutorization
	if err := c.ShouldBindJSON(&user); err != nil {
		u.logger.Error("Error in SignIn", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверка обязательных полей
	if strings.TrimSpace(user.Password) == "" ||
		strings.TrimSpace(user.Email) == "" {
		u.logger.Error("Error in SignIn", "error:", "One or more required fields are empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password, Email and PhoneNumber are required"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// проверка пользователя
	foundUser, err := u.userService.UserVerification(ctx, user)
	if err != nil {
		u.logger.Error("Error in SignIn", "error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// генерируем токен для найденного пользователя
	token, err := u.jwtService.GenerateToken(foundUser.Id, foundUser.Name, foundUser.Surname, foundUser.Email, foundUser.Role)
	if err != nil {
		u.logger.Error("Error in SignIn", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// устанавливаем куки
	c.SetCookie("cookie", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"user_id": foundUser.Id})
}

// выход из профиля
func (u *AuthHandler) SignOut(c *gin.Context) {
	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		u.logger.Error("Error in SignOut", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	u.logger.Info("User signed out", "user_id:", userID)

	c.JSON(http.StatusOK, gin.H{"message": "successful sign out"})
}
