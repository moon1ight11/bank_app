package accountshandlers

import (
	"bank_app/internal/api/helpers"
	"bank_app/internal/api/models"
	"bank_app/internal/monitoring"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// создание нового счёта
func (a *AccountsHandler) CreateAccount(c *gin.Context) {
	a.metrics.RecordOperation("create_account")

	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		a.metrics.RecordError(string(monitoring.ErrExtractUserId), "CreateAccount")
		a.logger.Error("Error in CreateAccount", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	var newAccount models.AccountCreate
	if err := c.ShouldBindJSON(&newAccount); err != nil {
		a.metrics.RecordError(string(monitoring.ErrBadRequest), "CreateAccount")
		a.logger.Error("Error in CreateAccount", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ba request"})
		return
	}

	newAccount.UserID = userID

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	accountId, err := a.accountsService.AccountAdd(ctx, newAccount)
	if err != nil {
		a.metrics.RecordError(string(monitoring.ErrBusinessLogic), "CreateAccount")
		a.logger.Error("Error in CreateAccount", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	a.logger.Info("Account created successfully", "userId", userID, "accountId", accountId)

	c.JSON(http.StatusCreated, gin.H{"account_id": accountId})
}

// получение списка счетов пользователя
func (a *AccountsHandler) GetAllUserAccounts(c *gin.Context) {
	a.metrics.RecordOperation("get_all_users_accounts")

	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		a.metrics.RecordError(string(monitoring.ErrExtractUserId), "GetAllUserAccounts")
		a.logger.Error("Error in GetAllUserAccounts", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	accounts, err := a.accountsService.AllAccountsGet(ctx, userID)
	if err != nil {
		a.metrics.RecordError(string(monitoring.ErrBusinessLogic), "GetAllUserAccounts")
		a.logger.Error("Error in GetAllUserAccounts", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

// получение одного счета пользователя
func (a *AccountsHandler) GetAccountById(c *gin.Context) {
	a.metrics.RecordOperation("get_account_by_id")

	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		a.metrics.RecordError(string(monitoring.ErrExtractUserId), "GetAccountById")
		a.logger.Error("Error in GetAccountById", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	idStr := c.Param("account_id")
	accountID, err := uuid.Parse(idStr)
	if err != nil {
		a.metrics.RecordError(string(monitoring.ErrParseUUID), "GetAccountById")
		a.logger.Error("Error in GetAccountById", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	account, err := a.accountsService.AccountGet(ctx, userID, accountID)
	if err != nil {
		a.metrics.RecordError(string(monitoring.ErrBusinessLogic), "GetAccountById")
		a.logger.Error("Error in GetAccountById", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}

// удаление счета
func (a *AccountsHandler) DeleteAccount(c *gin.Context) {
	a.metrics.RecordOperation("delete_account")

	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		a.metrics.RecordError(string(monitoring.ErrExtractUserId), "DeleteAccount")
		a.logger.Error("Error in DeleteAccount", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	idStr := c.Param("account_id")
	accountID, err := uuid.Parse(idStr)
	if err != nil {
		a.metrics.RecordError(string(monitoring.ErrParseUUID), "DeleteAccount")
		a.logger.Error("Error in DeleteAccount", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = a.accountsService.AccountDelete(ctx, userID, accountID)
	if err != nil {
		a.metrics.RecordError(string(monitoring.ErrBusinessLogic), "DeleteAccount")
		a.logger.Error("Error in DeleteAccount", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	a.logger.Info("Account deleted successfully", "userId", userID, "accountId", accountID)

	c.JSON(http.StatusOK, gin.H{"message": "successful delete"})
}
