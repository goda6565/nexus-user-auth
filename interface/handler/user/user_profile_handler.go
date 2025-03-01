package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/goda6565/nexus-user-auth/application/service/user/profile"
	"github.com/goda6565/nexus-user-auth/domain/user/entity"
	"github.com/goda6565/nexus-user-auth/interface/gen"
)

type UserProfileHandler struct {
	userProfileService profile.UserProfileService
}

func NewUserProfileHandler(userProfileService profile.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{
		userProfileService: userProfileService,
	}
}

func userToUserProfileResponse(user *entity.User) *gen.ProfileResponse {
	if user.AvatarURL().Value() == "" {
		return &gen.ProfileResponse{
			Uid:      user.ObjID().Value(),
			Email:    user.Email().Value(),
			Username: user.Username().Value(),
		}
	}
	avatarURL := user.AvatarURL().Value()
	return &gen.ProfileResponse{
		Uid:       user.ObjID().Value(),
		AvatarURL: &avatarURL,
		Email:     user.Email().Value(),
		Username:  user.Username().Value(),
	}
}

// UserGet: ユーザー情報取得(認可はミドルウェアで行う)
func (h *UserProfileHandler) GetUserProfile(c *gin.Context) {
	objID := c.Param("valified_uid") // ミドルウェアでセットされたユーザーIDを取得
	user, err := h.userProfileService.UserGet(objID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gen.ErrorResponse{Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, userToUserProfileResponse(user))
}

// UserUpdate: ユーザー情報更新(認可はミドルウェアで行う)
func (h *UserProfileHandler) UpdateUserProfile(c *gin.Context) {
	objID := c.Param("valified_uid") // ミドルウェアでセットされたユーザーIDを取得
	var req gen.UserProfileUpdateRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gen.ErrorResponse{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}

	user, err := h.userProfileService.UserUpdate(objID, *req.AvatarURL, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gen.ErrorResponse{Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, userToUserProfileResponse(user))
}

func (h *UserProfileHandler) DeleteUserProfile(c *gin.Context) {
	objID := c.Param("valified_uid") // ミドルウェアでセットされたユーザーIDを取得
	err := h.userProfileService.UserDelete(objID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gen.ErrorResponse{Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, nil)
}
