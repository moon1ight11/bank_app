package transactionshandlers

import (
	"bank_app/internal/api/helpers"
	"bank_app/internal/api/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// получение всех транзакций пользователя
func (t *TransactionsHandler) GetAllUserTransactions(c *gin.Context) {
	// записываем операцию в метрики
	t.metrics.RecordOperation("get_all_user_transactions")

	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		t.metrics.RecordError(err.Error(), "GetAllUserTransactions")
		t.logger.Error("Error in GetAllUserTransactions", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем список транзакций
	transactions, err := t.transactionsService.AllTransactionsGet(ctx, userID)
	if err != nil {
		t.metrics.RecordError(err.Error(), "GetAllUserTransactions")
		t.logger.Error("Error in GetAllUserTransactions", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// получение всех транзакций конкретного счета
func (t *TransactionsHandler) GetAllAccountTransactions(c *gin.Context) {
	// записываем операцию в метрики
	t.metrics.RecordOperation("get_all_account_transactions")

	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		t.metrics.RecordError(err.Error(), "GetAllAccountTransactions")
		t.logger.Error("Error in GetAllAccountTransactions", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем id счета из параметров
	idStr := c.Param("account_id")
	accountID, err := uuid.Parse(idStr)
	if err != nil {
		t.metrics.RecordError(err.Error(), "GetAllAccountTransactions")
		t.logger.Error("Error in GetAllAccountTransactions", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error in parse uuid"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем список транзакций
	transactions, err := t.transactionsService.AccountTransactionsGet(ctx, userID, accountID)
	if err != nil {
		t.metrics.RecordError(err.Error(), "GetAllAccountTransactions")
		t.logger.Error("Error in GetAllAccountTransactions", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account_transactions": transactions})
}

// информация о конкретной транзакции
func (t *TransactionsHandler) GetTransactionByID(c *gin.Context) {
	// записываем операцию в метрики
	t.metrics.RecordOperation("get_tansaction_by_id")

	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		t.metrics.RecordError(err.Error(), "GetTransactionByID")
		t.logger.Error("Error in GetTransactionByID", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем id транзакции из параметров
	idStr := c.Param("transaction_id")
	transactionID, err := uuid.Parse(idStr)
	if err != nil {
		t.metrics.RecordError(err.Error(), "GetTransactionByID")
		t.logger.Error("Error in GetTransactionByID", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error in parse uuid"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем транзакцию
	transaction, err := t.transactionsService.TransactionByIdGet(ctx, userID, transactionID)
	if err != nil {
		t.metrics.RecordError(err.Error(), "GetTransactionByID")
		t.logger.Error("Error in GetTransactionByID", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

// пополнение счета
func (t *TransactionsHandler) CreateIncomingTransaction(c *gin.Context) {
	// записываем операцию в метрики
	t.metrics.RecordOperation("create_incoming_transaction")

	// получаем с фронта тело транзакции
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		t.metrics.RecordError(err.Error(), "CreateIncomingTransaction")
		t.logger.Error("Error in CreateIncomingTransaction", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	transactionID, err := t.transactionsService.TransactionIncoming(ctx, transaction)
	if err != nil {
		t.metrics.RecordError(err.Error(), "CreateIncomingTransaction")
		t.logger.Error("Error in CreateIncomingTransaction", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	// записываем операцию в метрики
	t.metrics.RecordOperation("create_outcoming_transaction")

	// получаем с фронта тело транзакции
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		t.metrics.RecordError(err.Error(), "CreateOutcomingTransaction")
		t.logger.Error("Error in CreateOutcomingTransaction", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	transactionID, err := t.transactionsService.TransactionOutcoming(ctx, transaction)
	if err != nil {
		t.metrics.RecordError(err.Error(), "CreateOutcomingTransaction")
		t.logger.Error("Error in CreateOutcomingTransaction", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	// записываем операцию в метрики
	t.metrics.RecordOperation("create_transfer_transaction")

	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		t.metrics.RecordError(err.Error(), "CreateTransferTransaction")
		t.logger.Error("Error in CreateTransferTransaction", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем роль пользователя из контекста
	userRole, exist := c.Get("UserRole")
	if !exist {
		t.metrics.RecordError("Error in CreateTransferTransaction: users role not found", "CreateTransferTransaction")
		t.logger.Warn("Error in CreateTransferTransaction", "warn:", "users role not found")
		c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
		return
	}

	// запрещаем делать транзакции админам и верификаторам
	if userRole != models.RoleBasic {
		t.metrics.RecordError("Error in CreateTransferTransaction: try to make transfer by user/verificator", "CreateTransferTransaction")
		t.logger.Warn("Error in CreateTransferTransaction", "warn:", "try to make transfer by user/verificator")
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin and Verificator cant make transfer"})
		return
	}

	// получаем с фронта тело транзакции
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		t.metrics.RecordError(err.Error(), "CreateTransferTransaction")
		t.logger.Error("Error in CreateTransferTransaction", "error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// запрещаем делать перевод не от себя
	if transaction.UserFrom != userID {
		t.metrics.RecordError("Error in CreateTransferTransaction: transaction from foreign user", "CreateTransferTransaction")
		t.logger.Warn("Error in CreateTransferTransaction", "warn:", "transaction from foreign user")
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction from foreign user"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	transactionID, err := t.transactionsService.TransactionTransfer(ctx, transaction)
	if err != nil {
		t.metrics.RecordError(err.Error(), "CreateTransferTransaction")
		t.logger.Error("Error in CreateTransferTransaction", "error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
