package handlers

import (
	"bank_app/internal/jwt"
	"bank_app/internal/services"
	"bank_app/internal/storage/repos/users"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type AuthHandler struct {
	userService *services.UsersService
	jwtService  jwt.TokenService
}

func NewAuthHandler(userService *services.UsersService, jwtService jwt.TokenService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

// регистрация
func (u *AuthHandler) SignUp(c *gin.Context) {
	// получаем данные пользователя с фронта
	var User users.User
	if err := c.ShouldBindJSON(&User); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// проверка заполненности имени пользователя
	if strings.TrimSpace(User.Name) == "" {
		log.Println("Name is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Name is empty"})
		return
	}

	// проверка заполненности фамилии пользователя
	if strings.TrimSpace(User.Surname) == "" {
		log.Println("Surname is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Surname is empty"})
		return
	}

	// проверка заполненности почты
	if strings.TrimSpace(User.Surname) == "" {
		log.Println("Surname is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Surname is empty"})
		return
	}

	// проверка заполненности пароля
	if strings.TrimSpace(User.Password) == "" {
		log.Println("Pass is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Pass is empty"})
		return
	}

	// проверка заполненности номера телефона
	if strings.TrimSpace(User.PhoneNumber) == "" {
		log.Println("Phone number is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Phone number is empty"})
		return
	}

	// проверяем, нет ли пользователя с указанной почтой или номером телефона
	userCheck, err := u.userService.UserCheck(User.PhoneNumber, User.Email)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusConflict), gin.H{"error": err.Error()})
		return
	}

	// если нет - отклоняем
	if userCheck {
		log.Println("Email or phone number already exist")
		c.JSON((http.StatusConflict), gin.H{"error": "Email or phone number already exist"})
		return
	}

	// применяем роль по умолчанию
	User.Role = users.RoleUser

	// добавляем пользователя в БД
	userID, err := u.userService.UserAdd(User)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	// генерируем токен для нового пользователя
	token, err := u.jwtService.GenerateToken(userID, User.Name, User.Surname, User.Email, User.Role)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	// устанавливаем куки
	c.SetCookie("cookie", token, 3600, "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{"user_id": userID})
}

// авторизация
func (u *AuthHandler) SignIn(c *gin.Context) {
	// получаем данные пользователя с фронта
	var User users.User
	if err := c.ShouldBindJSON(&User); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// проверка заполненности пароля
	if strings.TrimSpace(User.Password) == "" {
		log.Println("Pass is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Pass is empty"})
		return
	}

	// проверка заполненности почты
	if strings.TrimSpace(User.Email) == "" {
		log.Println("Email is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Email is empty"})
		return
	}

	// проверка заполненности номера телефона
	if strings.TrimSpace(User.PhoneNumber) == "" {
		log.Println("Phone number is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Phone number is empty"})
		return
	}

	// проверка пользователя
	foundUser, err := u.userService.UserVerification(User)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusUnauthorized), gin.H{"error": err.Error()})
		return
	}

	// генерируем токен для найденного пользователя
	token, err := u.jwtService.GenerateToken(foundUser.ID, foundUser.Name, foundUser.Surname, foundUser.Email, foundUser.Role)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	// устанавливаем куки
	c.SetCookie("cookie", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"user_id": foundUser.ID})
}

// выход из профиля
func (u *AuthHandler) SignOut(c *gin.Context) {
	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
