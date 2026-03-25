package usersservice

import (
	"bank_app/internal/storage/repos/users"
)

type UsersService struct {
	usersRepo *users.Repo
}

func NewUsersService(usersRepo *users.Repo) *UsersService {
	return &UsersService{
		usersRepo: usersRepo,
	}
}