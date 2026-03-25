package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// удаление пользователя
func (db *Repo) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `
				DELETE FROM bank_app.users
				WHERE id = $1
			`
	_, err := db.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error in DeleteUser query: %w", err)
	}

	return nil
}
