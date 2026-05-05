package usershandlers

import (
	"bank_app/internal/api/helpers"
	"bank_app/internal/api/models"
	"bank_app/internal/monitoring"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// получение данных пользователя
func (u *UsersHandler) GetUser(c *gin.Context) {
	// записываем операцию в метрики
	u.metrics.RecordOperation("get_user")

	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrExtractUserId), "GetUser")
		u.logger.Error("Error in GetUser", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем пользователя
	user, err := u.userService.UserGet(ctx, userID)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrBusinessLogic), "GetUser")
		u.logger.Error("Error in GetUser", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// изменение данных пользователя
func (u *UsersHandler) UpdateUser(c *gin.Context) {
	// записываем операцию в метрики
	u.metrics.RecordOperation("update_user")

	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrExtractUserId), "UpdateUser")
		u.logger.Error("Error in UpdateUser", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
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
		u.metrics.RecordError(string(monitoring.ErrBadRequest), "UpdateUser")
		u.logger.Error("Error in UpdateUser", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// если обновляется имя - чтобы было не пустое
	if updatedUser.Name != nil {
		if strings.TrimSpace(*updatedUser.Name) == "" {
			u.metrics.RecordError(string(monitoring.ErrInvalidInput), "UpdateUser")
			u.logger.Error("Error in UpdateUser", "error:", "new name is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "New name is empty"})
			return
		}
	}

	// если обновляется пароль - чтобы не был пустым
	if updatedUser.Password != nil {
		if strings.TrimSpace(*updatedUser.Password) == "" {
			u.metrics.RecordError(string(monitoring.ErrInvalidInput), "UpdateUser")
			u.logger.Error("Error in UpdateUser", "error:", "new pass is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "New pass is empty"})
			return
		}
	}

	// если обновляется почта
	if updatedUser.Email != nil {
		if strings.TrimSpace(*updatedUser.Email) == "" {
			u.metrics.RecordError(string(monitoring.ErrInvalidInput), "UpdateUser")
			u.logger.Error("Error in UpdateUser", "error:", "new email is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "New email is empty"})
			return
		}
	}

	// если обновляется телефон
	if updatedUser.PhoneNumber != nil {
		if strings.TrimSpace(*updatedUser.PhoneNumber) == "" {
			u.metrics.RecordError(string(monitoring.ErrInvalidInput), "UpdateUser")
			u.logger.Error("Error in UpdateUser", "error:", "new phone number is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "New phone number is empty"})
			return
		}
	}

	// если обновляется временная зона
	if updatedUser.Timezone != nil {
		if strings.TrimSpace(*updatedUser.Timezone) == "" {
			u.metrics.RecordError(string(monitoring.ErrInvalidInput), "UpdateUser")
			u.logger.Error("Error in UpdateUser", "error:", "new timezone is empty")
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
		u.metrics.RecordError(string(monitoring.ErrBusinessLogic), "UpdateUser")
		u.logger.Error("Error in UpdateUser", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// получение структуры обновленного пользователя
	foundUser, err := u.userService.UserGet(ctx, updatedUser.ID)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrBusinessLogic), "UpdateUser")
		u.logger.Error("Error in UpdateUser", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	u.logger.Info("User updated successfully", "userId", userID)

	c.JSON(http.StatusOK, gin.H{"updated_user": foundUser})
}

// удаление пользователя
func (u *UsersHandler) DeleteUser(c *gin.Context) {
	// записываем операцию в метрики
	u.metrics.RecordOperation("delete_user")

	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrExtractUserId), "DeleteUser")
		u.logger.Error("Error in DeleteUser", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// удаляем пользователя
	err = u.userService.UserDelete(ctx, userID)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrBusinessLogic), "DeleteUser")
		u.logger.Error("Error in DeleteUser", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// обнуляем куки
	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	u.logger.Info("User deleted successfully", "userId", userID)

	c.JSON(http.StatusOK, gin.H{"message": "successful delete"})
}

// создание админа или верификатора
func (u *UsersHandler) CreateAdminOrVerificator(c *gin.Context) {
	// записываем операцию в метрики
	u.metrics.RecordOperation("create_admin/verif")

	// получаем с фронта пользователя
	var user models.UserRegister
	if err := c.ShouldBindJSON(&user); err != nil {
		u.metrics.RecordError(string(monitoring.ErrBadRequest), "CreateAdminOrVerificator")
		u.logger.Error("Error in CreateAdminOrVerificator", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// проверка заполненности всех обязательных полей
	if strings.TrimSpace(user.Name) == "" ||
		strings.TrimSpace(user.Surname) == "" ||
		strings.TrimSpace(user.Email) == "" ||
		strings.TrimSpace(user.Password) == "" ||
		strings.TrimSpace(user.PhoneNumber) == "" ||
		user.Role == "" {

		u.metrics.RecordError(string(monitoring.ErrInvalidInput), "CreateAdminOrVerificator")
		u.logger.Error("Error in CreateAdminOrVerificator", "error:", "one or more required fields are empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields (Name, Surname, Email, Password, PhoneNumber, Role) are required"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userId, err := u.userService.AdminAdd(ctx, user)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrBusinessLogic), "CreateAdminOrVerificator")
		u.logger.Error("Error in CreateAdminOrVerificator", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	u.logger.Info("User created successfully", "userId", userId, "role", user.Role)

	c.JSON(http.StatusCreated, gin.H{"admin_or_verificator_id": userId})
}

// получение списка юзеров с фильтром на роль
func (u *UsersHandler) GetAllUsers(c *gin.Context) {
	// записываем операцию в метрики
	u.metrics.RecordOperation("get_all_users")

	// получаем роль из параметров
	param := c.Query("role")
	role := models.Role(param)

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	users, err := u.userService.UsersByRoleGet(ctx, role)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrBusinessLogic), "GetAllUsers")
		u.logger.Error("Error in GetAllUsers", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// изменение роли пользователя
func (u *UsersHandler) ChangeRole(c *gin.Context) {
	// записываем операцию в метрики
	u.metrics.RecordOperation("change_role")

	// получаем с фронта новую роль
	var role models.ChangeRoleRequest
	if err := c.ShouldBindJSON(&role); err != nil {
		u.metrics.RecordError(string(monitoring.ErrBadRequest), "ChangeRole")
		u.logger.Error("Error in ChangeRole", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// получаем id пользователя из параметров
	idStr := c.Param("user_id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrParseUUID), "ChangeRole")
		u.logger.Error("Error in ChangeRole", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// меняем ему роль
	err = u.userService.RoleChange(ctx, userID, role.Role)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrBusinessLogic), "ChangeRole")
		u.logger.Error("Error in ChangeRole", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// получаем обновленного пользователя
	user, err := u.userService.UserGet(ctx, userID)
	if err != nil {
		u.metrics.RecordError(string(monitoring.ErrBusinessLogic), "ChangeRole")
		u.logger.Error("Error in ChangeRole", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	u.logger.Info("User role changed successfully",
		"userId", user.Id,
		"new_role", user.Role,
		"old_role", role,
	)

	c.JSON(http.StatusOK, gin.H{"user": user})
}
