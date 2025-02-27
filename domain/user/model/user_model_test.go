package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/goda6565/nexus-user-auth/domain/user/value"
	"github.com/goda6565/nexus-user-auth/domain/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func dummyUserEmail() *value.UserEmail {
	email := fmt.Sprintf("user_%s@example.com", uuid.New().String())
	emailObj, err := value.NewUserEmail(email)
	if err != nil {
		panic(err)
	}
	return emailObj
}

func dummyUserPassword() *value.UserPassword {
	password := uuid.New().String()
	passwordObj, err := value.NewUserPassword(password)
	if err != nil {
		panic(err)
	}
	return passwordObj
}

func dummyUserUsername() *value.UserUsername {
	username := fmt.Sprintf("username_%s", uuid.New().String())
	usernameObj, err := value.NewUserUsername(username)
	if err != nil {
		panic(err)
	}
	return usernameObj
}

func dummyUserAvatarURL() *value.UserAvatarURL {
	avatarURL := fmt.Sprintf("https://avatar.example.com/%s.png", uuid.New().String())
	avatarURLObj, err := value.NewUserAvatarURL(avatarURL)
	if err != nil {
		panic(err)
	}
	return avatarURLObj
}

func dummyUserRole() *value.UserRole {
	role, err := value.NewUserRole("user")
	if err != nil {
		panic(err)
	}
	return role
}

func dummyTimeObj() *utils.TimeObj {
	// 現在時刻を元に TimeObj を生成する（エラー処理は panic で簡略化）
	timeObj, err := utils.NewTimeObj(time.Now())
	if err != nil {
		panic(err)
	}
	return timeObj
}

func TestNewUser(t *testing.T) {
	// ダミーのパラメータを生成
	email := dummyUserEmail()
	password := dummyUserPassword()
	username := dummyUserUsername()
	avatarURL := dummyUserAvatarURL()
	emailVerifiedAt := dummyTimeObj()
	lastLoginAt := dummyTimeObj()
	role := dummyUserRole()

	u, err := NewUser(email, password, username, avatarURL, emailVerifiedAt, lastLoginAt, role)
	assert.NoError(t, err, "NewUser でエラーが返られないこと")
	assert.NotNil(t, u, "NewUser で生成されたユーザが nil でないこと")
	assert.NotNil(t, u.ObjID(), "生成されたユーザの ObjID が nil でないこと")
	// 各フィールドが正しくセットされているかの確認
	assert.Equal(t, email, u.Email(), "Email が正しくセットされていること")
	assert.Equal(t, password, u.Password(), "Password が正しくセットされていること")
	assert.Equal(t, username, u.Username(), "Username が正しくセットされていること")
	assert.Equal(t, avatarURL, u.AvatarURL(), "AvatarURL が正しくセットされていること")
	assert.Equal(t, role, u.Role(), "Role が正しくセットされていること")
}

func TestBuildUser(t *testing.T) {
	// 既存のユーザオブジェクトを用意する場合
	email := dummyUserEmail()
	password := dummyUserPassword()
	username := dummyUserUsername()
	avatarURL := dummyUserAvatarURL()
	emailVerifiedAt := dummyTimeObj()
	lastLoginAt := dummyTimeObj()
	role := dummyUserRole()

	// NewUser を使って一度ユーザを生成し、その ObjID を用いて BuildUser を呼び出す
	u1, err := NewUser(email, password, username, avatarURL, emailVerifiedAt, lastLoginAt, role)
	assert.NoError(t, err)

	u2, err := BuildUser(u1.ObjID(), email, password, username, avatarURL, emailVerifiedAt, lastLoginAt, role)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
	assert.Equal(t, u1.ObjID(), u2.ObjID(), "ObjID が一致していること")
}

func TestUserEquals(t *testing.T) {
	// 2つのユーザを生成
	email := dummyUserEmail()
	password := dummyUserPassword()
	username := dummyUserUsername()
	avatarURL := dummyUserAvatarURL()
	emailVerifiedAt := dummyTimeObj()
	lastLoginAt := dummyTimeObj()
	role := dummyUserRole()

	u1, err := NewUser(email, password, username, avatarURL, emailVerifiedAt, lastLoginAt, role)
	assert.NoError(t, err)
	u2, err := BuildUser(u1.ObjID(), email, password, username, avatarURL, emailVerifiedAt, lastLoginAt, role)
	assert.NoError(t, err)

	// Equals で同一のオブジェクトと判断されること
	equal, err := u1.Equals(u2)
	assert.NoError(t, err)
	assert.True(t, equal, "同じ ObjID を持つ場合 Equals は true を返す")

	// nil を引数にした場合、エラーが返ることを検証
	equal, err = u1.Equals(nil)
	assert.Error(t, err, "nil を引数にした場合はエラーが返ること")
	assert.False(t, equal, "nil を引数にした場合 Equals は false を返す")
}
