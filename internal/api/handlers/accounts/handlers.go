package accountshandlers

import (
	"bank_app/internal/api/helpers"
	"bank_app/internal/api/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// создание нового счёта
func (a *AccountsHandler) CreateAccount(c *gin.Context) {
	// записываем операцию в метрики
	a.metrics.RecordOperation("create_account")

	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		a.metrics.RecordError(err.Error(), "CreateAccount")
		a.logger.Error("Error in CreateAccount", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newAccount models.AccountCreate

	// получаем валюту счета с фронта
	if err := c.ShouldBindJSON(&newAccount); err != nil {
		a.metrics.RecordError(err.Error(), "CreateAccount")
		a.logger.Error("Error in CreateAccount", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// устанавливаем владельцем счета пользователя
	newAccount.UserID = userID

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// создаем счет
	accountId, err := a.accountsService.AccountAdd(ctx, newAccount)
	if err != nil {
		a.metrics.RecordError(err.Error(), "CreateAccount")
		a.logger.Error("Error in CreateAccount", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	a.logger.Info("Account created successfully", "userId", userID, "accountId", accountId)

	c.JSON(http.StatusCreated, gin.H{"account_id": accountId})
}

// список счетов пользователя
func (a *AccountsHandler) GetAllUserAccounts(c *gin.Context) {
	// записываем операцию в метрики
	a.metrics.RecordOperation("get_all_users_accounts")

	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		a.metrics.RecordError(err.Error(), "GetAllUserAccounts")
		a.logger.Error("Error in GetAllUserAccounts", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем список счетов конкретного пользователя
	accounts, err := a.accountsService.AllAccountsGet(ctx, userID)
	if err != nil {
		a.metrics.RecordError(err.Error(), "GetAllUserAccounts")
		a.logger.Error("Error in GetAllUserAccounts", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

// конкретный счет пользователя
func (a *AccountsHandler) GetAccountById(c *gin.Context) {
	// записываем операцию в метрики
	a.metrics.RecordOperation("get_account_by_id")

	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		a.metrics.RecordError(err.Error(), "GetAccountById")
		a.logger.Error("Error in GetAccountById", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем id счета из параметров
	idStr := c.Param("account_id")
	accountID, err := uuid.Parse(idStr)
	if err != nil {
		a.metrics.RecordError(err.Error(), "GetAccountById")
		a.logger.Error("Error in GetAccountById", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error in parse uuid"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем счет из БД
	account, err := a.accountsService.AccountGet(ctx, userID, accountID)
	if err != nil {
		a.metrics.RecordError(err.Error(), "GetAccountById")
		a.logger.Error("Error in GetAccountById", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}

// удаление счета
func (a *AccountsHandler) DeleteAccount(c *gin.Context) {
	// записываем операцию в метрики
	a.metrics.RecordOperation("delete_account")

	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		a.metrics.RecordError(err.Error(), "DeleteAccount")
		a.logger.Error("Error in DeleteAccount", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем id счета из параметров
	idStr := c.Param("account_id")
	accountID, err := uuid.Parse(idStr)
	if err != nil {
		a.metrics.RecordError(err.Error(), "DeleteAccount")
		a.logger.Error("Error in DeleteAccount", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error in parse uuid"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// удаляем счет
	err = a.accountsService.AccountDelete(ctx, userID, accountID)
	if err != nil {
		a.metrics.RecordError(err.Error(), "DeleteAccount")
		a.logger.Error("Error in DeleteAccount", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	a.logger.Info("Account deleted successfully", "userId", userID, "accountId", accountID)

	c.JSON(http.StatusOK, gin.H{"message": "successful delete"})
}
