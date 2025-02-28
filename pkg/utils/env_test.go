package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvDefault_VarNotSet(t *testing.T) {
	// テスト対象の環境変数名を一意なものにする
	key := "TEST_GET_ENV_DEFAULT_NOT_SET"
	// 環境変数が設定されていないことを保証
	err := os.Unsetenv(key)
	assert.NoError(t, err, "Unsetenv should not return an error")

	defaultVal := "default"
	val := GetEnvDefault(key, defaultVal)
	assert.Equal(t, defaultVal, val, "When the env variable is not set, the default value should be returned")
}

func TestGetEnvDefault_VarSet(t *testing.T) {
	key := "TEST_GET_ENV_DEFAULT_SET"
	expectedVal := "myvalue"
	err := os.Setenv(key, expectedVal)
	assert.NoError(t, err, "Setting env variable should not return an error")
	// テスト終了後にクリーンアップする
	defer func() {
		err := os.Unsetenv(key)
		assert.NoError(t, err, "Unsetenv should not return an error")
	}()

	defaultVal := "default"
	val := GetEnvDefault(key, defaultVal)
	assert.Equal(t, expectedVal, val, "When the env variable is set, its value should be returned")
}
