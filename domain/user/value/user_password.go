package value

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/goda6565/nexus-user-auth/errs"
	"github.com/goda6565/nexus-user-auth/pkg/utils"
)

// UserPassword はハッシュ化されたパスワード値を保持する値オブジェクトです。
// ドメイン内では生のパスワードを扱わず、常にハッシュ済みの値を利用します。
type UserPassword struct {
	hashed string
}

// Value はハッシュ化されたパスワード文字列を返します。
func (p *UserPassword) Value() string {
	return p.hashed
}

// NewUserPassword はプレーンなパスワードを受け取り、バリデーションとハッシュ化を行い
// UserPassword を生成します。
func NewUserPassword(plain string) (*UserPassword, error) {
	const minLength = 8  // パスワードの最低文字数
	const maxLength = 64 // パスワードの最大文字数

	length := utf8.RuneCountInString(plain)
	if length < minLength || length > maxLength {
		return nil, errs.NewDomainError(fmt.Sprintf("パスワードは%d文字以上%d文字以内でなければなりません。", minLength, maxLength))
	}

	// 英字と数字が含まれているかチェックする
	hasLetter, _ := regexp.MatchString(`[A-Za-z]`, plain)
	hasDigit, _ := regexp.MatchString(`[0-9]`, plain)
	if !hasLetter || !hasDigit {
		return nil, errs.NewDomainError("パスワードは英字と数字の両方を含む必要があります。")
	}

	hashed, err := utils.HashPassword(plain)
	if err != nil {
		return nil, errs.NewDomainError("パスワードのハッシュ化に失敗しました。")
	}

	return &UserPassword{hashed: hashed}, nil
}

// FromHashed は、永続化層から復元する際に、すでにハッシュ化されたパスワード値から
// UserPassword を構築します。
func FromHashed(hashed string) *UserPassword {
	return &UserPassword{hashed: hashed}
}

// Verify は、プレーンなパスワードと内部に保持しているハッシュ値を比較し、一致するかを判定します。
func (p *UserPassword) Verify(plain string) bool {
	err := utils.CheckPassword(p.hashed, plain)
	return err == nil
}
