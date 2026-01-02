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
) (uuid.UUID, error) {
	query := `
				INSERT INTO bank_app.users (name, surname, email, phone_number, password)
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
	).Scan(&userID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("Error in CreateUser query: %w", err)
	}

	return userID, nil
}

// // создание верификатора
// func (db *Repo) CreateVerificator(
// 	Name string,
// 	Surname string,
// 	Email string,
// 	PhoneNumber string,
// 	Password string,
// ) (uuid.UUID, error) {
// 	query := `
// 				INSERT INTO bank_app.users (name, surname, email, phone_number, password, role)
// 				VALUES ($1, $2, $3, $4, $5, $6)
// 				RETURNING id
// 			`
// 	Role := models.RoleVerificator
// 	var userID uuid.UUID

// 	err := db.DB.QueryRow(
// 		query,
// 		Name,
// 		Surname,
// 		Email,
// 		PhoneNumber,
// 		Password,
// 		Role,
// 	).Scan(&userID)

// 	if err != nil {
// 		return uuid.Nil, fmt.Errorf("Error in CreateVerificator query: %w", err)
// 	}

// 	return userID, nil
// }

// создание админа
func (db *Repo) CreateAdmin(
	Name string,
	Surname string,
	Email string,
	PhoneNumber string,
	Password string,
	Role string,
) (uuid.UUID, error) {
	query := `
				INSERT INTO bank_app.users (name, surname, email, phone_number, password, role)
				VALUES ($1, $2, $3, $4, $5, $6)
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
		return uuid.Nil, fmt.Errorf("Error in CreateAdmin query: %w", err)
	}

	return userID, nil
}
