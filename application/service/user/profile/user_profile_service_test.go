package profile_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/goda6565/nexus-user-auth/application/service/user/profile"
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

// UserProfileServiceTestSuite は UserProfileService のテストスイート
type UserProfileServiceTestSuite struct {
	suite.Suite
	mockRepo *mockUserRepository
	service  profile.UserProfileService
}

func TestUserProfileServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserProfileServiceTestSuite))
}

// SetupTest は各テスト前に実行されるセットアップ処理
func (suite *UserProfileServiceTestSuite) SetupTest() {
	suite.mockRepo = NewMockUserRepository()
	suite.service = profile.NewUserProfileService(suite.mockRepo)
}

// TestUserUpdate_Success は、ユーザー名とアバターURLの更新が成功するケース
func (suite *UserProfileServiceTestSuite) TestUserUpdate_Success() {
	objID := "user-123"
	username := "new_username"
	avatarURL := "https://example.com/avatar.png"

	// モックの設定
	oldUser := &entity.User{
		// モックのエンティティ作成
	}
	newUsername, _ := value.NewUserUsername(username)
	newAvatarURL, _ := value.NewUserAvatarURL(avatarURL)

	updatedUser := *oldUser
	updatedUser.ChangeUsername(newUsername)
	updatedUser.ChangeAvatarURL(newAvatarURL)

	suite.mockRepo.On("GetUserByObjID", objID).Return(oldUser, nil)
	suite.mockRepo.On("UpdateUser", mock.Anything).Return(&updatedUser, nil)

	// 実行
	result, err := suite.service.UserUpdate(objID, username, avatarURL)

	// 検証
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), username, result.Username().Value())
	assert.Equal(suite.T(), avatarURL, result.AvatarURL().Value())

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestUserUpdate_EmptyAvatarURL は、アバターURLが空の場合に変更されないことをテスト
func (suite *UserProfileServiceTestSuite) TestUserUpdate_EmptyAvatarURL() {
	objID := "user-123"
	username := "new_username"
	avatarURL := ""

	oldUser := &entity.User{
		// モックのエンティティ作成
	}
	newUsername, _ := value.NewUserUsername(username)

	updatedUser := *oldUser
	updatedUser.ChangeUsername(newUsername)

	suite.mockRepo.On("GetUserByObjID", objID).Return(oldUser, nil)
	suite.mockRepo.On("UpdateUser", mock.Anything).Return(&updatedUser, nil)

	// 実行
	result, err := suite.service.UserUpdate(objID, username, avatarURL)

	// 検証
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), username, result.Username().Value())
	assert.Equal(suite.T(), oldUser.AvatarURL(), result.AvatarURL()) // 変更されていない

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestUserUpdate_GetUserByObjID_Error は、ユーザー取得時にエラーが発生するケース
func (suite *UserProfileServiceTestSuite) TestUserUpdate_GetUserByObjID_Error() {
	objID := "user-123"
	username := "new_username"
	avatarURL := "https://example.com/avatar.png"

	suite.mockRepo.On("GetUserByObjID", objID).Return(nil, errs.NewServiceError("user not found"))

	// 実行
	result, err := suite.service.UserUpdate(objID, username, avatarURL)

	// 検証
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestUserUpdate_UpdateUser_Error は、ユーザー更新時にエラーが発生するケース
func (suite *UserProfileServiceTestSuite) TestUserUpdate_UpdateUser_Error() {
	objID := "user-123"
	username := "new_username"
	avatarURL := "https://example.com/avatar.png"

	oldUser := &entity.User{
		// モックのエンティティ作成
	}
	newUsername, _ := value.NewUserUsername(username)
	newAvatarURL, _ := value.NewUserAvatarURL(avatarURL)

	updatedUser := *oldUser
	updatedUser.ChangeUsername(newUsername)
	updatedUser.ChangeAvatarURL(newAvatarURL)

	suite.mockRepo.On("GetUserByObjID", objID).Return(oldUser, nil)
	suite.mockRepo.On("UpdateUser", mock.Anything).Return(nil, errs.NewServiceError("failed to update user"))

	// 実行
	result, err := suite.service.UserUpdate(objID, username, avatarURL)

	// 検証
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestUserUpdate_InvalidUsername は、無効なユーザー名のケース
func (suite *UserProfileServiceTestSuite) TestUserUpdate_InvalidUsername() {
	objID := "user-123"
	username := "" // 無効なユーザー名
	avatarURL := "https://example.com/avatar.png"

	oldUser := &entity.User{
		// モックのエンティティ作成
	}

	suite.mockRepo.On("GetUserByObjID", objID).Return(oldUser, nil)

	// 実行
	result, err := suite.service.UserUpdate(objID, username, avatarURL)

	// 検証
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)

	suite.mockRepo.AssertExpectations(suite.T())
}
