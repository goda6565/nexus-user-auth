package value

import (
	"fmt"
	"unicode/utf8"

	"github.com/goda6565/nexus-user-auth/errs"
)

type UserUsername struct {
	value string
}

func (n *UserUsername) Value() string {
	return n.value
}

// ユーザー名は3文字以上50文字以内である必要がある。
func NewUserUsername(value string) (*UserUsername, error) {
	const minLength = 3
	const maxLength = 50

	count := utf8.RuneCountInString(value)
	if count < minLength || count > maxLength {
		return nil, errs.NewDomainError(fmt.Sprintf("ユーザー名は%d文字以上%d文字以内でなければなりません。", minLength, maxLength))
	}
	return &UserUsername{value: value}, nil
}
