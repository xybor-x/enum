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

func Test_Enum_StringOf(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, enum.StringOf(RoleUser), "user")
	assert.Equal(t, enum.StringOf(RoleAdmin), "admin")
	assert.Equal(t, enum.StringOf(Role(2)), "Role::<<undefined>>")
}

func Test_Enum_MustStringOf(t *testing.T) {
	type Role int

	assert.Panics(t, func() { enum.MustStringOf(Role(0)) })
	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, enum.MustStringOf(RoleUser), "user")
	assert.Equal(t, enum.MustStringOf(RoleAdmin), "admin")
	assert.Panics(t, func() { enum.MustStringOf(Role(2)) })
}

func Test_Enum_EnumOf(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	userRole, _ := enum.EnumOf[Role]("user")
	assert.Equal(t, userRole, RoleUser)
	adminRole, _ := enum.EnumOf[Role]("admin")
	assert.Equal(t, adminRole, RoleAdmin)
	_, valid := enum.EnumOf[Role]("moderator")
	assert.False(t, valid)
}

func Test_Enum_MustEnumOf(t *testing.T) {
	type Role int

	assert.Panics(t, func() { enum.MustEnumOf[Role]("moderator") })

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, enum.MustEnumOf[Role]("user"), RoleUser)
	assert.Equal(t, enum.MustEnumOf[Role]("admin"), RoleAdmin)
	assert.Panics(t, func() { enum.MustEnumOf[Role]("moderator") })
}

func Test_Enum_Undefined(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.True(t, enum.IsValid(RoleUser))
	assert.True(t, enum.IsValid(RoleAdmin))

	_, ok := enum.EnumOf[Role]("moderator")
	assert.False(t, ok)
	// assert.False(t, enum.IsValid(moderator))
}

func Test_Enum_UndefinedEnum(t *testing.T) {
	type Role int

	moderator, _ := enum.EnumOf[Role]("moderator")
	assert.False(t, enum.IsValid(moderator))
	assert.False(t, enum.IsValid(Role(0)))
}

func Test_Enum_MarshalJSON(t *testing.T) {
	type Role int

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := enum.MarshalJSON(RoleUser)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"user"`), data)

	_, err = enum.MarshalJSON(Role(1))
	assert.ErrorContains(t, err, "unknown Role: 1")
}

func Test_Enum_UnmarshalJSON(t *testing.T) {
	type Role int

	var (
		RoleUser = enum.New[Role]("user")
	)

	var data Role

	err := enum.UnmarshalJSON([]byte(`"user"`), &data)
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	// Invalid data
	err = enum.UnmarshalJSON([]byte(`user"`), &data)
	assert.ErrorContains(t, err, "invalid character")

	// Invalid enum
	err = enum.UnmarshalJSON([]byte(`"admin"`), &data)
	assert.ErrorContains(t, err, "unknown Role string: admin")
}

func Test_Enum_All(t *testing.T) {
	type Role int

	assert.Nil(t, enum.All[Role]())

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	all := enum.All[Role]()
	assert.Contains(t, all, RoleUser)
	assert.Contains(t, all, RoleAdmin)
}

func Test_Enum_Non_Int(t *testing.T) {
	type Role byte

	assert.Nil(t, enum.All[Role]())

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	all := enum.All[Role]()
	assert.Contains(t, all, RoleUser)
	assert.Contains(t, all, RoleAdmin)
}
