package repository

import (
	"github.com/goda6565/nexus-user-auth/domain/user/model"
)

type UserRepository interface {
	// CreateUser: 新規ユーザーを作成
	CreateUser(user *model.User) (*model.User, error)

	// GetUserByEmail: 指定のメールアドレスでユーザーを取得
	GetUserByEmail(email string) (*model.User, error)

	// GetUserByObjID: 外部識別用のUUID (ObjID) でユーザーを取得
	GetUserByObjID(objID string) (*model.User, error)

	// UpdateUser: ユーザー情報を更新
	UpdateUser(user *model.User) (*model.User, error)

	// DeleteUser: ユーザーを論理削除
	DeleteUser(objID string) error
}
