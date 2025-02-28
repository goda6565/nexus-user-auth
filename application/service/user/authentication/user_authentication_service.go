package authentication

import (
	"github.com/goda6565/nexus-user-auth/domain/user/repository"
	"github.com/goda6565/nexus-user-auth/errs"
	"github.com/goda6565/nexus-user-auth/pkg/utils"
)

type UserAuthenticationService interface {
	// UserLogin: ユーザーログイン
	UserLogin(email string, password string) (accessToken string, refreshToken string, err error)
	// UserTokenRefresh: トークンリフレッシュ
	UserTokenRefresh(refreshToken string) (accessToken string, err error)
}

// userAuthenticationService は UserAuthenticationService の実装
type userAuthenticationService struct {
	userRepository repository.UserRepository
}

// NewUserAuthenticationService は UserAuthenticationService のインスタンスを作成
func NewUserAuthenticationService(userRepository repository.UserRepository) UserAuthenticationService {
	return &userAuthenticationService{
		userRepository: userRepository,
	}
}

// UserLogin はユーザー認証を行い、アクセストークンとリフレッシュトークンを発行
func (s *userAuthenticationService) UserLogin(email string, password string) (string, string, error) {
	// ユーザー取得
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return "", "", errs.NewServiceError("invalid email or password")
	}

	// パスワードの検証
	err = utils.CheckPassword(user.Password().Value(), password)
	if err != nil {
		return "", "", errs.NewServiceError("invalid email or password")
	}

	// トークン生成
	accessToken, refreshToken, err := utils.GenerateTokens(user.ObjID().Value())
	if err != nil {
		return "", "", errs.NewServiceError("failed to generate tokens")
	}

	return accessToken, refreshToken, nil
}

// UserTokenRefresh はリフレッシュトークンを用いて新しいアクセストークンを発行
func (s *userAuthenticationService) UserTokenRefresh(refreshToken string) (string, error) {
	// リフレッシュトークンを検証
	_, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", errs.NewServiceError("invalid refresh token")
	}

	// 新しいアクセストークンを発行
	newAccessToken, err := utils.RefreshAccessToken(refreshToken)
	if err != nil {
		return "", errs.NewServiceError("failed to refresh token")
	}

	return newAccessToken, nil
}
