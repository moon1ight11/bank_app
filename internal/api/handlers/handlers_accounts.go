package handlers

import (
	"bank_app/internal/jwt"
	"bank_app/internal/services"
	"bank_app/internal/storage/repos/accounts"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type AccountsHandler struct {
	accountsService *services.AccountsService
	jwtService      jwt.TokenService
}

func NewAccountsHandler(accountsService *services.AccountsService, jwtService jwt.TokenService) *AccountsHandler {
	return &AccountsHandler{
		accountsService: accountsService,
		jwtService:      jwtService,
	}
}

// создание нового счёта
func (a *AccountsHandler) CreateAccount(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	UserID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var NewAccount accounts.Account

	// устанавливаем владельцем счета пользователя
	NewAccount.OwnerID = UserID

	// создаем счет
	_, err := a.accountsService.AccountAdd(NewAccount)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"message": "sucessful"})
}

// список счетов пользователя
func (a *AccountsHandler) GetAllUserAccounts(c *gin.Context) {
	// получаем id из контекста
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

	// получаем список счетов конкретного пользователя
	Accounts, err := a.accountsService.AllAccountsGet(userID)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"accounts": Accounts})
}

// конкретный счет пользователя
func (a *AccountsHandler) GetAccountById(c *gin.Context) {
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

	// получаем счет из БД
	Account, err := a.accountsService.AccountGet(userID, accountID)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"account": Account})
}

// удаление счета
func (a *AccountsHandler) DeleteAccount(c *gin.Context) {
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

	// удаляем счет
	err = a.accountsService.AccountDelete(userID, accountID)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"message": "successfull delete"})
}

// пополнение счета
func (a *AccountsHandler) BalanceIncoming(c *gin.Context) {
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

	// из параметров получаем id счета, который пополняется
	idStr := c.Param("account_id")
	accountID, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": "Error in parse uuid"})
		return
	}

	// получаем с фронта сумму, на которую пополняется счет
	var amount float64
	if err := c.ShouldBindJSON(&amount); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// пополнение счёта + запись в операциях
	err = a.accountsService.AccountIncoming(userID, accountID, amount)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"message": "successfull incoming"})
}

// списание со счета
func (a *AccountsHandler) BalanceOutlay(c *gin.Context) {
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

	// из параметров получаем id счета, с которого снимаем
	idStr := c.Param("account_id")
	accountID, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": "Error in parse uuid"})
		return
	}

	// получаем с фронта сумму, которая снимается
	var amount float64
	if err := c.ShouldBindJSON(&amount); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// снятие со счета + запись в операции
	err = a.accountsService.AccountOutlay(userID, accountID, amount)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"message": "successfull outlay"})
}

// перевод
func (a *AccountsHandler) BalanceTransfer(c *gin.Context) {
	// получаем id пользователя из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userOutID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// из параметров получаем id счета, с которого снимаем
	idStr := c.Param("account_id")
	accountOutID, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": "Error in parse uuid"})
		return
	}

	// с фронта получаем id юзера и счета, на который переводим и сумму перевода
	var transferTo struct {
		userInID    uuid.UUID
		accountInID uuid.UUID
		amount      float64
	}

	if err := c.ShouldBindJSON(&transferTo); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// выполняем трансфер
	err = a.accountsService.AccountTransfer(transferTo.userInID, transferTo.accountInID, userOutID, accountOutID, transferTo.amount)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"message": "successfull transfer"})
}
