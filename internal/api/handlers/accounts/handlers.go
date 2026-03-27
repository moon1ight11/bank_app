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
	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		a.logger.Error("Error in CreateAccount", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newAccount models.AccountCreate

	// получаем валюту счета с фронта
	if err := c.ShouldBindJSON(&newAccount); err != nil {
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
		a.logger.Error("Error in CreateAccount", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"account_id": accountId})
}

// список счетов пользователя
func (a *AccountsHandler) GetAllUserAccounts(c *gin.Context) {
	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
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
		a.logger.Error("Error in GetAllUserAccounts", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

// конкретный счет пользователя
func (a *AccountsHandler) GetAccountById(c *gin.Context) {
	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		a.logger.Error("Error in GetAccountById", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем id счета из параметров
	idStr := c.Param("account_id")
	accountID, err := uuid.Parse(idStr)
	if err != nil {
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
		a.logger.Error("Error in GetAccountById", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}

// удаление счета
func (a *AccountsHandler) DeleteAccount(c *gin.Context) {
	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		a.logger.Error("Error in DeleteAccount", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем id счета из параметров
	idStr := c.Param("account_id")
	accountID, err := uuid.Parse(idStr)
	if err != nil {
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
		a.logger.Error("Error in DeleteAccount", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successful delete"})
}
