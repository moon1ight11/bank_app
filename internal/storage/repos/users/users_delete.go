package users

import (
	"fmt"
	"github.com/google/uuid"
)

// удаление пользователя
func (db *Repo) DeleteUser(userID uuid.UUID) error {
	query := `
				DELETE FROM bank_app.users
				WHERE id = $1
			`
	_, err := db.DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("Error in DeleteUser query: %w", err)
	}

	return nil
}
