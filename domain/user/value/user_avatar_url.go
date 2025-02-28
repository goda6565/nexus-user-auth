package value

import (
	"fmt"
	"net/url"

	"github.com/goda6565/nexus-user-auth/errs"
)

type UserAvatarURL struct {
	value string
}

func (u *UserAvatarURL) Value() string {
	return u.value
}

func NewUserAvatarURL(value string) (*UserAvatarURL, error) {
	_, err := url.ParseRequestURI(value)
	if err != nil {
		return nil, errs.NewDomainError(fmt.Sprintf("アバターURLの形式が正しくありません: %s", value))
	}
	return &UserAvatarURL{value: value}, nil
}
