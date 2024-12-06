package enum_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xybor-x/enum"
)

func Test_Enum_New(t *testing.T) {
	type Role int
	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, RoleUser, Role(0))
	assert.Equal(t, RoleAdmin, Role(1))

	type File int
	var (
		FileImage = enum.New[File]("image")
		FilePDF   = enum.New[File]("pdf")
	)

	assert.Equal(t, FileImage, File(0))
	assert.Equal(t, FilePDF, File(1))

	assert.NotEqual(t, File(0), Role(0))
}

func Test_Enum_Map(t *testing.T) {
	type Role int
	const (
		RoleUser Role = iota
		RoleAdmin
	)

	assert.Equal(t, enum.StringOf(RoleUser), "Role::<<undefined>>")
	assert.Equal(t, enum.StringOf(RoleAdmin), "Role::<<undefined>>")

	var (
		_ = enum.Map(RoleUser, "user")
		_ = enum.Map(RoleAdmin, "admin")
	)

	assert.Equal(t, enum.StringOf(RoleUser), "user")
	assert.Equal(t, enum.StringOf(RoleAdmin), "admin")
}

func Test_Enum_Map_Duplicated(t *testing.T) {
	type Role int

	var (
		RoleUser = enum.New[Role]("user")
	)

	const (
		RoleAdmin Role = iota
	)

	assert.Equal(t, enum.StringOf(RoleUser), "user")
	assert.Panics(t, func() { enum.Map(RoleAdmin, "admin") })
}

func Test_Enum_String(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, enum.StringOf(RoleUser), "user")
	assert.Equal(t, enum.StringOf(RoleAdmin), "admin")
}

func Test_Enum_EnumOf(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, enum.EnumOf[Role]("user"), RoleUser)
	assert.Equal(t, enum.EnumOf[Role]("admin"), RoleAdmin)
}

func Test_Enum_Undefined(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	moderator := enum.EnumOf[Role]("moderator")
	assert.NotEqual(t, moderator, RoleUser)
	assert.NotEqual(t, moderator, RoleAdmin)
	assert.Equal(t, enum.StringOf(moderator), "Role::<<undefined>>")
}

func Test_Enum_UndefinedEnum(t *testing.T) {
	type Role int

	moderator := enum.EnumOf[Role]("moderator")
	assert.Equal(t, enum.StringOf(moderator), "Role::<<undefined>>")
}
