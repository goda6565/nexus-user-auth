package value

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/goda6565/nexus-user-auth/errs"
)

type UserObjID struct {
	value string
}

func (ins *UserObjID) Value() string {
	return ins.value
}

func (ins *UserObjID) Equals(value *UserObjID) bool {
	return ins.value == value.Value()
}

func NewUserObjID(value string) (*UserObjID, error) {
	const LENGTH int = 36
	const REGEXP string = "(^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$)"
	if utf8.RuneCountInString(value) != LENGTH {
		return nil, errs.NewDomainError(fmt.Sprintf("オブジェクトIDの文字数は%d文字でなければなりません。", LENGTH))
	}
	if !regexp.MustCompile(REGEXP).Match([]byte(value)) {
		return nil, errs.NewDomainError("オブジェクトIDはUUID形式でなければなりません。")
	}
	return &UserObjID{value: value}, nil
}
