package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/goda6565/nexus-user-auth/application/service/user/registration"
	"github.com/goda6565/nexus-user-auth/domain/user/entity"
	"github.com/goda6565/nexus-user-auth/interface/gen"
)

type UserRegistrationHandler struct {
	userRegistrationService registration.UserRegistrationService
}

func NewUserRegistrationHandler(userRegistrationService registration.UserRegistrationService) *UserRegistrationHandler {
	return &UserRegistrationHandler{
		userRegistrationService: userRegistrationService,
	}
}

func userToUserRegistrationResponse(user *entity.User) *gen.RegisterResponse {
	return &gen.RegisterResponse{
		Uid:      user.ObjID().Value(),
		Email:    user.Email().Value(),
		Username: user.Username().Value(),
	}
}

// UserRegister: ユーザー登録
func (h *UserRegistrationHandler) UserRegister(c *gin.Context) {
	var req gen.UserRegisterRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gen.ErrorResponse{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}

	user, err := h.userRegistrationService.UserRegister(req.Email, req.Password, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gen.ErrorResponse{Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, userToUserRegistrationResponse(user))
}
