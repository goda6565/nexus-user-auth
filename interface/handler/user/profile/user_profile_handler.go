package profile

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
	// user.AvatarURL() の nil チェックを追加
	var avatar string
	if user.AvatarURL() != nil {
		avatar = user.AvatarURL().Value()
	}

	if avatar == "" {
		return &gen.ProfileResponse{
			Uid:      user.ObjID().Value(),
			Email:    user.Email().Value(),
			Username: user.Username().Value(),
		}
	}
	return &gen.ProfileResponse{
		Uid:       user.ObjID().Value(),
		AvatarURL: &avatar,
		Email:     user.Email().Value(),
		Username:  user.Username().Value(),
	}
}

// getValidatedUID: Gin の Context から認証済みユーザーIDを取得するヘルパー関数
func getValidatedUID(c *gin.Context) (string, bool) {
	objID, exists := c.Get("validated_uid")
	if !exists {
		return "", false
	}

	objIDStr, ok := objID.(string)
	if !ok || objIDStr == "" {
		return "", false
	}

	return objIDStr, true
}

// GetUserProfile: ユーザー情報取得
func (h *UserProfileHandler) GetUserProfile(c *gin.Context) {
	objID, exists := getValidatedUID(c)
	if !exists {
		c.JSON(http.StatusBadRequest, gen.ErrorResponse{
			Message: "ユーザーIDが取得できません",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, err := h.userProfileService.UserGet(objID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gen.ErrorResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gen.ErrorResponse{
			Message: "ユーザーが見つかりません",
			Code:    http.StatusNotFound,
		})
		return
	}
	response := userToUserProfileResponse(user)
	c.JSON(http.StatusOK, response)
}

// UpdateUserProfile: ユーザー情報更新
func (h *UserProfileHandler) UpdateUserProfile(c *gin.Context) {
	objID, exists := getValidatedUID(c)
	if !exists {
		c.JSON(http.StatusBadRequest, gen.ErrorResponse{
			Message: "ユーザーIDが取得できません",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req gen.UserProfileUpdateRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gen.ErrorResponse{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, err := h.userProfileService.UserUpdate(objID, req.Username, *req.AvatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gen.ErrorResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	response := userToUserProfileResponse(user)
	c.JSON(http.StatusOK, response)
}

// DeleteUserProfile: ユーザー削除
func (h *UserProfileHandler) DeleteUserProfile(c *gin.Context) {
	objID, exists := getValidatedUID(c)
	if !exists {
		c.JSON(http.StatusBadRequest, gen.ErrorResponse{
			Message: "ユーザーIDが取得できません",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err := h.userProfileService.UserDelete(objID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gen.ErrorResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, nil)
}
