package handlers

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/api/models"
	"bank_app/internal/services/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type UsersHandler struct {
	userService *user.UsersService
	jwtService  jwt.TokenService
}

func NewUsersHandler(userService *user.UsersService, jwtService jwt.TokenService) *UsersHandler {
	return &UsersHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

// получение данных пользователя
func (u *UsersHandler) GetUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	UserId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	user, err := u.userService.UserGet(UserId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"user": user})
}

// изменение данных пользователя
func (u *UsersHandler) UpdateUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	UserId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var UpdatedUser models.UserUpdate

	// устанавливаем user_id
	UpdatedUser.ID = UserId

	// получаем обновленного пользователя с фронта
	if err := c.ShouldBindJSON(&UpdatedUser); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// если обновляется имя - чтобы было не пустое
	if UpdatedUser.Name != nil {
		if strings.TrimSpace(*UpdatedUser.Name) == "" {
			log.Println("New name is empty")
			c.JSON((http.StatusBadRequest), gin.H{"error": "New name is empty"})
			return
		}
	}

	// если обновляется пароль - чтобы не был пустым
	if UpdatedUser.Password != nil {
		if strings.TrimSpace(*UpdatedUser.Password) == "" {
			log.Println("New pass is empty")
			c.JSON((http.StatusBadRequest), gin.H{"error": "New pass is empty"})
			return
		}
	}

	// если обновляется почта
	if UpdatedUser.Email != nil {
		// проверяем, похожа ли новая почта на почту
		pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, err := regexp.MatchString(pattern, *UpdatedUser.Email)
		if err != nil {
			log.Println("Error in MatchString", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// если нет - отклоняем
		if !matched {
			log.Println("New email not looks like email")
			c.JSON(http.StatusBadRequest, gin.H{"error": "New email not looks like email"})
			return
		}
	}

	// если обновляется телефон
	if UpdatedUser.PhoneNumber != nil {
		if strings.TrimSpace(*UpdatedUser.PhoneNumber) == "" {
			log.Println("New phone number is empty")
			c.JSON((http.StatusBadRequest), gin.H{"error": "New phone number is empty"})
			return
		}
	}

	// если обновляется временная зона
	if UpdatedUser.Timezone != nil {
		if strings.TrimSpace(*UpdatedUser.Timezone) == "" {
			log.Println("New timezone is empty")
			c.JSON((http.StatusBadRequest), gin.H{"error": "New timezone is empty"})
			return
		}
	}

	// проверка уникальности новых почты и телефона
	if UpdatedUser.Email != nil || UpdatedUser.PhoneNumber != nil {
		var userNumber, userEmail string

		if UpdatedUser.PhoneNumber != nil {
			userNumber = *UpdatedUser.PhoneNumber
		}

		if UpdatedUser.Email != nil {
			userEmail = *UpdatedUser.Email
		}

		userCheck, err := u.userService.UserCheck(userNumber, userEmail)
		if err != nil {
			log.Println(err)
			c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
			return
		}

		if userCheck {
			log.Println("Phone number or email already exist")
			c.JSON((http.StatusUnauthorized), gin.H{"error": "Phone number or email already exist"})
			return
		}
	}

	// обновляем нужные поля
	err := u.userService.UserUpdate(
		UpdatedUser.Name,
		UpdatedUser.Surname,
		UpdatedUser.Password,
		UpdatedUser.Email,
		UpdatedUser.PhoneNumber,
		UpdatedUser.Timezone,
		UpdatedUser.ID,
	)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	// получение структуры обновленного пользователя
	foundUser, err := u.userService.UserGet(UpdatedUser.ID)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"UpdatedUser": foundUser})
}

// удаление пользователя
func (u *UsersHandler) DeleteUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	UserId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// находим пользователя по id
	_, err := u.userService.UserGet(UserId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// удаляем пользователя
	err = u.userService.UserDelete(UserId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"message": "delete is complete"})
}

// создание админа или верификатора
func (u *UsersHandler) CreateAdminOrVerificator(c *gin.Context) {
	// получаем с фронта пользователя
	var user models.UserRegister
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// проверка заполненности имени пользователя
	if strings.TrimSpace(user.Name) == "" {
		log.Println("Name is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Name is empty"})
		return
	}

	// проверка заполненности фамилии пользователя
	if strings.TrimSpace(user.Surname) == "" {
		log.Println("Surname is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Surname is empty"})
		return
	}

	// проверка заполненности почты
	if strings.TrimSpace(user.Surname) == "" {
		log.Println("Surname is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Surname is empty"})
		return
	}

	// проверка заполненности пароля
	if strings.TrimSpace(user.Password) == "" {
		log.Println("Pass is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Pass is empty"})
		return
	}

	// проверка заполненности номера телефона
	if strings.TrimSpace(user.PhoneNumber) == "" {
		log.Println("Phone number is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Phone number is empty"})
		return
	}

	// проверяем указана ли роль
	if user.Role == "" {
		log.Println("Role is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Role is empty"})
		return
	}

	// проверяем, нет ли пользователя с указанной почтой или номером телефона
	userCheck, err := u.userService.UserCheck(user.PhoneNumber, user.Email)
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

	userId, err := u.userService.AdminCreate(user)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"newUserId": userId})

	// // если указанная роль - админ
	// if user.Role == models.RoleAdmin {
	// 	userID, err := u.userService.AdminCreate(user)
	// 	if err != nil {
	// 		log.Println(err)
	// 		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{"newUserId": userID})

	// 	// если указанная роль - верификатор
	// } else if user.Role == models.RoleVerificator {
	// 	userID, err := u.userService.VerificatorCreate(user)
	// 	if err != nil {
	// 		log.Println(err)
	// 		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{"newUserId": userID})

	// } else {
	// 	log.Println("wrong role")
	// 	c.JSON((http.StatusBadRequest), gin.H{"error": "wrong role"})
	// 	return
	// }
}

// получение списка юзеров с фильтром на роль
func (u *UsersHandler) GetAllUsers(c *gin.Context) {
	// получаем роль из параметров
	param := c.Query("role")
	role := models.Role(param)

	users, err := u.userService.UsersByRoleGet(role)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
