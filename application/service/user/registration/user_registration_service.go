package registration

import (
	"github.com/goda6565/nexus-user-auth/domain/user/entity"
	"github.com/goda6565/nexus-user-auth/domain/user/repository"
	"github.com/goda6565/nexus-user-auth/domain/user/value"
	"github.com/goda6565/nexus-user-auth/errs"
)

type UserRegistrationService interface {
	// UserRegister: ユーザー登録
	UserRegister(email string, password string, username string) (*entity.User, error)
}

// UserRegistrationServiceの実装
type userRegistrationService struct {
	userRepository repository.UserRepository
}

// NewUserRegistrationService: UserRegistrationServiceを生成
func NewUserRegistrationService(userRepository repository.UserRepository) UserRegistrationService {
	return &userRegistrationService{
		userRepository: userRepository,
	}
}

// UserRegister: ユーザー登録
func (s *userRegistrationService) UserRegister(email string, password string, username string) (*entity.User, error) {
	// ユーザーエンティティの生成
	emailValue, err := value.NewUserEmail(email)
	if err != nil {
		return nil, errs.NewServiceError("failed to create user email")
	}
	passwordValue, err := value.NewUserPassword(password)
	if err != nil {
		return nil, errs.NewServiceError("failed to create user password")
	}
	usernameValue, err := value.NewUserUsername(username)
	if err != nil {
		return nil, errs.NewServiceError("failed to create user username")
	}

	userEntity, err := entity.NewUser(emailValue, passwordValue, usernameValue)
	if err != nil {
		return nil, err
	}

	// ユーザー作成
	createdUser, err := s.userRepository.CreateUser(userEntity)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
