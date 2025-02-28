package authentication_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/goda6565/nexus-user-auth/application/service/user/authentication"
	"github.com/goda6565/nexus-user-auth/domain/user/entity"
	"github.com/goda6565/nexus-user-auth/domain/user/value"
	"github.com/goda6565/nexus-user-auth/errs"
	"github.com/goda6565/nexus-user-auth/pkg/utils"
)

// --- モックリポジトリ ---

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) CreateUser(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepository) GetUserByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepository) GetUserByObjID(objID string) (*entity.User, error) {
	args := m.Called(objID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepository) DeleteUser(objID string) error {
	args := m.Called(objID)
	return args.Error(0)
}

// --- テストスイート ---

type AuthServiceTestSuite struct {
	suite.Suite
	mockRepo *mockUserRepository
	authServ authentication.UserAuthenticationService
	testUser *entity.User
}

// SetupSuite: JWT_SECRET_KEY のセットアップ
func (suite *AuthServiceTestSuite) SetupSuite() {
	err := os.Setenv("JWT_SECRET_KEY", "mysecret")
	suite.Require().NoError(err, "環境変数の設定に失敗してはいけない")
}

// TearDownSuite: 後片付け
func (suite *AuthServiceTestSuite) TearDownSuite() {
	err := os.Unsetenv("JWT_SECRET_KEY")
	suite.Require().NoError(err, "環境変数の後片付けに失敗してはいけない")
}

// SetupTest: 各テスト前のセットアップ
func (suite *AuthServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mockUserRepository)
	suite.authServ = authentication.NewUserAuthenticationService(suite.mockRepo)

	// テスト用ユーザー作成
	emailVal, _ := value.NewUserEmail("test@example.com")
	usernameVal, _ := value.NewUserUsername("testuser")
	// パスワードは bcrypt ハッシュ済み文字列として保存している前提
	hashedPwd, _ := utils.HashPassword("correct-password")
	passwordVal := value.FromHashed(hashedPwd)

	testUser, err := entity.NewUser(emailVal, passwordVal, usernameVal)
	suite.Require().NoError(err)
	suite.testUser = testUser
}

// --- テストケース ---

// UserLogin 成功パターン
func (suite *AuthServiceTestSuite) TestUserLogin_Success() {
	email := "test@example.com"
	password := "correct-password"

	// モック: GetUserByEmail が testUser を返す
	suite.mockRepo.On("GetUserByEmail", email).Return(suite.testUser, nil)

	accessToken, refreshToken, err := suite.authServ.UserLogin(email, password)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), accessToken)
	assert.NotEmpty(suite.T(), refreshToken)

	suite.mockRepo.AssertExpectations(suite.T())
}

// UserLogin: 存在しないメールアドレスの場合
func (suite *AuthServiceTestSuite) TestUserLogin_InvalidEmail() {
	email := "nonexistent@example.com"
	password := "any-password"

	suite.mockRepo.On("GetUserByEmail", email).Return(nil, errs.NewServiceError("user not found"))

	accessToken, refreshToken, err := suite.authServ.UserLogin(email, password)
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), accessToken)
	assert.Empty(suite.T(), refreshToken)

	suite.mockRepo.AssertExpectations(suite.T())
}

// UserLogin: パスワード不一致の場合
func (suite *AuthServiceTestSuite) TestUserLogin_InvalidPassword() {
	email := "test@example.com"
	wrongPassword := "wrong-password"

	suite.mockRepo.On("GetUserByEmail", email).Return(suite.testUser, nil)

	accessToken, refreshToken, err := suite.authServ.UserLogin(email, wrongPassword)
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), accessToken)
	assert.Empty(suite.T(), refreshToken)

	suite.mockRepo.AssertExpectations(suite.T())
}

// UserTokenRefresh: 成功パターン
func (suite *AuthServiceTestSuite) TestUserTokenRefresh_Success() {
	// まず、トークン生成でリフレッシュトークンを取得
	objID := suite.testUser.ObjID().Value()
	_, refreshToken, err := utils.GenerateTokens(objID)
	suite.Require().NoError(err)
	suite.Require().NotEmpty(refreshToken)

	// 実際に新しいアクセストークンを発行
	newAccessToken, err := suite.authServ.UserTokenRefresh(refreshToken)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), newAccessToken)
}

// UserTokenRefresh: 無効なリフレッシュトークンの場合
func (suite *AuthServiceTestSuite) TestUserTokenRefresh_InvalidToken() {
	invalidToken := "invalid.refresh.token"
	newAccessToken, err := suite.authServ.UserTokenRefresh(invalidToken)
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), newAccessToken)
}

// --- Suite の実行 ---
func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
