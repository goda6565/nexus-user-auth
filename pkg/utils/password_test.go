package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err, "HashPassword should not error")
	assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
}

func TestCheckPassword(t *testing.T) {
	password := "password"
	hashedPassword, _ := HashPassword(password)

	err := CheckPassword(hashedPassword, password)
	assert.NoError(t, err, "CheckPassword should not error")
}

func TestCheckPassword_Invalid(t *testing.T) {
	password := "password"
	hashedPassword, _ := HashPassword(password)

	err := CheckPassword(hashedPassword, "invalid")
	assert.Error(t, err, "CheckPassword should error")
}
