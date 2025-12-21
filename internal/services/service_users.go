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
func (u *UsersService) UserVerification(User users.User) (users.User, error) {
	foundUser, err := u.usersRepo.GetUserByEmail(User.Email)
	if err != nil {
		return users.User{}, err
	}

	if foundUser.PhoneNumber != User.PhoneNumber {
		return users.User{}, fmt.Errorf("wrong phone number")
	}

	if foundUser.Password != User.Password {
		return users.User{}, fmt.Errorf("passwords not match")
	}

	return foundUser, nil
}

// добавление пользователя
func (u *UsersService) UserAdd(User users.User) (uuid.UUID, error) {
	UserID, err := u.usersRepo.CreateUser(User)
	if err != nil {
		return uuid.Nil, err
	}

	return UserID, nil
}

// получение пользователя
func (u *UsersService) UserGet(UserID uuid.UUID) (users.User, error) {
	user, err := u.usersRepo.GetUserByID(UserID)
	if err != nil {
		return users.User{}, err
	}

	return user, nil
}

// обновление пользователя
func (u *UsersService) UserUpdate(
	name *string,
	surname *string,
	password *string,
	email *string,
	phone *string,
	tz *string,
	ID uuid.UUID,
) error {
	// открываем транзакцию
	transaction, err := u.usersRepo.DB.Begin()
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	if name != nil {
		err = u.usersRepo.UpdateName(*name, ID, transaction)
		if err != nil {
			return err
		}
	}

	if surname != nil {
		err = u.usersRepo.UpdateSurname(*surname, ID, transaction)
		if err != nil {
			return err
		}
	}

	if password != nil {
		err = u.usersRepo.UpdatePass(*password, ID, transaction)
		if err != nil {
			return err
		}
	}

	if email != nil {
		err = u.usersRepo.UpdateEmail(*email, ID, transaction)
		if err != nil {
			return err
		}
	}

	if phone != nil {
		err = u.usersRepo.UpdatePhone(*phone, ID, transaction)
		if err != nil {
			return err
		}
	}

	if tz != nil {
		err = u.usersRepo.UpdateTZ(*tz, ID, transaction)
		if err != nil {
			return err
		}
	}

	// если все ок - подтверждаем транзакцию
	transaction.Commit()
	return nil
}

// удаление пользователя
func (u *UsersService) UserDelete(userID uuid.UUID) error {
	err := u.usersRepo.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}

// создание верификатора
func (u *UsersService) VerificatorCreate(verificator users.User) (uuid.UUID, error) {
	verificatorID, err := u.usersRepo.CreateVerificator(verificator)
	if err != nil {
		return uuid.Nil, err
	}

	return verificatorID, nil
}

// создание админа
func (u *UsersService) AdminCreate(admin users.User) (uuid.UUID, error) {
	adminID, err := u.usersRepo.CreateAdmin(admin)
	if err != nil {
		return uuid.Nil, err
	}

	return adminID, nil
}

// список админов
func (u *UsersService) AdminsGet() ([]users.User, error) {
	admins, err := u.usersRepo.GetUsersByRole(users.RoleAdmin)
	if err != nil {
		return nil, err
	}

	return admins, nil
}

// список верификаторов
func (u *UsersService) VerificatorsGet() ([]users.User, error) {
	verificators, err := u.usersRepo.GetUsersByRole(users.RoleVerificator)
	if err != nil {
		return nil, err
	}

	return verificators, nil
}
