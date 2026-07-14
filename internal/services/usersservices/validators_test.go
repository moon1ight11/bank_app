package usersservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"simple", "ivan@example.com", true},
		{"with dots", "ivan.petrov@example.com", true},
		{"with plus", "ivan+tag@example.com", true},
		{"with numbers", "ivan123@example.com", true},
		{"subdomain", "ivan@mail.example.com", true},
		{"no at sign", "ivanexample.com", false},
		{"no domain", "ivan@", false},
		{"no username", "@example.com", false},
		{"empty string", "", false},
		{"with spaces", "ivan @example.com", false},
		{"double at", "ivan@@example.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidEmail(tt.email)
			assert.Equal(t, tt.valid, result, "email: %s", tt.email)
		})
	}
}

func TestIsValidPhoneNumber(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		valid bool
	}{
		{"russian +7", "+79991234567", true},
		{"russian 8", "89991234567", true},
		{"russian with spaces", "+7 999 123 45 67", true},
		{"russian with dashes", "+7-999-123-45-67", true},
		{"russian with parens", "+7(999)1234567", true},
		{"simple digits", "1234567890", true},
		{"11 digits", "12345678901", true},
		{"empty", "", false},
		{"letters", "abc", false},
		{"too short", "123", false},
		{"special chars", "+7(999)123-45-67!", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidPhoneNumber(tt.phone)
			assert.Equal(t, tt.valid, result, "phone: %s", tt.phone)
		})
	}
}

func TestIsValidTimezone(t *testing.T) {
	tests := []struct {
		name  string
		tz    string
		valid bool
	}{
		{"UTC+3", "UTC+03:00", true},
		{"UTC-5", "UTC-05:00", true},
		{"UTC+0", "UTC+00:00", true},
		{"UTC+14", "UTC+14:00", true},
		{"UTC-12", "UTC-12:00", true},
		{"UTC+3:30", "UTC+03:30", true},
		{"no prefix", "+03:00", false},
		{"no colon", "UTC+0300", false},
		{"out of range high", "UTC+25:00", false},
		{"out of range low", "UTC-13:00", false},
		{"out of range minutes", "UTC+03:60", false},
		{"empty", "", false},
		{"random string", "Europe/Moscow", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidTimezone(tt.tz)
			assert.Equal(t, tt.valid, result, "tz: %s", tt.tz)
		})
	}
}
