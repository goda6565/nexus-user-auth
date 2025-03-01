package profile

import (
	"fmt"

	"github.com/goda6565/nexus-user-auth/domain/user/entity"
	"github.com/goda6565/nexus-user-auth/domain/user/repository"
	"github.com/goda6565/nexus-user-auth/domain/user/value"
	"github.com/goda6565/nexus-user-auth/errs"
)

type UserProfileService interface {
	// UserGet: ユーザー情報取得(認可はミドルウェアで行う)
	UserGet(objID string) (*entity.User, error)
	// UserUpdate: ユーザー情報更新(認可はミドルウェアで行う)
	UserUpdate(objID string, username string, avatarURL string) (*entity.User, error)
	// UserDelete: ユーザー削除(認可はミドルウェアで行う)
	UserDelete(objID string) error
}

type userProfileService struct {
	userRepository repository.UserRepository
}

func NewUserProfileService(userRepository repository.UserRepository) UserProfileService {
	return &userProfileService{
		userRepository: userRepository,
	}
}

// UserGet: ユーザー情報取得(認可はミドルウェアで行う)
func (s *userProfileService) UserGet(objID string) (*entity.User, error) {
	user, err := s.userRepository.GetUserByObjID(objID)
	if err != nil {
		return nil, errs.NewServiceError("failed to get user")
	}
	return user, nil
}

func (s *userProfileService) UserUpdate(objID string, username string, avatarURL string) (*entity.User, error) {
	// ユーザーを取得
	user, err := s.userRepository.GetUserByObjID(objID)
	if err != nil {
		return nil, errs.NewServiceError("failed to get user")
	}

	// ユーザー名の更新
	fmt.Println("username: ", username) 
	newUsername, err := value.NewUserUsername(username)
	if err != nil {
		return nil, errs.NewServiceError("failed to create user username")
	}
	user.ChangeUsername(newUsername)

	// アバターURLの更新（空文字なら変更しない）
	if avatarURL != "" {
		newAvatarURL, err := value.NewUserAvatarURL(avatarURL)
		if err != nil {
			return nil, errs.NewServiceError("failed to create user avatar url")
		}
		user.ChangeAvatarURL(newAvatarURL)
	}

	// 更新をリポジトリに保存
	updatedUser, err := s.userRepository.UpdateUser(user)
	if err != nil {
		return nil, errs.NewServiceError("failed to update user in repository")
	}

	return updatedUser, nil
}

func (s *userProfileService) UserDelete(objID string) error {
	err := s.userRepository.DeleteUser(objID)
	if err != nil {
		return errs.NewServiceError("failed to delete user")
	}
	return nil
}
