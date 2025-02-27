package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserEmail_Valid(t *testing.T) {
	validEmail := "user@example.com"
	email, err := NewUserEmail(validEmail)
	assert.NoError(t, err, "正しいメールアドレスの場合、エラーが返られないこと")
	assert.NotNil(t, email)
	assert.Equal(t, validEmail, email.Value())
}

func TestNewUserEmail_Empty(t *testing.T) {
	email, err := NewUserEmail("")
	assert.Error(t, err, "空文字の場合、エラーが返ること")
	assert.Nil(t, email)
}

func TestNewUserEmail_InvalidFormat(t *testing.T) {
	invalidEmail := "not-an-email"
	email, err := NewUserEmail(invalidEmail)
	assert.Error(t, err, "不正な形式の場合、エラーが返ること")
	assert.Nil(t, email)
}

func TestUserEmail_Equals(t *testing.T) {
	emailStr := "user@example.com"
	email1, err1 := NewUserEmail(emailStr)
	email2, err2 := NewUserEmail(emailStr)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.True(t, email1.Equals(email2), "同じメールアドレスの場合 Equals は true を返す")
}
