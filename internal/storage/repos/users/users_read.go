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
