package authentication

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/goda6565/nexus-user-auth/application/service/user/authentication"
	"github.com/goda6565/nexus-user-auth/interface/gen"
)

type UserAuthenticationHandler struct {
	userAuthenticationService authentication.UserAuthenticationService
}

func NewUserAuthenticationHandler(userAuthenticationService authentication.UserAuthenticationService) *UserAuthenticationHandler {
	return &UserAuthenticationHandler{
		userAuthenticationService: userAuthenticationService,
	}
}

// UserLogin: ユーザーログイン
func (h *UserAuthenticationHandler) UserLogin(c *gin.Context) {
	var req gen.UserLoginRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gen.ErrorResponse{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}

	accessToken, refreshToken, err := h.userAuthenticationService.UserLogin(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gen.ErrorResponse{Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gen.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// UserTokenRefresh: トークンリフレッシュ
func (h *UserAuthenticationHandler) UserTokenRefresh(c *gin.Context) {
	var req gen.TokenRefreshRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gen.ErrorResponse{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}

	accessToken, err := h.userAuthenticationService.UserTokenRefresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gen.ErrorResponse{Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gen.TokenRefreshResponse{
		AccessToken: accessToken,
	})
}
