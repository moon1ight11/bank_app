package jwt

import (
	"bank_app/internal/api/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testUserID = uuid.MustParse("b6609ddd-95f4-42f0-993e-7f07f3fa1d5b")

func TestCustomFieldsValidate_Success(t *testing.T) {
	claims := &Claims{
		UserId:      &testUserID,
		UserName:    "Иван",
		UserSurname: "Иванов",
		UserEmail:   "ivan@example.com",
	}
	err := claims.CustomFieldsValidate()
	assert.NoError(t, err)
}

func TestCustomFieldsValidate_NilUserID(t *testing.T) {
	claims := &Claims{
		UserId:      nil,
		UserName:    "Иван",
		UserSurname: "Иванов",
		UserEmail:   "ivan@example.com",
	}
	err := claims.CustomFieldsValidate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user id is empty")
}

func TestCustomFieldsValidate_EmptyName(t *testing.T) {
	claims := &Claims{
		UserId:      &testUserID,
		UserName:    "",
		UserSurname: "Иванов",
		UserEmail:   "ivan@example.com",
	}
	err := claims.CustomFieldsValidate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user name")
}

func TestCustomFieldsValidate_EmptySurname(t *testing.T) {
	claims := &Claims{
		UserId:      &testUserID,
		UserName:    "Иван",
		UserSurname: "",
		UserEmail:   "ivan@example.com",
	}
	err := claims.CustomFieldsValidate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user surname")
}

func TestCustomFieldsValidate_EmptyEmail(t *testing.T) {
	claims := &Claims{
		UserId:      &testUserID,
		UserName:    "Иван",
		UserSurname: "Иванов",
		UserEmail:   "",
	}
	err := claims.CustomFieldsValidate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not looks like email")
}

func TestCustomFieldsValidate_InvalidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"no at sign", "ivanexample.com"},
		{"no domain", "ivan@"},
		{"no username", "@example.com"},
		{"with spaces", "ivan @example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := &Claims{
				UserId:      &testUserID,
				UserName:    "Иван",
				UserSurname: "Иванов",
				UserEmail:   tt.email,
			}
			err := claims.CustomFieldsValidate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "not looks like email")
		})
	}
}

func TestGenerateToken_Success(t *testing.T) {
	service := NewJWTService("test-secret", 1*time.Hour)

	token, err := service.GenerateToken(
		testUserID,
		"Иван",
		"Иванов",
		"ivan@example.com",
		models.RoleBasic,
	)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateToken_ParseBack(t *testing.T) {
	service := NewJWTService("test-secret", 1*time.Hour)

	token, err := service.GenerateToken(
		testUserID,
		"Иван",
		"Иванов",
		"ivan@example.com",
		models.RoleBasic,
	)
	require.NoError(t, err)

	claims := &Claims{}
	parsedToken, err := service.ParseToken(token, claims)
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, testUserID, *claims.UserId)
	assert.Equal(t, "Иван", claims.UserName)
	assert.Equal(t, "Иванов", claims.UserSurname)
	assert.Equal(t, "ivan@example.com", claims.UserEmail)
	assert.Equal(t, models.RoleBasic, claims.Role)
}

func TestGenerateToken_InvalidClaims(t *testing.T) {
	service := NewJWTService("test-secret", 1*time.Hour)

	_, err := service.GenerateToken(
		testUserID,
		"Иван",
		"Иванов",
		"",
		models.RoleBasic,
	)
	require.Error(t, err)
}

func TestParseToken_InvalidToken(t *testing.T) {
	service := NewJWTService("test-secret", 1*time.Hour)

	claims := &Claims{}
	_, err := service.ParseToken("invalid.token.string", claims)
	require.Error(t, err)
}

func TestParseToken_WrongSecret(t *testing.T) {
	service1 := NewJWTService("secret-1", 1*time.Hour)
	service2 := NewJWTService("secret-2", 1*time.Hour)

	token, err := service1.GenerateToken(
		testUserID,
		"Иван",
		"Иванов",
		"ivan@example.com",
		models.RoleBasic,
	)
	require.NoError(t, err)

	claims := &Claims{}
	_, err = service2.ParseToken(token, claims)
	require.Error(t, err)
}

func TestParseToken_ExpiredToken(t *testing.T) {
	service := NewJWTService("test-secret", 1*time.Nanosecond)

	token, err := service.GenerateToken(
		testUserID,
		"Иван",
		"Иванов",
		"ivan@example.com",
		models.RoleBasic,
	)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	claims := &Claims{}
	_, err = service.ParseToken(token, claims)
	require.Error(t, err)
}
