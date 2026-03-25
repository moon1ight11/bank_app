package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// обновление имени пользователя
func (db *Repo) UpdateName(ctx context.Context, newName string, userID uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.users
				SET name = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.ExecContext(ctx, query, newName, userID)
	if err != nil {
		return fmt.Errorf("error in UpdateName query: %w", err)
	}

	return nil
}

// обновление фамилии
func (db *Repo) UpdateSurname(ctx context.Context, newSurname string, userID uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.users
				SET surname = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.ExecContext(ctx, query, newSurname, userID)
	if err != nil {
		return fmt.Errorf("error in UpdateSurname query: %w", err)
	}

	return nil
}

// обновление почты
func (db *Repo) UpdateEmail(ctx context.Context, newEmail string, userID uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.users
				SET email = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.ExecContext(ctx, query, newEmail, userID)
	if err != nil {
		return fmt.Errorf("error in UpdateEmail query: %w", err)
	}

	return nil
}

// обновление пароля
func (db *Repo) UpdatePass(ctx context.Context, newPass string, userID uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.users
				SET password = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.ExecContext(ctx, query, newPass, userID)
	if err != nil {
		return fmt.Errorf("error in UpdatePass query: %w", err)
	}

	return nil
}

// обновление номера телефона
func (db *Repo) UpdatePhone(ctx context.Context, newPhone string, userID uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.users
				SET phone_number = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.ExecContext(ctx, query, newPhone, userID)
	if err != nil {
		return fmt.Errorf("error in UpdatePhone query: %w", err)
	}

	return nil
}

// обновление временной зоны
func (db *Repo) UpdateTZ(ctx context.Context, newTZ string, userID uuid.UUID, tx *sql.Tx) error {
	query := `
				UPDATE bank_app.users
				SET timezone = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := tx.ExecContext(ctx, query, newTZ, userID)
	if err != nil {
		return fmt.Errorf("error in UpdateTZ query: %w", err)
	}

	return nil
}

// изменение роли
func (db *Repo) UpdateRole(ctx context.Context, role string, userID uuid.UUID) error {
	query := `
				UPDATE bank_app.users
				SET role = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.ExecContext(ctx, query, role, userID)
	if err != nil {
		return fmt.Errorf("error in UpdateRole query: %w", err)
	}

	return nil
}
