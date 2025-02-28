package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserPassword_Valid(t *testing.T) {
	validPassword := "Abcd1234"
	pw, err := NewUserPassword(validPassword)
	assert.NoError(t, err, "有効なパスワードの場合、エラーが返られないこと")
	assert.NotNil(t, pw, "UserPassword オブジェクトが生成されること")

	// ハッシュ値はプレーンなパスワードとは一致しないはず
	assert.NotEqual(t, validPassword, pw.Value(), "ハッシュ値はプレーンなパスワードと異なるはず")

	// 正しいパスワードで検証できること
	assert.True(t, pw.Verify(validPassword), "Verify メソッドは正しいパスワードに対して true を返す")
}

func TestNewUserPassword_TooShort(t *testing.T) {
	shortPassword := "Ab12"
	pw, err := NewUserPassword(shortPassword)
	assert.Error(t, err, "短すぎるパスワードはエラーになること")
	assert.Nil(t, pw, "UserPassword オブジェクトは生成されない")
}

func TestNewUserPassword_MissingDigit(t *testing.T) {
	// 英字はあるが数字が含まれていない
	invalidPassword := "Password"
	pw, err := NewUserPassword(invalidPassword)
	assert.Error(t, err, "数字が含まれていないパスワードはエラーになること")
	assert.Nil(t, pw, "UserPassword オブジェクトは生成されない")
}

func TestNewUserPassword_MissingLetter(t *testing.T) {
	// 数字はあるが英字が含まれていない
	invalidPassword := "12345678"
	pw, err := NewUserPassword(invalidPassword)
	assert.Error(t, err, "英字が含まれていないパスワードはエラーになること")
	assert.Nil(t, pw, "UserPassword オブジェクトは生成されない")
}
