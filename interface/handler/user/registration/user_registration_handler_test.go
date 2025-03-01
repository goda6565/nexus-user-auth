package registration_test

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
	. "github.com/goda6565/nexus-user-auth/interface/handler/user/registration"
)

// --- モックの UserRegistrationService ---
type mockUserRegistrationService struct {
	mock.Mock
}

// サービスメソッド: UserRegister(email, password, username) (*entity.User, error)
func (m *mockUserRegistrationService) UserRegister(email, password, username string) (*entity.User, error) {
	args := m.Called(email, password, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

// --- テストスイート ---
type UserRegistrationHandlerTestSuite struct {
	suite.Suite
	handler     *UserRegistrationHandler
	mockService *mockUserRegistrationService
}

func (suite *UserRegistrationHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode) // テストモードに設定
	suite.mockService = new(mockUserRegistrationService)
	suite.handler = NewUserRegistrationHandler(suite.mockService)
}

// 正常系テスト: 正しいJSONを渡し、サービス側で正常にユーザー登録が完了した場合
func (suite *UserRegistrationHandlerTestSuite) TestUserRegisterSuccess() {
	reqBody := gen.UserRegisterRequestBody{
		Email:    "test@example.com",
		Password: "password",
		Username: "username",
	}
	bodyBytes, err := json.Marshal(reqBody)
	suite.Require().NoError(err)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// モックサービスに正常系の戻り値を返すよう設定
	email, err := value.NewUserEmail("test@example.com")
	suite.Require().NoError(err)
	password, err := value.NewUserPassword("password30")
	suite.Require().NoError(err)
	username, err := value.NewUserUsername("username")
	suite.Require().NoError(err)
	fakeUser, err := entity.NewUser(email, password, username)
	suite.Require().NoError(err)

	suite.mockService.
		On("UserRegister", reqBody.Email, reqBody.Password, reqBody.Username).
		Return(fakeUser, nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.handler.UserRegister(c)

	suite.Equal(http.StatusOK, w.Code)

	var resp gen.RegisterResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Require().NoError(err)
	suite.NotEmpty(resp.Uid)
	suite.Equal("test@example.com", resp.Email)
	suite.Equal("username", resp.Username)

	suite.mockService.AssertExpectations(suite.T())
}

// テスト: リクエストのJSONが不正な場合（バインドエラー）
func (suite *UserRegistrationHandlerTestSuite) TestUserRegisterInvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.handler.UserRegister(c)

	suite.Equal(http.StatusBadRequest, w.Code)

	var errResp gen.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	suite.Require().NoError(err)
	suite.Equal(http.StatusBadRequest, errResp.Code)
	suite.NotEmpty(errResp.Message)
}

// テスト: サービス側でエラーが発生した場合
func (suite *UserRegistrationHandlerTestSuite) TestUserRegisterServiceError() {
	reqBody := gen.UserRegisterRequestBody{
		Email:    "test@example.com",
		Password: "password",
		Username: "username",
	}
	bodyBytes, err := json.Marshal(reqBody)
	suite.Require().NoError(err)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	serviceErr := errors.New("registration failed")
	suite.mockService.
		On("UserRegister", reqBody.Email, reqBody.Password, reqBody.Username).
		Return(nil, serviceErr)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.handler.UserRegister(c)

	suite.Equal(http.StatusInternalServerError, w.Code)

	var errResp gen.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &errResp)
	suite.Require().NoError(err)
	suite.Equal(http.StatusInternalServerError, errResp.Code)
	suite.Equal(serviceErr.Error(), errResp.Message)

	suite.mockService.AssertExpectations(suite.T())
}

func TestUserRegistrationHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserRegistrationHandlerTestSuite))
}
