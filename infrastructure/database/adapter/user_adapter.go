package adapter

import (
	"time"

	"github.com/goda6565/nexus-user-auth/domain/timeobj"
	userEntity "github.com/goda6565/nexus-user-auth/domain/user/entity"
	"github.com/goda6565/nexus-user-auth/domain/user/value"
	"github.com/goda6565/nexus-user-auth/errs"
	"github.com/goda6565/nexus-user-auth/infrastructure/database/models"
)

// UserAdapter は、ドメインのユーザーと永続化用モデル間の変換を行うためのインターフェースです。
type UserAdapter interface {
	// Convert は、ドメインエンティティから GORM モデルへ変換します。
	Convert(source *userEntity.User) any
	// ReBuild は、GORM モデルからドメインエンティティへ再構築します。
	ReBuild(source any) (*userEntity.User, error)
}

// userAdapterImpl は、UserAdapter の実装です。
type userAdapterImpl struct{}

// NewUserAdapter は、UserAdapter の実装を返します。
func NewUserAdapter() UserAdapter {
	return &userAdapterImpl{}
}

func (a *userAdapterImpl) Convert(source *userEntity.User) any {
	// EmailVerifiedAt, LastLoginAt は nil チェックを行い、存在すれば time.Time に変換
	var emailVerifiedAt *time.Time
	if source.EmailVerifiedAt() != nil {
		t := source.EmailVerifiedAt().Value()
		emailVerifiedAt = &t
	}
	var lastLoginAt *time.Time
	if source.LastLoginAt() != nil {
		t := source.LastLoginAt().Value()
		lastLoginAt = &t
	}
	// AvatarURL が nil なら空文字とする
	avatar := ""
	if source.AvatarURL() != nil {
		avatar = source.AvatarURL().Value()
	}

	return &models.User{
		ObjID:           source.ObjID().Value(),
		Email:           source.Email().Value(),
		Password:        source.Password().Value(),
		Username:        source.Username().Value(),
		AvatarURL:       avatar,
		EmailVerifiedAt: emailVerifiedAt,
		LastLoginAt:     lastLoginAt,
		Role:            source.Role().Value(),
	}
}

func (a *userAdapterImpl) ReBuild(source any) (*userEntity.User, error) {
	userModel, ok := source.(*models.User)
	if !ok {
		return nil, errs.NewInfraError("*models.User以外の値が指定されました。")
	}

	// 各値オブジェクトを生成
	email, err := value.NewUserEmail(userModel.Email)
	if err != nil {
		return nil, err
	}
	username, err := value.NewUserUsername(userModel.Username)
	if err != nil {
		return nil, err
	}
	var avatarURL *value.UserAvatarURL
	if userModel.AvatarURL != "" {
		avatarURL, err = value.NewUserAvatarURL(userModel.AvatarURL)
		if err != nil {
			return nil, err
		}
	}
	// timeobj.TimeObj に変換（存在すれば）
	var emailVerifiedAt *timeobj.TimeObj
	if userModel.EmailVerifiedAt != nil {
		emailVerifiedAt, err = timeobj.NewTimeObj(*userModel.EmailVerifiedAt)
		if err != nil {
			return nil, err
		}
	}
	var lastLoginAt *timeobj.TimeObj
	if userModel.LastLoginAt != nil {
		lastLoginAt, err = timeobj.NewTimeObj(*userModel.LastLoginAt)
		if err != nil {
			return nil, err
		}
	}
	role, err := value.NewUserRole(userModel.Role)
	if err != nil {
		return nil, err
	}
	objID, err := value.NewUserObjID(userModel.ObjID)
	if err != nil {
		return nil, err
	}

	// BuildUser は、既存データからドメインエンティティを再構築するためのファクトリ関数です。
	return userEntity.BuildUser(objID, email, value.FromHashed(userModel.Password), username, avatarURL, emailVerifiedAt, lastLoginAt, role)
}
