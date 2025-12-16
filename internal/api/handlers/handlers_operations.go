package handlers

import (
	"bank_app/internal/jwt"
	"bank_app/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OperationsHandler struct {
	operationsService *services.OperationsService
	jwtService        jwt.TokenService
}

func NewOperationsHandler(operationsService *services.OperationsService, jwtService jwt.TokenService) *OperationsHandler {
	return &OperationsHandler{
		operationsService: operationsService,
		jwtService:        jwtService,
	}
}

// получение всех операций пользователя
func (o *OperationsHandler) GetAllUserOperations(c *gin.Context) {
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

	operations, err := o.operationsService.AllOperationsGet(userID)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"operations": operations})
}

// получение всех операций конкретного счета
func (o *OperationsHandler) GetAllAccountOperations(c *gin.Context) {
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

	operations, err := o.operationsService.AccountOperationsGet(userID, accountID)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"operations": operations})
}

// информация о конкретной операции
func (o *OperationsHandler) GetOperationByID(c *gin.Context) {
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

	// получаем id операции из параметров
	idStr := c.Param("operation_id")
	operationID, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": "Error in parse uuid"})
		return
	}

	operation, err := o.operationsService.OperationByIdGet(userID, operationID)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"operation": operation})
}
