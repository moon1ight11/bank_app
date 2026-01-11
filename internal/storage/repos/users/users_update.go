package users

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// обновление имени пользователя
func (db *Repo) UpdateName(newName string, userID uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.users
				SET name = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.Exec(query, newName, userID)
	if err != nil {
		return fmt.Errorf("Error in UpdateName query: %w", err)
	}

	return nil
}

// обновление фамилии
func (db *Repo) UpdateSurname(newSurname string, userID uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.users
				SET surname = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.Exec(query, newSurname, userID)
	if err != nil {
		return fmt.Errorf("Error in UpdateSurname query: %w", err)
	}

	return nil
}

// обновление почты
func (db *Repo) UpdateEmail(newEmail string, userID uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.users
				SET email = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.Exec(query, newEmail, userID)
	if err != nil {
		return fmt.Errorf("Error in UpdateEmail query: %w", err)
	}

	return nil
}

// обновление пароля
func (db *Repo) UpdatePass(newPass string, userID uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.users
				SET password = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.Exec(query, newPass, userID)
	if err != nil {
		return fmt.Errorf("Error in UpdatePass query: %w", err)
	}

	return nil
}

// обновление номера телефона
func (db *Repo) UpdatePhone(newPhone string, userID uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.users
				SET phone_number = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.Exec(query, newPhone, userID)
	if err != nil {
		return fmt.Errorf("Error in UpdatePhone query: %w", err)
	}

	return nil
}

// обновление временной зоны
func (db *Repo) UpdateTZ(newTZ string, userID uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.users
				SET timezone = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.Exec(query, newTZ, userID)
	if err != nil {
		return fmt.Errorf("Error in UpdateTZ query: %w", err)
	}

	return nil
}

// изменение роли
func (db *Repo) UpdateRole(role string, userID uuid.UUID) error {
	query := `
				UPDATE bank_app.users
				SET role = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, role, userID)
	if err != nil {
		return fmt.Errorf("Error in UpdateRole query: %w", err)
	}

	return nil
}
