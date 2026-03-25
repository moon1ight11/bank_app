package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// получение данных о пользователе по id
func (db *Repo) GetUserByID(ctx context.Context, userId uuid.UUID) (GetUser, error) {
	query := `
				SELECT id, name, surname, email, phone_number, timezone, role
				FROM bank_app.users
				WHERE id = $1
			`
	var user GetUser

	err := db.DB.QueryRowContext(ctx, query, userId).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.PhoneNumber, &user.Timezone, &user.Role)
	if err != nil {
		return GetUser{}, fmt.Errorf("error in GetUserByID query: %w", err)
	}

	return user, nil
}

// получение данных о пользователе по email
func (db *Repo) GetUserByEmail(ctx context.Context, userEmail string) (GetUser, error) {
	query := `
				SELECT id, name, surname, password, email, phone_number, timezone, role
				FROM bank_app.users
				WHERE email = $1
			`
	var user GetUser

	err := db.DB.QueryRowContext(ctx, query, userEmail).Scan(&user.ID, &user.Name, &user.Surname, &user.Password, &user.Email, &user.PhoneNumber, &user.Timezone, &user.Role)
	if err != nil {
		return GetUser{}, fmt.Errorf("error in GetUserByEmail query: %w", err)
	}

	return user, nil
}

// проверка свободности почты
func (db *Repo) CheckUserEmail(ctx context.Context, userEmail string) (bool, error) {
	query := `
				SELECT EXISTS(
					SELECT 1
					FROM bank_app.users
					WHERE email = $1)
			`
	var exist bool

	err := db.DB.QueryRowContext(ctx, query, userEmail).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("error in CheckUserEmail query: %w", err)
	}

	return exist, nil
}

// проверка свободности почты с транзакцией
func (db *Repo) CheckUserEmailTx(ctx context.Context, userEmail string, tx *sql.Tx) (bool, error) {
	query := `
				SELECT EXISTS(
					SELECT 1
					FROM bank_app.users
					WHERE email = $1)
			`
	var exist bool

	err := tx.QueryRowContext(ctx, query, userEmail).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("error in CheckUserEmail query: %w", err)
	}

	return exist, nil
}

// проверка свободности номера телефона
func (db *Repo) CheckUserPhoneNumber(ctx context.Context, phoneNumber string) (bool, error) {
	query := `
				SELECT EXISTS(
					SELECT 1
					FROM bank_app.users
					WHERE phone_number = $1)
			`
	var exist bool

	err := db.DB.QueryRowContext(ctx, query, phoneNumber).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("error in CheckUserPhoneNumber query: %w", err)
	}

	return exist, nil
}

// проверка свободности номера телефона с транзакцией
func (db *Repo) CheckUserPhoneNumberTx(ctx context.Context, phoneNumber string, tx *sql.Tx) (bool, error) {
	query := `
				SELECT EXISTS(
					SELECT 1
					FROM bank_app.users
					WHERE phone_number = $1)
			`
	var exist bool

	err := tx.QueryRowContext(ctx, query, phoneNumber).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("error in CheckUserPhoneNumber query: %w", err)
	}

	return exist, nil
}

// список пользователей с заданной ролью
func (db *Repo) GetUsersByRole(ctx context.Context, role string) ([]GetUser, error) {
	query := `
				SELECT id, name, surname, email, phone_number, timezone, role
				FROM bank_app.users
				WHERE role = $1
			`
	var users []GetUser
	rows, err := db.DB.QueryContext(ctx, query, role)
	if err != nil {
		return nil, fmt.Errorf("error in GetUsersByRole query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user GetUser
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.PhoneNumber,
			&user.Timezone,
			&user.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("error in GetUsersByRole scan: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
