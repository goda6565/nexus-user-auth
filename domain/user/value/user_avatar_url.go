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
	if value == "" {
		// 空の場合は、値がない状態として許容する
		return &UserAvatarURL{value: value}, nil
	}
	// 値がある場合は、URL形式の検証を行う
	_, err := url.ParseRequestURI(value)
	if err != nil {
		return nil, errs.NewDomainError(fmt.Sprintf("アバターURLの形式が正しくありません: %s", value))
	}
	return &UserAvatarURL{value: value}, nil
}
