package profile_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/goda6565/nexus-user-auth/domain/user/entity"
	"github.com/goda6565/nexus-user-auth/domain/user/value"
	"github.com/goda6565/nexus-user-auth/interface/gen"
	. "github.com/goda6565/nexus-user-auth/interface/handler/user/profile"
)

// --- モックの UserProfileService ---
type mockUserProfileService struct {
	mock.Mock
}

func (m *mockUserProfileService) UserGet(id string) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserProfileService) UserUpdate(id, avatarURL, username string) (*entity.User, error) {
	args := m.Called(id, avatarURL, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserProfileService) UserDelete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// --- テスト用ユーティリティ ---
// ユーザーIDは entity.NewUser で自動生成されるため、生成後に fakeUser.ObjID().Value() で取得します。
// AvatarURL は値が存在する場合、ChangeAvatarURL 経由で設定します（実装に合わせてください）。
func createFakeUser(emailStr, usernameStr, avatar string) *entity.User {
	email, _ := value.NewUserEmail(emailStr)
	password, _ := value.NewUserPassword("dummy") // テスト用ダミー
	username, _ := value.NewUserUsername(usernameStr)
	fakeUser, _ := entity.NewUser(email, password, username)
	if avatar != "" {
		avatarURL, _ := value.NewUserAvatarURL(avatar)
		fakeUser.ChangeAvatarURL(avatarURL)
	}
	return fakeUser
}

// --- テストスイート ---
type UserProfileHandlerTestSuite struct {
	suite.Suite
	handler     *UserProfileHandler
	mockService *mockUserProfileService
}

func (suite *UserProfileHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockService = new(mockUserProfileService)
	suite.handler = NewUserProfileHandler(suite.mockService)
}

// ----- GetUserProfile のテスト -----

// 正常系: AvatarURL が空の場合
func (suite *UserProfileHandlerTestSuite) TestGetUserProfile_NoAvatar() {
	fakeUser := createFakeUser("test@example.com", "username", "")
	userID := fakeUser.ObjID().Value()

	// モックで、指定されたuidに対して生成したユーザーを返すよう設定
	suite.mockService.
		On("UserGet", userID).
		Return(fakeUser, nil)

	// テスト用のGinコンテキストを作成し、パラメータを設定
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "validated_uid", Value: userID}}

	// ハンドラーの GetUserProfile を呼び出す
	suite.handler.GetUserProfile(c)

	// HTTPステータスコードの検証
	suite.Equal(http.StatusOK, w.Code)

	// レスポンスボディをデコードし、各フィールドが期待通りになっているか検証
	var resp gen.ProfileResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Require().NoError(err)
	suite.Equal(userID, resp.Uid)
	suite.Equal("test@example.com", resp.Email)
	suite.Equal("username", resp.Username)
	// AvatarURL が空の場合は、レスポンスでは nil となるはず
	suite.Nil(resp.AvatarURL)

	// モックの期待通り呼び出されたかを検証
	suite.mockService.AssertExpectations(suite.T())
}

// 正常系: AvatarURL が設定されている場合
func (suite *UserProfileHandlerTestSuite) TestGetUserProfile_WithAvatar() {
	avatar := "https://example.com/avatar.png"
	fakeUser := createFakeUser("user2@example.com", "user2", avatar)
	userID := fakeUser.ObjID().Value()

	suite.mockService.
		On("UserGet", userID).
		Return(fakeUser, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "validated_uid", Value: userID}}

	suite.handler.GetUserProfile(c)

	suite.Equal(http.StatusOK, w.Code)
	var resp gen.ProfileResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Require().NoError(err)
	suite.Equal(userID, resp.Uid)
	suite.Equal("user2@example.com", resp.Email)
	suite.Equal("user2", resp.Username)
	suite.NotNil(resp.AvatarURL)
	suite.Equal(avatar, *resp.AvatarURL)

	suite.mockService.AssertExpectations(suite.T())
}

// エラー系: サービス側でエラー発生
func (suite *UserProfileHandlerTestSuite) TestGetUserProfile_ServiceError() {
	userID := "789" // この値はモックの期待値として利用します
	serviceErr := errors.New("user not found")
	suite.mockService.
		On("UserGet", userID).
		Return(nil, serviceErr)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "validated_uid", Value: userID}}

	suite.handler.GetUserProfile(c)

	suite.Equal(http.StatusInternalServerError, w.Code)
	var errResp gen.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	suite.Require().NoError(err)
	suite.Equal(http.StatusInternalServerError, errResp.Code)
	suite.Equal(serviceErr.Error(), errResp.Message)

	suite.mockService.AssertExpectations(suite.T())
}

// ----- UpdateUserProfile のテスト -----

// バインドエラー: 不正な JSON を渡した場合
func (suite *UserProfileHandlerTestSuite) TestUpdateUserProfile_InvalidJSON() {
	userID := "123"
	req := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "validated_uid", Value: userID}}
	c.Request = req

	suite.handler.UpdateUserProfile(c)

	suite.Equal(http.StatusBadRequest, w.Code)
	var errResp gen.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	suite.Require().NoError(err)
	suite.Equal(http.StatusBadRequest, errResp.Code)
	suite.NotEmpty(errResp.Message)
}

// エラー系: サービス側でエラー発生
func (suite *UserProfileHandlerTestSuite) TestUpdateUserProfile_ServiceError() {
	userID := "123"
	reqBody := gen.UserProfileUpdateRequestBody{
		AvatarURL: ptr("https://example.com/newavatar.png"),
		Username:  "newusername",
	}
	bodyBytes, err := json.Marshal(reqBody)
	suite.Require().NoError(err)

	req := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "validated_uid", Value: userID}}
	c.Request = req

	serviceErr := errors.New("update failed")
	suite.mockService.
		On("UserUpdate", userID, *reqBody.AvatarURL, reqBody.Username).
		Return(nil, serviceErr)

	suite.handler.UpdateUserProfile(c)

	suite.Equal(http.StatusInternalServerError, w.Code)
	var errResp gen.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &errResp)
	suite.Require().NoError(err)
	suite.Equal(http.StatusInternalServerError, errResp.Code)
	suite.Equal(serviceErr.Error(), errResp.Message)

	suite.mockService.AssertExpectations(suite.T())
}

// 正常系: ユーザー情報更新成功
func (suite *UserProfileHandlerTestSuite) TestUpdateUserProfile_Success() {
	// サービス側が返す fakeUser の uid を利用する
	fakeUser := createFakeUser("test@example.com", "newusername", "https://example.com/newavatar.png")
	userID := fakeUser.ObjID().Value()

	reqBody := gen.UserProfileUpdateRequestBody{
		AvatarURL: ptr("https://example.com/newavatar.png"),
		Username:  "newusername",
	}
	bodyBytes, err := json.Marshal(reqBody)
	suite.Require().NoError(err)

	req := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// コンテキストの uid として、fakeUser の uid を設定
	c.Params = []gin.Param{{Key: "validated_uid", Value: userID}}
	c.Request = req

	suite.mockService.
		On("UserUpdate", userID, *reqBody.AvatarURL, reqBody.Username).
		Return(fakeUser, nil)

	suite.handler.UpdateUserProfile(c)

	suite.Equal(http.StatusOK, w.Code)
	var resp gen.ProfileResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Require().NoError(err)
	suite.Equal(userID, resp.Uid)
	suite.Equal("test@example.com", resp.Email)
	suite.Equal("newusername", resp.Username)
	suite.NotNil(resp.AvatarURL)
	suite.Equal("https://example.com/newavatar.png", *resp.AvatarURL)

	suite.mockService.AssertExpectations(suite.T())
}

// ----- DeleteUserProfile のテスト -----

// エラー系: サービス側でエラー発生
func (suite *UserProfileHandlerTestSuite) TestDeleteUserProfile_ServiceError() {
	userID := "123"
	serviceErr := errors.New("delete failed")
	suite.mockService.
		On("UserDelete", userID).
		Return(serviceErr)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "validated_uid", Value: userID}}

	suite.handler.DeleteUserProfile(c)

	suite.Equal(http.StatusInternalServerError, w.Code)
	var errResp gen.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	suite.Require().NoError(err)
	suite.Equal(http.StatusInternalServerError, errResp.Code)
	suite.Equal(serviceErr.Error(), errResp.Message)

	suite.mockService.AssertExpectations(suite.T())
}

// 正常系: ユーザー削除成功
func (suite *UserProfileHandlerTestSuite) TestDeleteUserProfile_Success() {
	userID := "123"
	suite.mockService.
		On("UserDelete", userID).
		Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "validated_uid", Value: userID}}

	suite.handler.DeleteUserProfile(c)

	suite.Equal(http.StatusOK, w.Code)
	// レスポンスボディが nil または空であることを検証
	suite.Empty(w.Body)

	suite.mockService.AssertExpectations(suite.T())
}

// --- ヘルパー関数 ---
func ptr(s string) *string {
	return &s
}

func TestUserProfileHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserProfileHandlerTestSuite))
}
