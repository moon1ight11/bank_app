package transactionshandlers

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

// пополнение счета
func (t *TransactionsHandler) CreateIncomingTransaction(c *gin.Context) {
	t.metrics.RecordOperation("create_incoming_transaction")

	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		t.metrics.RecordError(string(monitoring.ErrBadRequest), "CreateIncomingTransaction")
		t.logger.Error("Error in CreateIncomingTransaction", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	transactionID, err := t.transactionsService.TransactionIncoming(ctx, transaction)
	if err != nil {
		t.metrics.RecordError(string(monitoring.ErrBusinessLogic), "CreateIncomingTransaction")
		t.logger.Error("Error in CreateIncomingTransaction", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	t.logger.Info("Incoming transaction created successfully",
		"id", transactionID,
		"account_to", transaction.AccountTo,
		"amount", transaction.Amount,
		"currency", transaction.Currency,
	)

	c.JSON(http.StatusOK, gin.H{"transaction_id": transactionID})
}

// списание со счета
func (t *TransactionsHandler) CreateOutcomingTransaction(c *gin.Context) {
	t.metrics.RecordOperation("create_outcoming_transaction")

	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		t.metrics.RecordError(string(monitoring.ErrBadRequest), "CreateOutcomingTransaction")
		t.logger.Error("Error in CreateOutcomingTransaction", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	transactionID, err := t.transactionsService.TransactionOutcoming(ctx, transaction)
	if err != nil {
		t.metrics.RecordError(string(monitoring.ErrBusinessLogic), "CreateOutcomingTransaction")
		t.logger.Error("Error in CreateOutcomingTransaction", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	t.logger.Info("Outcoming transaction created successfully",
		"id", transactionID,
		"account_from", transaction.AccountFrom,
		"amount", transaction.Amount,
		"currency", transaction.Currency,
	)

	c.JSON(http.StatusOK, gin.H{"transaction_id": transactionID})
}

// трансфер
func (t *TransactionsHandler) CreateTransferTransaction(c *gin.Context) {
	t.metrics.RecordOperation("create_transfer_transaction")

	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		t.metrics.RecordError(string(monitoring.ErrExtractUserId), "CreateTransferTransaction")
		t.logger.Error("Error in CreateTransferTransaction", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	userRole, exist := c.Get("UserRole")
	if !exist {
		t.metrics.RecordError(string(monitoring.ErrInternal), "CreateTransferTransaction")
		t.logger.Warn("Error in CreateTransferTransaction", "warn:", "users role not found")
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if userRole != models.RoleBasic {
		t.metrics.RecordError(string(monitoring.ErrForbidden), "CreateTransferTransaction")
		t.logger.Warn("Error in CreateTransferTransaction", "warn:", "try to make transfer by user/verificator")
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin and Verificator cant make transfer"})
		return
	}

	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		t.metrics.RecordError(string(monitoring.ErrBadRequest), "CreateTransferTransaction")
		t.logger.Error("Error in CreateTransferTransaction", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if transaction.UserFrom != userID {
		t.metrics.RecordError(string(monitoring.ErrForbidden), "CreateTransferTransaction")
		t.logger.Warn("Error in CreateTransferTransaction", "warn:", "transaction from foreign user")
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction from foreign user"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	transactionID, err := t.transactionsService.TransactionTransfer(ctx, transaction)
	if err != nil {
		t.metrics.RecordError(string(monitoring.ErrBusinessLogic), "CreateTransferTransaction")
		t.logger.Error("Error in CreateTransferTransaction", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	t.logger.Info("Transfer transaction created successfully",
		"id", transactionID,
		"account_from", transaction.AccountFrom,
		"account_to", transaction.AccountTo,
		"amount", transaction.Amount,
		"currency", transaction.Currency,
	)

	c.JSON(http.StatusOK, gin.H{"transaction_id": transactionID})
}

// получение всех транзакций пользователя
func (t *TransactionsHandler) GetAllUserTransactions(c *gin.Context) {
	t.metrics.RecordOperation("get_all_user_transactions")

	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		t.metrics.RecordError(string(monitoring.ErrExtractUserId), "GetAllUserTransactions")
		t.logger.Error("Error in GetAllUserTransactions", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	transactions, err := t.transactionsService.AllTransactionsGet(ctx, userID)
	if err != nil {
		t.metrics.RecordError(string(monitoring.ErrBusinessLogic), "GetAllUserTransactions")
		t.logger.Error("Error in GetAllUserTransactions", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// получение всех транзакций по счету
func (t *TransactionsHandler) GetAllAccountTransactions(c *gin.Context) {
	t.metrics.RecordOperation("get_all_account_transactions")

	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		t.metrics.RecordError(string(monitoring.ErrExtractUserId), "GetAllAccountTransactions")
		t.logger.Error("Error in GetAllAccountTransactions", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	idStr := c.Param("account_id")
	accountID, err := uuid.Parse(idStr)
	if err != nil {
		t.metrics.RecordError(string(monitoring.ErrParseUUID), "GetAllAccountTransactions")
		t.logger.Error("Error in GetAllAccountTransactions", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	transactions, err := t.transactionsService.AccountTransactionsGet(ctx, userID, accountID)
	if err != nil {
		t.metrics.RecordError(string(monitoring.ErrBusinessLogic), "GetAllAccountTransactions")
		t.logger.Error("Error in GetAllAccountTransactions", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account_transactions": transactions})
}

// получение транзакции по айди
func (t *TransactionsHandler) GetTransactionByID(c *gin.Context) {
	t.metrics.RecordOperation("get_tansaction_by_id")

	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		t.metrics.RecordError(string(monitoring.ErrExtractUserId), "GetTransactionByID")
		t.logger.Error("Error in GetTransactionByID", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	idStr := c.Param("transaction_id")
	transactionID, err := uuid.Parse(idStr)
	if err != nil {
		t.metrics.RecordError(string(monitoring.ErrParseUUID), "GetTransactionByID")
		t.logger.Error("Error in GetTransactionByID", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	transaction, err := t.transactionsService.TransactionByIdGet(ctx, userID, transactionID)
	if err != nil {
		t.metrics.RecordError(string(monitoring.ErrBusinessLogic), "GetTransactionByID")
		t.logger.Error("Error in GetTransactionByID", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}
