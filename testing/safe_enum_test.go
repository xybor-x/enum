package testing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xybor-x/enum"
)

func TestNewExtendedSafe(t *testing.T) {
	type Role struct{ enum.SafeEnum[int] }

	var (
		RoleUser  = enum.NewExtendedSafe[Role]("user")
		RoleAdmin = enum.NewExtendedSafe[Role]("admin")
	)

	assert.Equal(t, "user", RoleUser.String())
	assert.Equal(t, "admin", RoleAdmin.String())

	user, ok := enum.FromString[Role]("user")
	assert.True(t, ok)
	assert.Equal(t, RoleUser, user)

	assert.Equal(t, []Role{RoleUser, RoleAdmin}, enum.All[Role]())
}

func TestSafeEnumPrintZeroStruct(t *testing.T) {
	type role any
	type Role = enum.SafeEnum[role]

	var (
		_ = enum.NewSafe[Role]("user")
	)

	type User struct {
		Role Role
	}

	assert.Equal(t, "{<nil>}", fmt.Sprint(User{}))
}
