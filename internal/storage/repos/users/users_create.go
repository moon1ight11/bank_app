package users

import (
	"fmt"
	"github.com/google/uuid"
)

// создание пользователя
func (db *Repo) CreateUser(NewUser User, password string) (uuid.UUID, error) {
	query := `
				INSERT INTO bank_app.users (name, surname, email, phone_number, password)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id
			`
	var userID uuid.UUID
	err := db.DB.QueryRow(query, NewUser.Name, NewUser.Surname, NewUser.Email, NewUser.PhoneNumber, password).Scan(&userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error in CreateUser query: %w", err)
	}

	return userID, nil
}
