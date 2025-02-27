package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserPassword_Valid(t *testing.T) {
	validPassword := "Abcd1234"
	pw, err := NewUserPassword(validPassword)
	assert.NoError(t, err, "有効なパスワードの場合、エラーが返られないこと")
	assert.NotNil(t, pw)
	assert.Equal(t, validPassword, pw.Value())
}

func TestNewUserPassword_TooShort(t *testing.T) {
	shortPassword := "Ab12"
	pw, err := NewUserPassword(shortPassword)
	assert.Error(t, err, "短すぎるパスワードはエラーになること")
	assert.Nil(t, pw)
}

func TestNewUserPassword_MissingDigit(t *testing.T) {
	// 英字はあるが数字が含まれていない
	invalidPassword := "Password"
	pw, err := NewUserPassword(invalidPassword)
	assert.Error(t, err, "数字が含まれていないパスワードはエラーになること")
	assert.Nil(t, pw)
}

func TestNewUserPassword_MissingLetter(t *testing.T) {
	// 数字はあるが英字が含まれていない
	invalidPassword := "12345678"
	pw, err := NewUserPassword(invalidPassword)
	assert.Error(t, err, "英字が含まれていないパスワードはエラーになること")
	assert.Nil(t, pw)
}
