package handlers

import (
	"bank_app/internal/jwt"
	"bank_app/internal/services"
	"bank_app/internal/storage/repos/transactions"
	"bank_app/internal/storage/repos/users"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionsHandler struct {
	transactionsService *services.TransactionsService
	jwtService          jwt.TokenService
}

func NewTransactionsHandler(transactionsService *services.TransactionsService, jwtService jwt.TokenService) *TransactionsHandler {
	return &TransactionsHandler{
		transactionsService: transactionsService,
		jwtService:          jwtService,
	}
}

// получение всех транзакций пользователя
func (t *TransactionsHandler) GetAllUserTransactions(c *gin.Context) {
	// получаем id пользователя из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	transactions, err := t.transactionsService.AllTransactionsGet(userID)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"transactions": transactions})
}

// получение всех транзакций конкретного счета
func (t *TransactionsHandler) GetAllAccountTransactions(c *gin.Context) {
	// получаем id пользователя из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// получаем id счета из параметров
	idStr := c.Param("account_id")
	accountID, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": "Error in parse uuid"})
		return
	}

	transactions, err := t.transactionsService.AccountTransactionsGet(userID, accountID)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"transactions": transactions})
}

// информация о конкретной транзакции
func (t *TransactionsHandler) GetTransactionByID(c *gin.Context) {
	// получаем id пользователя из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// получаем id транзакции из параметров
	idStr := c.Param("transaction_id")
	transactionID, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": "Error in parse uuid"})
		return
	}

	transaction, err := t.transactionsService.TransactionByIdGet(userID, transactionID)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"transaction": transaction})
}

// пополнение счета
func (t *TransactionsHandler) CreateIncomingTransaction(c *gin.Context) {
	// получаем с фронта тело транзакции
	var transaction transactions.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	transactionID, err := t.transactionsService.TransactionIncoming(transaction)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"transactionID": transactionID})
}

// списание со счета
func (t *TransactionsHandler) CreateOutcomingTransaction(c *gin.Context) {
	// получаем с фронта тело транзакции
	var transaction transactions.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	transactionID, err := t.transactionsService.TransactionOutcoming(transaction)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"transactionID": transactionID})
}

// трансфер
func (t *TransactionsHandler) CreateTransferTransaction(c *gin.Context) {
	// получаем роль пользователя из контекста
	userRole, exist := c.Get("UserRole")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
		return
	}

	// запрещаем делать транзакции админам и верификаторам
	if userRole != users.RoleUser {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin and Verificator cant make transfer"})
		return
	}

	// получаем с фронта тело транзакции
	var transaction transactions.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	transactionID, err := t.transactionsService.TransactionTransfer(transaction)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"transactionID": transactionID})
}
