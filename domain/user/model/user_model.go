package model

import (
	"github.com/goda6565/nexus-user-auth/domain/timeobj"
	"github.com/goda6565/nexus-user-auth/domain/user/value"
	"github.com/goda6565/nexus-user-auth/errs"
	"github.com/google/uuid"
)

type User struct {
	objID           *value.UserObjID
	email           *value.UserEmail
	password        *value.UserPassword
	username        *value.UserUsername
	avatarURL       *value.UserAvatarURL
	emailVerifiedAt *timeobj.TimeObj
	lastLoginAt     *timeobj.TimeObj
	role            *value.UserRole
}

func (ins *User) ObjID() *value.UserObjID {
	return ins.objID
}

func (ins *User) Email() *value.UserEmail {
	return ins.email
}

func (ins *User) Password() *value.UserPassword {
	return ins.password
}

func (ins *User) Username() *value.UserUsername {
	return ins.username
}

func (ins *User) AvatarURL() *value.UserAvatarURL {
	return ins.avatarURL
}

func (ins *User) EmailVerifiedAt() *timeobj.TimeObj {
	return ins.emailVerifiedAt
}

func (ins *User) LastLoginAt() *timeobj.TimeObj {
	return ins.lastLoginAt
}

func (ins *User) Role() *value.UserRole {
	return ins.role
}

// 同一性の確認
func (ins *User) Equals(obj *User) (bool, error) {
	if obj == nil {
		return false, errs.NewDomainError("引数でnilが指定されました。")
	}
	result := ins.objID.Equals(obj.ObjID())
	return result, nil
}

func NewUser(email *value.UserEmail, password *value.UserPassword, username *value.UserUsername, avatarURL *value.UserAvatarURL, emailVerifiedAt *timeobj.TimeObj, lastLoginAt *timeobj.TimeObj, role *value.UserRole) (*User, error) {
	if uid, err := uuid.NewRandom(); err != nil { // UUIDを生成する
		return nil, errs.NewDomainError(err.Error())
	} else {
		if id, err := value.NewUserObjID(uid.String()); err != nil {
			return nil, errs.NewDomainError(err.Error())
		} else {
			return &User{
				objID:           id,
				email:           email,
				password:        password,
				username:        username,
				avatarURL:       avatarURL,
				emailVerifiedAt: emailVerifiedAt,
				lastLoginAt:     lastLoginAt,
				role:            role,
			}, nil
		}
	}
}

func BuildUser(objID *value.UserObjID, email *value.UserEmail, password *value.UserPassword, username *value.UserUsername, avatarURL *value.UserAvatarURL, emailVerifiedAt *timeobj.TimeObj, lastLoginAt *timeobj.TimeObj, role *value.UserRole) (*User, error) {
	return &User{
		objID:           objID,
		email:           email,
		password:        password,
		username:        username,
		avatarURL:       avatarURL,
		emailVerifiedAt: emailVerifiedAt,
		lastLoginAt:     lastLoginAt,
		role:            role,
	}, nil
}
