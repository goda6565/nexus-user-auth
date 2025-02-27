package value

import (
	"regexp"

	"github.com/goda6565/nexus-user-auth/errs"
)

type UserEmail struct {
	value string
}

func (e *UserEmail) Value() string {
	return e.value
}

func (e *UserEmail) Equals(other *UserEmail) bool {
	return e.value == other.value
}

// 空の場合や無効な形式の場合はエラーを返す。
func NewUserEmail(value string) (*UserEmail, error) {
	// 空文字チェック
	if value == "" {
		return nil, errs.NewDomainError("メールアドレスは必須です。")
	}

	// シンプルなメールアドレス形式の正規表現
	const emailRegex = `^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`
	matched, err := regexp.MatchString(emailRegex, value)
	if err != nil {
		return nil, errs.NewDomainError("メールアドレスの形式が正しくありません。")
	}
	if !matched {
		return nil, errs.NewDomainError("メールアドレスの形式が正しくありません。")
	}

	return &UserEmail{value: value}, nil
}
