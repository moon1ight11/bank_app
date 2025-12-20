package users

import (
	"bank_app/internal/storage"
	"github.com/google/uuid"
)

type Repo struct {
	storage.DataBase
}

func NewUsersRepo(db *storage.DataBase) *Repo {
	return &Repo{DataBase: *db}
}

type User struct {
	ID          uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Email       string    `json:"email" binding:"required,email"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"password"`
	Timezone    string    `json:"timezone"`
	Role        Role      `json:"roles"`
}

type Role string

const (
	RoleUser        Role = "User"
	RoleVerificator Role = "Verificator"
	RoleAdmin       Role = "Admin"
)
