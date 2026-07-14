package authentificationhandlers

import (
	"bank_app/internal/api/helpers"
	"bank_app/internal/api/models"
	"bank_app/internal/monitoring"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// регистрация
func (u *AuthHandler) SignUp(c *gin.Context) {
	u.metrics.RecordOperation("sign_up")

	var user models.UserRegister
	if err := c.ShouldBindJSON(&user); err != nil {
		u.metrics.RecordError(string(monitoring.ErrBadRequest), "SignUp")
		u.logger.Error("Error in SignUp", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	user.Role = models.RoleBasic

	if strings.TrimSpace(user.Name) == "" ||
		strings.TrimSpace(user.Surname) == "" ||
		strings.TrimSpace(user.Email) == "" ||
		strings.TrimSpace(user.Password) == "" ||
		strings.TrimSpace(user.PhoneNumber) == "" {

		u.metrics.RecordError(string(monitoring.ErrInvalidInput), "SignUp")
		u.logger.Error("Error in SignUp", "error:", "One or more required fields are empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields (Name, Surname, Email, Password, PhoneNumber) are required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID, err := u.userService.UserAdd(ctx, user)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrBusinessLogic), "SignUp")
		u.logger.Error("Error in SignUp", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	token, err := u.jwtService.GenerateToken(userID, user.Name, user.Surname, user.Email, user.Role)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrInternal), "SignUp")
		u.logger.Error("Error in SignUp", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.SetCookie("cookie", token, 3600, "/", "", false, true)

	u.logger.Info("User signed up successfully", "userId", userID)

	c.JSON(http.StatusCreated, gin.H{"user_id": userID})
}

// авторизация
func (u *AuthHandler) SignIn(c *gin.Context) {
	u.metrics.RecordOperation("sign_in")

	var user models.UserAutorization
	if err := c.ShouldBindJSON(&user); err != nil {
		u.metrics.RecordError(string(monitoring.ErrBadRequest), "SignIn")
		u.logger.Error("Error in SignIn", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if strings.TrimSpace(user.Password) == "" || strings.TrimSpace(user.Email) == "" {
		u.metrics.RecordError(string(monitoring.ErrInvalidInput), "SignIn")
		u.logger.Error("Error in SignIn", "error:", "One or more required fields are empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password, Email and PhoneNumber are required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	foundUser, err := u.userService.UserVerification(ctx, user)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrForbidden), "SignIn")
		u.logger.Error("Error in SignIn", "error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	token, err := u.jwtService.GenerateToken(foundUser.Id, foundUser.Name, foundUser.Surname, foundUser.Email, foundUser.Role)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrInternal), "SignIn")
		u.logger.Error("Error in SignIn", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.SetCookie("cookie", token, 3600, "/", "", false, true)

	u.logger.Info("User signed in successfully", "userId", foundUser.Id)

	c.JSON(http.StatusOK, gin.H{"user_id": foundUser.Id})
}

// логаут
func (u *AuthHandler) SignOut(c *gin.Context) {
	u.metrics.RecordOperation("sign_out")

	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrExtractUserId), "SignOut")
		u.logger.Error("Error in SignOut", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	u.logger.Info("User signed out successfully", "userId:", userID)

	c.JSON(http.StatusOK, gin.H{"message": "successful sign out"})
}
