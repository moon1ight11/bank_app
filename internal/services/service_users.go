package services

import (
	"bank_app/internal/storage/repos/users"
	"fmt"

	"github.com/google/uuid"
)

type UsersService struct {
	usersRepo *users.Repo
}

func NewUsersService(usersRepo *users.Repo) *UsersService {
	return &UsersService{usersRepo: usersRepo}
}

// проверка данных пользователя
func (u *UsersService) UserCheck(phoneNumber string, userEmail string) (bool, error) {
	numberExist, err := u.usersRepo.CheckUserPhoneNumber(phoneNumber)
	if err != nil {
		return true, err
	}

	emailExist, err := u.usersRepo.CheckUserEmail(userEmail)
	if err != nil {
		return true, err
	}

	if numberExist {
		return true, nil
	}

	if emailExist {
		return true, nil
	}

	return false, nil
}

// верификация пользователя
func (u *UsersService) UserVerification(User users.User) (bool, error) {
	foundUser, err := u.usersRepo.GetUserByEmail(User.Email)
	if err != nil {
		return false, err
	}

	if foundUser.Password != User.Password {
		return false, fmt.Errorf("passwords not match")
	}

	return true, nil
}

// добавление пользователя
func (u *UsersService) UserAdd(User users.User) (uuid.UUID, error) {
	UserID, err := u.usersRepo.CreateUser(User)
	if err != nil {
		return uuid.Nil, err
	}

	return UserID, nil
}
