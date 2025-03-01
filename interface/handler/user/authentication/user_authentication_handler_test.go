package authentication_test

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

	"github.com/goda6565/nexus-user-auth/interface/gen"
	. "github.com/goda6565/nexus-user-auth/interface/handler/user/authentication"
)

// --- モックの UserAuthenticationService ---
type mockUserAuthenticationService struct {
	mock.Mock
}

func (m *mockUserAuthenticationService) UserLogin(email, password string) (string, string, error) {
	args := m.Called(email, password)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *mockUserAuthenticationService) UserTokenRefresh(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}

// --- テストスイート ---
type UserAuthenticationHandlerTestSuite struct {
	suite.Suite
	handler     *UserAuthenticationHandler
	mockService *mockUserAuthenticationService
}

func (suite *UserAuthenticationHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockService = new(mockUserAuthenticationService)
	suite.handler = NewUserAuthenticationHandler(suite.mockService)
}

// ----- UserLogin のテスト -----

// 正常系: 正しいJSONを渡し、ログインに成功する場合
func (suite *UserAuthenticationHandlerTestSuite) TestUserLogin_Success() {
	reqBody := gen.UserLoginRequestBody{
		Email:    "test@example.com",
		Password: "password123",
	}
	bodyBytes, err := json.Marshal(reqBody)
	suite.Require().NoError(err)

	accessToken := "access_token_value"
	refreshToken := "refresh_token_value"
	suite.mockService.
		On("UserLogin", reqBody.Email, reqBody.Password).
		Return(accessToken, refreshToken, nil)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.handler.UserLogin(c)

	suite.Equal(http.StatusOK, w.Code)
	var resp gen.LoginResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Require().NoError(err)
	suite.Equal(accessToken, resp.AccessToken)
	suite.Equal(refreshToken, resp.RefreshToken)
	suite.mockService.AssertExpectations(suite.T())
}

// バインドエラー: 不正なJSONの場合
func (suite *UserAuthenticationHandlerTestSuite) TestUserLogin_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.handler.UserLogin(c)

	suite.Equal(http.StatusBadRequest, w.Code)
	var errResp gen.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	suite.Require().NoError(err)
	suite.Equal(http.StatusBadRequest, errResp.Code)
	suite.NotEmpty(errResp.Message)
}

// サービスエラー: ログイン処理でエラーが発生した場合
func (suite *UserAuthenticationHandlerTestSuite) TestUserLogin_ServiceError() {
	reqBody := gen.UserLoginRequestBody{
		Email:    "test@example.com",
		Password: "password123",
	}
	bodyBytes, err := json.Marshal(reqBody)
	suite.Require().NoError(err)

	serviceErr := errors.New("login failed")
	suite.mockService.
		On("UserLogin", reqBody.Email, reqBody.Password).
		Return("", "", serviceErr)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.handler.UserLogin(c)

	suite.Equal(http.StatusInternalServerError, w.Code)
	var errResp gen.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &errResp)
	suite.Require().NoError(err)
	suite.Equal(http.StatusInternalServerError, errResp.Code)
	suite.Equal(serviceErr.Error(), errResp.Message)
	suite.mockService.AssertExpectations(suite.T())
}

// ----- UserTokenRefresh のテスト -----

// 正常系: 正しいJSONを渡し、トークンリフレッシュに成功する場合
func (suite *UserAuthenticationHandlerTestSuite) TestUserTokenRefresh_Success() {
	reqBody := gen.TokenRefreshRequestBody{
		RefreshToken: "old_refresh_token",
	}
	bodyBytes, err := json.Marshal(reqBody)
	suite.Require().NoError(err)

	newAccessToken := "new_access_token_value"
	suite.mockService.
		On("UserTokenRefresh", reqBody.RefreshToken).
		Return(newAccessToken, nil)

	req := httptest.NewRequest(http.MethodPost, "/token/refresh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.handler.UserTokenRefresh(c)

	suite.Equal(http.StatusOK, w.Code)
	var resp gen.TokenRefreshResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Require().NoError(err)
	suite.Equal(newAccessToken, resp.AccessToken)
	suite.mockService.AssertExpectations(suite.T())
}

// バインドエラー: 不正なJSONの場合
func (suite *UserAuthenticationHandlerTestSuite) TestUserTokenRefresh_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/token/refresh", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.handler.UserTokenRefresh(c)

	suite.Equal(http.StatusBadRequest, w.Code)
	var errResp gen.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	suite.Require().NoError(err)
	suite.Equal(http.StatusBadRequest, errResp.Code)
	suite.NotEmpty(errResp.Message)
}

// サービスエラー: トークンリフレッシュ処理でエラーが発生した場合
func (suite *UserAuthenticationHandlerTestSuite) TestUserTokenRefresh_ServiceError() {
	reqBody := gen.TokenRefreshRequestBody{
		RefreshToken: "old_refresh_token",
	}
	bodyBytes, err := json.Marshal(reqBody)
	suite.Require().NoError(err)

	serviceErr := errors.New("refresh failed")
	suite.mockService.
		On("UserTokenRefresh", reqBody.RefreshToken).
		Return("", serviceErr)

	req := httptest.NewRequest(http.MethodPost, "/token/refresh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.handler.UserTokenRefresh(c)

	suite.Equal(http.StatusInternalServerError, w.Code)
	var errResp gen.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &errResp)
	suite.Require().NoError(err)
	suite.Equal(http.StatusInternalServerError, errResp.Code)
	suite.Equal(serviceErr.Error(), errResp.Message)
	suite.mockService.AssertExpectations(suite.T())
}

func TestUserAuthenticationHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserAuthenticationHandlerTestSuite))
}
