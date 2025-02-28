package utils

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// setupEnv は環境変数 JWT_SECRET_KEY のセットアップと後片付けを行います。
func setupEnv(t *testing.T) {
	t.Helper()
	err := os.Setenv("JWT_SECRET_KEY", "mysecret")
	assert.NoError(t, err, "Setting JWT_SECRET_KEY should not error")

	t.Cleanup(func() {
		err := os.Unsetenv("JWT_SECRET_KEY")
		assert.NoError(t, err, "Unsetenv should not return an error")
	})
}

// TestGenerateTokens は、アクセストークンとリフレッシュトークンの生成テスト
func TestGenerateTokens(t *testing.T) {
	setupEnv(t)

	objID := "123"

	accessToken, refreshToken, err := GenerateTokens(objID)
	assert.NoError(t, err, "トークン生成中にエラーが発生してはいけない")
	assert.NotEmpty(t, accessToken, "アクセストークンが生成されるべき")
	assert.NotEmpty(t, refreshToken, "リフレッシュトークンが生成されるべき")
}

// TestValidateAccessToken は、アクセストークンの検証テスト
func TestValidateAccessToken(t *testing.T) {
	setupEnv(t)

	objID := "123"

	accessToken, _, err := GenerateTokens(objID)
	assert.NoError(t, err)

	claims, err := ValidateToken(accessToken)
	assert.NoError(t, err, "有効なトークンの検証中にエラーが発生してはいけない")
	assert.Equal(t, objID, claims.ObjID, "検証結果のユーザーIDが一致する")
}

// TestValidateInvalidAccessToken は、無効なアクセストークンの検証テスト
func TestValidateInvalidAccessToken(t *testing.T) {
	setupEnv(t)

	invalidToken := "this.is.an.invalid.token"
	claims, err := ValidateToken(invalidToken)
	assert.Error(t, err, "無効なトークンの場合はエラーが返る")
	assert.Nil(t, claims, "無効なトークンの場合、claims は nil でなければならない")
}

// TestAccessTokenExpiration は、期限切れのアクセストークンの検証テスト
func TestAccessTokenExpiration(t *testing.T) {
	setupEnv(t)

	// オリジナルの時間取得関数を保持
	originalTimeNow := timeNowFunc
	defer func() { timeNowFunc = originalTimeNow }()

	// 25時間前の時刻を設定（アクセストークンの期限は 24時間なので、1時間前に期限切れ）
	timeNowFunc = func() time.Time { return time.Now().Add(-25 * time.Hour) }

	accessToken, _, err := GenerateTokens("123")
	assert.NoError(t, err)

	claims, err := ValidateToken(accessToken)
	assert.Error(t, err, "期限切れのトークンは検証時にエラーとなる")
	assert.Nil(t, claims)
}

// TestValidateRefreshToken は、リフレッシュトークンの検証テスト
func TestValidateRefreshToken(t *testing.T) {
	setupEnv(t)

	objID := "123"

	_, refreshToken, err := GenerateTokens(objID)
	assert.NoError(t, err)

	claims, err := ValidateRefreshToken(refreshToken)
	assert.NoError(t, err, "有効なリフレッシュトークンの検証中にエラーが発生してはいけない")
	assert.Equal(t, objID, claims.ObjID, "検証結果のユーザーIDが一致する")
}

// TestValidateInvalidRefreshToken は、無効なリフレッシュトークンの検証テスト
func TestValidateInvalidRefreshToken(t *testing.T) {
	setupEnv(t)

	invalidToken := "this.is.an.invalid.refresh.token"
	claims, err := ValidateRefreshToken(invalidToken)
	assert.Error(t, err, "無効なリフレッシュトークンの場合はエラーが返る")
	assert.Nil(t, claims, "無効なトークンの場合、claims は nil でなければならない")
}

// TestRefreshTokenExpiration は、期限切れのリフレッシュトークンの検証テスト
func TestRefreshTokenExpiration(t *testing.T) {
	setupEnv(t)

	originalTimeNow := timeNowFunc
	defer func() { timeNowFunc = originalTimeNow }()

	// 8日前の時刻を設定（リフレッシュトークンの期限は 7日間なので、1日超過）
	timeNowFunc = func() time.Time { return time.Now().Add(-8 * 24 * time.Hour) }

	_, refreshToken, err := GenerateTokens("123")
	assert.NoError(t, err)

	claims, err := ValidateRefreshToken(refreshToken)
	assert.Error(t, err, "期限切れのリフレッシュトークンは検証時にエラーとなる")
	assert.Nil(t, claims)
}

// TestRefreshAccessToken は、リフレッシュトークンを使ったアクセストークン更新テスト
func TestRefreshAccessToken(t *testing.T) {
	setupEnv(t)

	objID := "123"

	_, refreshToken, err := GenerateTokens(objID)
	assert.NoError(t, err)

	newAccessToken, err := RefreshAccessToken(refreshToken)
	assert.NoError(t, err, "リフレッシュトークンでアクセストークンを更新できるべき")
	assert.NotEmpty(t, newAccessToken, "新しいアクセストークンが生成されるべき")
}

// TestRefreshAccessTokenWithInvalidToken は、無効なリフレッシュトークンでのアクセストークン更新テスト
func TestRefreshAccessTokenWithInvalidToken(t *testing.T) {
	setupEnv(t)

	invalidToken := "this.is.an.invalid.refresh.token"

	newAccessToken, err := RefreshAccessToken(invalidToken)
	assert.Error(t, err, "無効なリフレッシュトークンではアクセストークンを更新できない")
	assert.Empty(t, newAccessToken, "エラー時には新しいアクセストークンは生成されるべきではない")
}
