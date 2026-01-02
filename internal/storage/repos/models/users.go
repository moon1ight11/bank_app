package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID
	Name        string
	Surname     string
	Email       string
	PhoneNumber string
	Password    string
	Timezone    string
	Role        Role
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
