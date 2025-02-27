package value

import (
	"fmt"

	"github.com/goda6565/nexus-user-auth/errs"
)

const (
	Admin       = "admin"
	RegularUser = "user"
)

type UserRole struct {
	value string
}

func (r *UserRole) Value() string {
	return r.value
}

func NewUserRole(value string) (*UserRole, error) {
	// ユーザーロールの値が正しいかチェックする
	switch value {
	case Admin, RegularUser:
		return &UserRole{value: value}, nil
	default:
		return nil, errs.NewDomainError(fmt.Sprintf("無効なユーザーロール: %s", value))
	}
}
