package usersservice

import (
	"bank_app/internal/api/models"
	"bank_app/internal/storage/repos/users"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newIntegrationService() *UsersService {
	repo := users.NewUsersRepo(testDB)
	return NewUsersService(repo, nil)
}

func TestIntegration_UserAdd_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newIntegrationService()

	user := models.UserRegister{
		Name:        "Пётр",
		Surname:     "Петров",
		PhoneNumber: "+79998887766",
		UserAutorization: models.UserAutorization{
			Email:    "petr@example.com",
			Password: "password123",
		},
		Role: models.RoleBasic,
	}

	id, err := svc.UserAdd(context.Background(), user)
	require.NoError(t, err)
	assert.NotEqual(t, id.String(), "")

	found, err := svc.UserGet(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, "Пётр", found.Name)
	assert.Equal(t, "petr@example.com", found.Email)
}

func TestIntegration_UserAdd_DuplicateEmail(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newIntegrationService()

	user := models.UserRegister{
		Name:        "Дубль",
		Surname:     "Дублёв",
		PhoneNumber: "+79991112233",
		UserAutorization: models.UserAutorization{
			Email:    "duplicate@example.com",
			Password: "password123",
		},
		Role: models.RoleBasic,
	}

	_, err := svc.UserAdd(context.Background(), user)
	require.NoError(t, err)

	user.PhoneNumber = "+79991112244"
	_, err = svc.UserAdd(context.Background(), user)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "userCheck is failed")
}

func TestIntegration_UserVerification_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newIntegrationService()

	user := models.UserRegister{
		Name:        "Авториз",
		Surname:     "Авторизов",
		PhoneNumber: "+79993334455",
		UserAutorization: models.UserAutorization{
			Email:    "auth@example.com",
			Password: "mypassword",
		},
		Role: models.RoleBasic,
	}
	id, err := svc.UserAdd(context.Background(), user)
	require.NoError(t, err)

	found, err := svc.UserVerification(context.Background(), models.UserAutorization{
		Email:    "auth@example.com",
		Password: "mypassword",
	})
	require.NoError(t, err)
	assert.Equal(t, id, found.Id)
}

func TestIntegration_UserVerification_WrongPassword(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newIntegrationService()

	user := models.UserRegister{
		Name:        "Неверный",
		Surname:     "Пароль",
		PhoneNumber: "+79994445566",
		UserAutorization: models.UserAutorization{
			Email:    "wrongpass@example.com",
			Password: "correctpass",
		},
		Role: models.RoleBasic,
	}
	_, err := svc.UserAdd(context.Background(), user)
	require.NoError(t, err)

	_, err = svc.UserVerification(context.Background(), models.UserAutorization{
		Email:    "wrongpass@example.com",
		Password: "wrongpass",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestIntegration_UserDelete_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newIntegrationService()

	user := models.UserRegister{
		Name:        "Удаляемый",
		Surname:     "Удаляев",
		PhoneNumber: "+79995556677",
		UserAutorization: models.UserAutorization{
			Email:    "delete@example.com",
			Password: "password123",
		},
		Role: models.RoleBasic,
	}
	id, err := svc.UserAdd(context.Background(), user)
	require.NoError(t, err)

	err = svc.UserDelete(context.Background(), id)
	require.NoError(t, err)

	_, err = svc.UserGet(context.Background(), id)
	require.Error(t, err)
}

func TestIntegration_RoleChange_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	svc := newIntegrationService()

	user := models.UserRegister{
		Name:        "Роль",
		Surname:     "Ролев",
		PhoneNumber: "+79996667788",
		UserAutorization: models.UserAutorization{
			Email:    "role@example.com",
			Password: "password123",
		},
		Role: models.RoleBasic,
	}
	id, err := svc.UserAdd(context.Background(), user)
	require.NoError(t, err)

	err = svc.RoleChange(context.Background(), id, models.RoleVerificator)
	require.NoError(t, err)

	found, err := svc.UserGet(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, models.RoleVerificator, found.Role)
}
