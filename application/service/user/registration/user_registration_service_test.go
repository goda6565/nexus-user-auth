package registration_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/goda6565/nexus-user-auth/application/service/user/registration"
	"github.com/goda6565/nexus-user-auth/domain/user/entity"
	"github.com/goda6565/nexus-user-auth/domain/user/value"
	"github.com/goda6565/nexus-user-auth/errs"
)

// モックリポジトリ（UserRepository のテスト用実装）
type mockUserRepository struct {
	mock.Mock
}

func NewMockUserRepository() *mockUserRepository {
	return &mockUserRepository{}
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

// テストスイート
type UserServiceTestSuite struct {
	suite.Suite
	userService registration.UserRegistrationService
	repo        *mockUserRepository
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.repo = NewMockUserRepository()
	suite.userService = registration.NewUserRegistrationService(suite.repo)
}

// 正常な登録処理のテスト
func (suite *UserServiceTestSuite) TestRegister_Success() {
	email := "test@example.com"
	password := "Password123!"
	username := "testuser"

	// 値オブジェクト生成
	emailValue, err := value.NewUserEmail(email)
	suite.NoError(err)
	passwordValue, err := value.NewUserPassword(password)
	suite.NoError(err)
	usernameValue, err := value.NewUserUsername(username)
	suite.NoError(err)

	userEntity, err := entity.NewUser(emailValue, passwordValue, usernameValue)
	suite.NoError(err)

	suite.repo.
		On("CreateUser", mock.AnythingOfType("*entity.User")).
		Return(userEntity, nil)

	createdUser, err := suite.userService.UserRegister(email, password, username)
	suite.NoError(err)
	suite.NotNil(createdUser)
	suite.Equal(email, createdUser.Email().Value())
}

// 不正なメールの場合のテスト
func (suite *UserServiceTestSuite) TestRegister_InvalidEmail() {
	createdUser, err := suite.userService.UserRegister("invalid-email", "Password123!", "testuser")
	suite.Error(err)
	suite.Nil(createdUser)
}

// 不正なパスワードの場合のテスト
func (suite *UserServiceTestSuite) TestRegister_InvalidPassword() {
	createdUser, err := suite.userService.UserRegister("test@example.com", "short", "testuser")
	suite.Error(err)
	suite.Nil(createdUser)
}

// 不正なユーザー名の場合のテスト
func (suite *UserServiceTestSuite) TestRegister_InvalidUsername() {
	createdUser, err := suite.userService.UserRegister("test@example.com", "Password123!", "")
	suite.Error(err)
	suite.Nil(createdUser)
}

// リポジトリエラー発生時のテスト
func (suite *UserServiceTestSuite) TestRegister_RepositoryError() {
	suite.repo.
		On("CreateUser", mock.AnythingOfType("*entity.User")).
		Return(nil, errs.NewServiceError("repository error"))

	createdUser, err := suite.userService.UserRegister("test@example.com", "Password123!", "testuser")
	suite.Error(err)
	suite.Nil(createdUser)
	suite.Contains(err.Error(), "repository")
}
