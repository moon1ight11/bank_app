package users

import (
	"fmt"
	"github.com/google/uuid"
)

// создание пользователя
func (db *Repo) CreateUser(
	Name string,
	Surname string,
	Email string,
	PhoneNumber string,
	Password string,
	Role string,
) (uuid.UUID, error) {
	query := `
				INSERT INTO bank_app.users (name, surname, email, phone_number, password, role)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id
			`
	var userID uuid.UUID

	err := db.DB.QueryRow(
		query,
		Name,
		Surname,
		Email,
		PhoneNumber,
		Password,
		Role,
	).Scan(&userID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("Error in CreateUser query: %w", err)
	}

	return userID, nil
} 