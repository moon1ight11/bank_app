package users

import (
	"fmt"
	"github.com/google/uuid"
)

// получение данных о пользователе по id
func (db *Repo) GetUserByID(userId uuid.UUID) (User, error) {
	query := `
				SELECT id, name, surname, email, phone_number, timezone
				FROM bank_app.users
				WHERE id = $1
			`
	var user User

	err := db.DB.QueryRow(query, userId).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.PhoneNumber, &user.Timezone)
	if err != nil {
		return User{}, fmt.Errorf("Error in GetUserByID query: %w", err)
	}

	return user, nil
}

// получение данных о пользователе по email
func (db *Repo) GetUserByEmail(userEmail string) (User, error) {
	query := `
				SELECT id, name, surname, password, email, phone_number, timezone
				FROM bank_app.users
				WHERE email = $1
			`
	var user User

	err := db.DB.QueryRow(query, userEmail).Scan(&user.ID, &user.Name, &user.Surname, &user.Password, &user.Email, &user.PhoneNumber, &user.Timezone)
	if err != nil {
		return User{}, fmt.Errorf("Error in GetUserByEmail query: %w", err)
	}

	return user, nil
}

// проверка свободности почты
func (db *Repo) CheckUserEmail(userEmail string) (bool, error) {
	query := `
				SELECT EXISTS(
					SELECT 1
					FROM bank_app.users
					WHERE email = $1)
			`
	var exist bool

	err := db.DB.QueryRow(query, userEmail).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("Error in CheckUserEmail query: %w", err)
	}

	return exist, nil
}

// проверка свободности номера телефона
func (db *Repo) CheckUserPhoneNumber(phoneNumber string) (bool, error) {
	query := `
				SELECT EXISTS(
					SELECT 1
					FROM bank_app.users
					WHERE phone_number = $1)
			`
	var exist bool

	err := db.DB.QueryRow(query, phoneNumber).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("Error in CheckUserPhoneNumber query: %w", err)
	}

	return exist, nil
}

// список пользователей с заданной ролью
func (db *Repo) GetUsersByRole(role Role) ([]User, error) {
	query := `
				SELECT id, name, surname, email, timezone
				FROM bank_app.users
				WHERE role = $1
			`
	var users []User
	rows, err := db.DB.Query(query, role)
	if err != nil {
		return nil, fmt.Errorf("error in GetUsersByRole query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.Timezone,
		)
		if err != nil {
			return nil, fmt.Errorf("error in GetUsersByRole scan: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
