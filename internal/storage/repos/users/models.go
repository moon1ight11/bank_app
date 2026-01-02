package users

import (
	"bank_app/internal/storage/repos/models"
	"github.com/google/uuid"
)

type GetUser struct {
	ID          uuid.UUID
	Name        string
	Surname     string
	Email       string
	PhoneNumber string
	Password    string
	Timezone    string
	Role        models.Role
}
