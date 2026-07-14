package usersservice

import (
	"bank_app/internal/api/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"time"
)

// создание базового пользователя
func (u *UsersService) UserAdd(ctx context.Context, user models.UserRegister) (uuid.UUID, error) {
	if isValidEmail(user.Email) == false {
		return uuid.Nil, fmt.Errorf("error in UserAdd: email is not valid")
	}

	if isValidPhoneNumber(user.PhoneNumber) == false {
		return uuid.Nil, fmt.Errorf("error in UserAdd: phone number is not valid")
	}

	check, err := u.UserCheck(ctx, user.PhoneNumber, user.Email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in UserAdd: %w", err)
	}

	if check {
		return uuid.Nil, fmt.Errorf("error in UserAdd: userCheck is failed")
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in AddUser hash password: %w", err)
	}

	UserID, err := u.usersRepo.CreateUser(ctx, user.Name, user.Surname, user.Email, user.PhoneNumber, string(hashPass), string(user.Role))
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in UserAdd: %w", err)
	}

	cacheKey := fmt.Sprintf("user:%s", UserID.String())
	if u.cacheService != nil {
		if err := u.cacheService.Set(ctx, cacheKey, UserID, 10*time.Minute); err != nil {
			return uuid.Nil, fmt.Errorf("error in userAdd: %w; cachekey %s not set", err, cacheKey)
		}
	}

	return UserID, nil
}

// создание админа/верификатора
func (u *UsersService) AdminAdd(ctx context.Context, admin models.UserRegister) (uuid.UUID, error) {
	if isValidEmail(admin.Email) == false {
		return uuid.Nil, fmt.Errorf("error in AdminAdd: email is not valid")
	}

	if isValidPhoneNumber(admin.PhoneNumber) == false {
		return uuid.Nil, fmt.Errorf("error in AdminAdd: phone number is not valid")
	}

	check, err := u.UserCheck(ctx, admin.PhoneNumber, admin.Email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in AdminAdd: %w", err)
	}
	if check {
		return uuid.Nil, fmt.Errorf("error in AdminAdd: userCheck is failed")
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in AddUser hash password: %w", err)
	}

	adminID, err := u.usersRepo.CreateUser(
		ctx,
		admin.Name,
		admin.Surname,
		admin.Email,
		admin.PhoneNumber,
		string(hashPass),
		string(admin.Role),
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in AdminAdd: %w", err)
	}

	cacheKey := fmt.Sprintf("user:%s", adminID.String())
	if u.cacheService != nil {
		if err := u.cacheService.Set(ctx, cacheKey, adminID, 10*time.Minute); err != nil {
			return uuid.Nil, fmt.Errorf("error in adminAdd: %w; cachekey %s not set", err, cacheKey)
		}
	}

	return adminID, nil
}

// проверка уникальности данных пользователя
func (u *UsersService) UserCheck(ctx context.Context, phoneNumber string, userEmail string) (bool, error) {
	numberExist, err := u.usersRepo.CheckUserPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return true, fmt.Errorf("error in UserCheck: %w", err)
	}

	emailExist, err := u.usersRepo.CheckUserEmail(ctx, userEmail)
	if err != nil {
		return true, fmt.Errorf("error in UserCheck: %w", err)
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
func (u *UsersService) UserVerification(ctx context.Context, userAutoriz models.UserAutorization) (models.UserGet, error) {
	foundUser, err := u.usersRepo.GetUserByEmail(ctx, userAutoriz.Email)
	if err != nil {
		return models.UserGet{}, fmt.Errorf("error in UserVerification: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userAutoriz.Password))
	if err != nil {
		return models.UserGet{}, fmt.Errorf("error in UserVerification: invalid credentials")
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

// получение конкретного пользователя
func (u *UsersService) UserGet(ctx context.Context, userID uuid.UUID) (models.UserGet, error) {
	cacheKey := fmt.Sprintf("user:%s", userID.String())
	if u.cacheService != nil {
		var cachedUser models.UserGet
		err := u.cacheService.Get(ctx, cacheKey, &cachedUser)
		if err == nil {
			return cachedUser, nil
		}
	}

	userRepo, err := u.usersRepo.GetUserByID(ctx, userID)
	if err != nil {
		return models.UserGet{}, fmt.Errorf("error in UserGet: %w", err)
	}

	var user models.UserGet
	user.Id = userRepo.ID
	user.Name = userRepo.Name
	user.Surname = userRepo.Surname
	user.Email = userRepo.Email
	user.PhoneNumber = userRepo.PhoneNumber
	user.Timezone = userRepo.Timezone
	user.Role = models.Role(userRepo.Role)

	if u.cacheService != nil {
		if err := u.cacheService.Set(ctx, cacheKey, userID, 10*time.Minute); err != nil {
			return models.UserGet{}, fmt.Errorf("error in userGet: %w; cachekey %s not set", err, cacheKey)
		}
	}

	return user, nil
}

// получение списка пользователей с указанной ролью
func (u *UsersService) UsersByRoleGet(ctx context.Context, role models.Role) ([]models.UserGet, error) {
	usersRepo, err := u.usersRepo.GetUsersByRole(ctx, string(role))
	if err != nil {
		return nil, fmt.Errorf("error in UsersByRoleGet: %w", err)
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
	ctx context.Context,
	name *string,
	surname *string,
	password *string,
	email *string,
	phone *string,
	tz *string,
	ID uuid.UUID,
) error {
	transaction, err := u.usersRepo.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error in UserUpdate: %w", err)
	}
	defer transaction.Rollback()

	if name != nil {
		err = u.usersRepo.UpdateName(ctx, *name, ID, transaction)
		if err != nil {
			return fmt.Errorf("error in UserUpdate: %w", err)
		}
	}

	if surname != nil {
		err = u.usersRepo.UpdateSurname(ctx, *surname, ID, transaction)
		if err != nil {
			return fmt.Errorf("error in UserUpdate: %w", err)
		}
	}

	if password != nil {
		hashPass, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error in UserUpdate hash password: %w", err)
		}

		err = u.usersRepo.UpdatePass(ctx, string(hashPass), ID, transaction)
		if err != nil {
			return fmt.Errorf("error in UserUpdate: %w", err)
		}
	}

	if email != nil {
		if isValidEmail(*email) == false {
			return fmt.Errorf("error in UserUpdate: email is not valid")
		}

		exist, err := u.usersRepo.CheckUserEmailTx(ctx, *email, transaction)
		if err != nil {
			return fmt.Errorf("error in UserUpdate: %w", err)
		}
		if exist {
			return fmt.Errorf("error in UserUpdate: CheckEmail is failed")
		}

		err = u.usersRepo.UpdateEmail(ctx, *email, ID, transaction)
		if err != nil {
			return fmt.Errorf("error in UserUpdate: %w", err)
		}
	}

	if phone != nil {
		if isValidPhoneNumber(*phone) == false {
			return fmt.Errorf("error in UserUpdate: phone number is not valid")
		}

		exist, err := u.usersRepo.CheckUserPhoneNumberTx(ctx, *phone, transaction)
		if err != nil {
			return fmt.Errorf("error in UserUpdate: %w", err)
		}
		if exist {
			return fmt.Errorf("error in UserUpdate: CheckPhone is failed")
		}

		err = u.usersRepo.UpdatePhone(ctx, *phone, ID, transaction)
		if err != nil {
			return fmt.Errorf("error in UserUpdate: %w", err)
		}
	}

	if tz != nil {
		if isValidTimezone(*tz) == false {
			return fmt.Errorf("error in UserUpdate: timezone is not valid")
		}

		err = u.usersRepo.UpdateTZ(ctx, *tz, ID, transaction)
		if err != nil {
			return fmt.Errorf("error in UserUpdate: %w", err)
		}
	}

	if err = transaction.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	cacheKey := fmt.Sprintf("user:%s", ID.String())

	if u.cacheService != nil {
		if err := u.cacheService.Delete(ctx, cacheKey); err != nil {
			return fmt.Errorf("error in userUpdate: %w; cachekey %s not delete", err, cacheKey)
		}
		if err := u.cacheService.Set(ctx, cacheKey, ID, 10*time.Minute); err != nil {
			return fmt.Errorf("error in userUpdate: %w; cachekey %s not set", err, cacheKey)
		}
	}

	return nil
}

// удаление пользователя
func (u *UsersService) UserDelete(ctx context.Context, userID uuid.UUID) error {
	err := u.usersRepo.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("error in UserDelete: %w", err)
	}

	cacheKey := fmt.Sprintf("user:%s", userID.String())
	if u.cacheService != nil {
		if err := u.cacheService.Delete(ctx, cacheKey); err != nil {
			return fmt.Errorf("error in userUpdate: %w; cachekey %s not delete", err, cacheKey)
		}
	}

	return nil
}

// изменение роли пользователя
func (u *UsersService) RoleChange(ctx context.Context, userID uuid.UUID, role models.Role) error {
	user, err := u.UserGet(ctx, userID)
	if err != nil {
		return fmt.Errorf("error in RoleChange: %w", err)
	}

	if user.Role == role {
		return fmt.Errorf("error in RoleChange: role already set")
	}

	err = u.usersRepo.UpdateRole(ctx, string(role), userID)
	if err != nil {
		return fmt.Errorf("error in RoleChange: %w", err)
	}

	cacheKey := fmt.Sprintf("user:%s", userID.String())

	if u.cacheService != nil {
		if err := u.cacheService.Delete(ctx, cacheKey); err != nil {
			return fmt.Errorf("error in userUpdate: %w; cachekey %s not delete", err, cacheKey)
		}
	}

	return nil
}

// проверка номера телефона
func isValidPhoneNumber(phone string) bool {
	cleaned := strings.TrimSpace(phone)
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	patterns := []string{
		`^\+?[0-9]{10,15}$`,
		`^8[0-9]{10}$`,
		`^\+7[0-9]{10}$`,
		`^[0-9]{10,15}$`,
	}

	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, cleaned)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// проверка временной зоны
func isValidTimezone(timezone string) bool {
	pattern := `^UTC\+((0[0-9]|1[0-4]):[0-5][0-9]|14:00)$|^UTC-((0[0-9]|1[0-2]):[0-5][0-9])$`
	matched, err := regexp.MatchString(pattern, timezone)
	if err != nil {
		return false
	}

	if !matched {
		return false
	}

	return true
}

// проверка почты
func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}

	if !matched {
		return false
	}

	return true
}
