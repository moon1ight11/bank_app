package user

import (
	"bank_app/internal/api/models"
	"fmt"
	"github.com/google/uuid"
)

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
func (u *UsersService) UserVerification(User models.UserAutorization) (models.UserGet, error) {
	foundUser, err := u.usersRepo.GetUserByEmail(User.Email)
	if err != nil {
		return models.UserGet{}, err
	}

	if foundUser.PhoneNumber != User.PhoneNumber {
		return models.UserGet{}, fmt.Errorf("wrong phone number")
	}

	if foundUser.Password != User.Password {
		return models.UserGet{}, fmt.Errorf("passwords not match")
	}

	var user models.UserGet
	user.Id = foundUser.ID
	user.Name = foundUser.Name
	user.Surname = foundUser.Surname
	user.Email = foundUser.Email
	user.PhoneNumber = foundUser.PhoneNumber
	user.Role = models.Role(foundUser.Role)

	return user, nil
}

// добавление базового пользователя
func (u *UsersService) UserAdd(User models.UserRegister) (uuid.UUID, error) {
	User.Role = models.RoleBasic

	UserID, err := u.usersRepo.CreateUser(User.Name, User.Surname, User.Email, User.PhoneNumber, User.Password, string(User.Role))
	if err != nil {
		return uuid.Nil, err
	}

	return UserID, nil
}

// создание админа/верификатора
func (u *UsersService) AdminAdd(admin models.UserRegister) (uuid.UUID, error) {
	adminID, err := u.usersRepo.CreateUser(
		admin.Name,
		admin.Surname,
		admin.Email,
		admin.PhoneNumber,
		admin.Password,
		string(admin.Role),
	)
	if err != nil {
		return uuid.Nil, err
	}

	return adminID, nil
}

// получение конкретного пользователя
func (u *UsersService) UserGet(UserID uuid.UUID) (models.UserGet, error) {
	User, err := u.usersRepo.GetUserByID(UserID)
	if err != nil {
		return models.UserGet{}, err
	}

	var user models.UserGet
	user.Id = User.ID
	user.Name = User.Name
	user.Surname = User.Surname
	user.Email = User.Email
	user.PhoneNumber = User.PhoneNumber
	user.Timezone = User.Timezone
	user.Role = models.Role(User.Role)

	return user, nil
}

// список пользователей с указанной ролью
func (u *UsersService) UsersByRoleGet(role models.Role) ([]models.UserGet, error) {
	usersRepo, err := u.usersRepo.GetUsersByRole(string(role))
	if err != nil {
		return nil, err
	}

	var usersApi []models.UserGet

	for i := range usersRepo {
		var userApi models.UserGet

		userApi.Id = usersRepo[i].ID
		userApi.Name = usersRepo[i].Name
		userApi.Surname = usersRepo[i].Surname
		userApi.Email = usersRepo[i].Email
		userApi.PhoneNumber = usersRepo[i].PhoneNumber
		userApi.Timezone = usersRepo[i].Timezone
		userApi.Role = models.Role(usersRepo[i].Role)

		usersApi = append(usersApi, userApi)
	}

	return usersApi, nil
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

// изменение роли пользователя
func (u *UsersService) RoleChange(userID uuid.UUID, role models.Role) error {
	// получаем пользователя по id
	user, err := u.UserGet(userID)
	if err != nil {
		return err
	}

	// проверяем, не совпадает ли его роль с новой
	if user.Role == role {
		return fmt.Errorf("role already set")
	}

	// меняем роль в БД
	err = u.usersRepo.UpdateRole(string(role), userID)
	if err != nil {
		return err
	}

	return nil
}
