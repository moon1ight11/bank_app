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
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Role    Role   `json:"role"`
	UserAutorization
}

// авторизация юзера
type UserAutorization struct {
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

// обновление юзера
type UserUpdate struct {
	ID          uuid.UUID `json:"user_id"`
	Name        *string   `json:"name"`
	Surname     *string   `json:"surname"`
	Password    *string   `json:"password"`
	PhoneNumber *string   `json:"phone_number"`
	Email       *string   `json:"email"`
	Timezone    *string   `json:"timezone"`
}

type UserGet struct {
	Id       uuid.UUID `json:"user_id"`
	Timezone string    `json:"timezone"`
	UserRegister
}
