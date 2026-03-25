package models

import "github.com/google/uuid"

// роли юзеров
type Role string

const (
	RoleBasic       Role = "Basic"
	RoleVerificator Role = "Verificator"
	RoleAdmin       Role = "Admin"
)

// регистрация юзера
type UserRegister struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`
	Surname     string `json:"surname" binding:"required,min=2,max=50"`
	PhoneNumber string `json:"phone_number"`
	Role        Role   `json:"role"`
	UserAutorization
}

// авторизация юзера
type UserAutorization struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password"`
}

// обновление юзера
type UserUpdate struct {
	ID          uuid.UUID `json:"user_id"`
	Name        *string   `json:"name"`
	Surname     *string   `json:"surname"`
	Password    *string   `json:"password"`
	PhoneNumber *string   `json:"phone_number"`
	Email       *string   `json:"email" binding:"required,email"`
	Timezone    *string   `json:"timezone"`
}

// получение пользователя
type UserGet struct {
	Id       uuid.UUID `json:"user_id"`
	Timezone string    `json:"timezone"`
	UserRegister
}
