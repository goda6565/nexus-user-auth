package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserUsername_Valid(t *testing.T) {
	validUsername := "john_doe"
	username, err := NewUserUsername(validUsername)
	assert.NoError(t, err, "正しいユーザー名の場合、エラーが返られないこと")
	assert.NotNil(t, username)
	assert.Equal(t, validUsername, username.Value())
}

func TestNewUserUsername_TooShort(t *testing.T) {
	invalidUsername := "ab" // 3文字未満
	username, err := NewUserUsername(invalidUsername)
	assert.Error(t, err, "短すぎるユーザー名はエラーになること")
	assert.Nil(t, username)
}

func TestNewUserUsername_TooLong(t *testing.T) {
	// 50文字を超えるユーザー名
	longName := ""
	for i := 0; i < 51; i++ {
		longName += "a"
	}
	username, err := NewUserUsername(longName)
	assert.Error(t, err, "長すぎるユーザー名はエラーになること")
	assert.Nil(t, username)
}

func TestNewUserUsername_LengthBoundary(t *testing.T) {
	// 3文字、50文字の境界値テスト
	minName := "abc" // 3文字
	username, err := NewUserUsername(minName)
	assert.NoError(t, err)
	assert.NotNil(t, username)
	assert.Equal(t, minName, username.Value())

	// 50文字のユーザー名
	longName := ""
	for i := 0; i < 50; i++ {
		longName += "a"
	}
	username, err = NewUserUsername(longName)
	assert.NoError(t, err)
	assert.NotNil(t, username)
	assert.Equal(t, longName, username.Value())
}
