package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserRole_Valid(t *testing.T) {
	roleAdmin, err := NewUserRole(Admin)
	assert.NoError(t, err, "admin は有効なロールであること")
	assert.NotNil(t, roleAdmin)
	assert.Equal(t, Admin, roleAdmin.Value())

	roleUser, err := NewUserRole(RegularUser)
	assert.NoError(t, err, "user は有効なロールであること")
	assert.NotNil(t, roleUser)
	assert.Equal(t, RegularUser, roleUser.Value())
}

func TestNewUserRole_Invalid(t *testing.T) {
	role, err := NewUserRole("unknown")
	assert.Error(t, err, "不正なロールはエラーになること")
	assert.Nil(t, role)
}
