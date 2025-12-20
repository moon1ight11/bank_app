package users

import (
	"fmt"
	"github.com/google/uuid"
)

// создание пользователя
func (db *Repo) CreateUser(NewUser User) (uuid.UUID, error) {
	query := `
				INSERT INTO bank_app.users (name, surname, email, phone_number, password)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id
			`
	var userID uuid.UUID

	err := db.DB.QueryRow(
		query, 
		NewUser.Name, 
		NewUser.Surname, 
		NewUser.Email, 
		NewUser.PhoneNumber, 
		NewUser.Password,
		).Scan(&userID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("Error in CreateUser query: %w", err)
	}

	return userID, nil
}

// создание верификатора
func (db *Repo) CreateVerificator(NewVerificator User) (uuid.UUID, error) {
	query := `
				INSERT INTO bank_app.users (name, surname, email, phone_number, password, role)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id
			`
	NewVerificator.Role = RoleVerificator
	var userID uuid.UUID

	err := db.DB.QueryRow(
		query, 
		NewVerificator.Name,
		NewVerificator.Surname, 
		NewVerificator.Email, 
		NewVerificator.PhoneNumber, 
		NewVerificator.Password, 
		NewVerificator.Role,
		).Scan(&userID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("Error in CreateVerificator query: %w", err)
	}

	return userID, nil
}

// создание админа
func (db *Repo) CreateAdmin(NewAdmin User) (uuid.UUID, error) {
	query := `
				INSERT INTO bank_app.users (name, surname, email, phone_number, password, role)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id
			`
	NewAdmin.Role = RoleAdmin
	var userID uuid.UUID

	err := db.DB.QueryRow(
		query, 
		NewAdmin.Name,
		NewAdmin.Surname, 
		NewAdmin.Email, 
		NewAdmin.PhoneNumber, 
		NewAdmin.Password, 
		NewAdmin.Role,
		).Scan(&userID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("Error in CreateAdmin query: %w", err)
	}

	return userID, nil
}
