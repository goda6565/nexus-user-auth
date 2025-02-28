package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserAvatarURL_Valid(t *testing.T) {
	validURL := "https://avatar.example.com/123.png"
	avatar, err := NewUserAvatarURL(validURL)
	assert.NoError(t, err, "正しいURLの場合、エラーが返られないこと")
	assert.NotNil(t, avatar)
	assert.Equal(t, validURL, avatar.Value())
}

func TestNewUserAvatarURL_Invalid(t *testing.T) {
	invalidURL := "not-a-url"
	avatar, err := NewUserAvatarURL(invalidURL)
	assert.Error(t, err, "不正なURLの場合、エラーが返ること")
	assert.Nil(t, avatar, "不正なURLの場合、nil が返ること")
}
