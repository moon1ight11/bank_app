package transactionshandlers

import (
	"bank_app/internal/api/helpers"
	"bank_app/internal/api/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// получение всех транзакций пользователя
func (t *TransactionsHandler) GetAllUserTransactions(c *gin.Context) {
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

	// получаем список транзакций
	transactions, err := t.transactionsService.AllTransactionsGet(ctx, userID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// получение всех транзакций конкретного счета
func (t *TransactionsHandler) GetAllAccountTransactions(c *gin.Context) {
	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		log.Println("Error in EAVCUI", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем id счета из параметров
	idStr := c.Param("account_id")
	accountID, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error in parse uuid"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем список транзакций
	transactions, err := t.transactionsService.AccountTransactionsGet(ctx, userID, accountID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// информация о конкретной транзакции
func (t *TransactionsHandler) GetTransactionByID(c *gin.Context) {
	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		log.Println("Error in EAVCUI", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем id транзакции из параметров
	idStr := c.Param("transaction_id")
	transactionID, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error in parse uuid"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем транзакцию
	transaction, err := t.transactionsService.TransactionByIdGet(ctx, userID, transactionID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

// пополнение счета
func (t *TransactionsHandler) CreateIncomingTransaction(c *gin.Context) {
	// получаем с фронта тело транзакции
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	transactionID, err := t.transactionsService.TransactionIncoming(ctx, transaction)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactionID": transactionID})
}

// списание со счета
func (t *TransactionsHandler) CreateOutcomingTransaction(c *gin.Context) {
	// получаем с фронта тело транзакции
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	transactionID, err := t.transactionsService.TransactionOutcoming(ctx, transaction)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactionID": transactionID})
}

// трансфер
func (t *TransactionsHandler) CreateTransferTransaction(c *gin.Context) {
	// получаем userID из контекста
	userID, err := helpers.ExtractAndValidateContextUserId(c)
	if err != nil {
		log.Println("Error in EAVCUI", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем роль пользователя из контекста
	userRole, exist := c.Get("UserRole")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
		return
	}

	// запрещаем делать транзакции админам и верификаторам
	if userRole != models.RoleBasic {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin and Verificator cant make transfer"})
		return
	}

	// получаем с фронта тело транзакции
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// запрещаем делать перевод не от себя
	if transaction.UserFrom != userID {
		log.Println("transaction from foreign user")
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction from foreign user"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	transactionID, err := t.transactionsService.TransactionTransfer(ctx, transaction)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactionID": transactionID})
}
