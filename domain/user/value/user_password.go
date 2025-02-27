package value

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/goda6565/nexus-user-auth/errs"
)

// UserPassword はユーザーのパスワードを表す値オブジェクトです。
type UserPassword struct {
	value string
}

func (p *UserPassword) Value() string {
	return p.value
}

func NewUserPassword(value string) (*UserPassword, error) {
	const minLength = 8  // パスワードの最低文字数
	const maxLength = 64 // パスワードの最大文字数

	length := utf8.RuneCountInString(value) // パスワードの文字数を取得
	if length < minLength || length > maxLength {
		return nil, errs.NewDomainError(fmt.Sprintf("パスワードは%d文字以上%d文字以内でなければなりません。", minLength, maxLength))
	}

	// 英字と数字が含まれているかチェックする
	hasLetter, _ := regexp.MatchString(`[A-Za-z]`, value)
	hasDigit, _ := regexp.MatchString(`[0-9]`, value)
	if !hasLetter || !hasDigit {
		return nil, errs.NewDomainError("パスワードは英字と数字の両方を含む必要があります。")
	}

	return &UserPassword{value: value}, nil
}
