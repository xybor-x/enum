package testing_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xybor-x/enum"
	"github.com/xybor-x/enum/testing/proto"
)

func TestProtoNew(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user", proto.ProtoRole_User)
		RoleAdmin = enum.New[Role]("admin", proto.ProtoRole_Admin)
	)

	assert.Equal(t, Role(0), RoleUser)
	assert.Equal(t, Role(1), RoleAdmin)

	r, ok := enum.From[Role]("user")
	assert.True(t, ok)
	assert.Equal(t, RoleUser, r)

	r, ok = enum.From[Role](proto.ProtoRole_User)
	assert.True(t, ok)
	assert.Equal(t, RoleUser, r)

	r, ok = enum.From[Role]("admin")
	assert.True(t, ok)
	assert.Equal(t, RoleAdmin, r)

	r, ok = enum.From[Role](proto.ProtoRole_Admin)
	assert.True(t, ok)
	assert.Equal(t, RoleAdmin, r)

	assert.Equal(t, "user", enum.MustTo[string](RoleUser))
	assert.Equal(t, proto.ProtoRole_User, enum.MustTo[proto.ProtoRole](RoleUser))

	assert.Equal(t, "admin", enum.MustTo[string](RoleAdmin))
	assert.Equal(t, proto.ProtoRole_Admin, enum.MustTo[proto.ProtoRole](RoleAdmin))
}

func TestProtoMap(t *testing.T) {
	type Role int

	const (
		RoleUser Role = iota
		RoleAdmin
	)

	var (
		_ = enum.Map(RoleUser, "user", proto.ProtoRole_User)
		_ = enum.Map(RoleAdmin, "admin", proto.ProtoRole_Admin)
	)

	assert.Equal(t, Role(0), RoleUser)
	assert.Equal(t, Role(1), RoleAdmin)

	r, ok := enum.From[Role]("user")
	assert.True(t, ok)
	assert.Equal(t, RoleUser, r)

	r, ok = enum.From[Role](proto.ProtoRole_User)
	assert.True(t, ok)
	assert.Equal(t, RoleUser, r)

	r, ok = enum.From[Role]("admin")
	assert.True(t, ok)
	assert.Equal(t, RoleAdmin, r)

	r, ok = enum.From[Role](proto.ProtoRole_Admin)
	assert.True(t, ok)
	assert.Equal(t, RoleAdmin, r)

	assert.Equal(t, "user", enum.MustTo[string](RoleUser))
	assert.Equal(t, proto.ProtoRole_User, enum.MustTo[proto.ProtoRole](RoleUser))

	assert.Equal(t, "admin", enum.MustTo[string](RoleAdmin))
	assert.Equal(t, proto.ProtoRole_Admin, enum.MustTo[proto.ProtoRole](RoleAdmin))
}

func TestProtoMapOnlyProtoEnum(t *testing.T) {
	type Role int

	const (
		RoleUser Role = iota
		RoleAdmin
	)

	var (
		_ = enum.Map(RoleUser, proto.ProtoRole_User)
		_ = enum.Map(RoleAdmin, proto.ProtoRole_Admin)
	)

	assert.Equal(t, Role(0), RoleUser)
	assert.Equal(t, Role(1), RoleAdmin)

	r, ok := enum.From[Role]("User")
	assert.True(t, ok)
	assert.Equal(t, RoleUser, r)

	r, ok = enum.From[Role](proto.ProtoRole_User)
	assert.True(t, ok)
	assert.Equal(t, RoleUser, r)

	r, ok = enum.From[Role]("Admin")
	assert.True(t, ok)
	assert.Equal(t, RoleAdmin, r)

	r, ok = enum.From[Role](proto.ProtoRole_Admin)
	assert.True(t, ok)
	assert.Equal(t, RoleAdmin, r)

	assert.Equal(t, "User", enum.MustTo[string](RoleUser))
	assert.Equal(t, proto.ProtoRole_User, enum.MustTo[proto.ProtoRole](RoleUser))

	assert.Equal(t, "Admin", enum.MustTo[string](RoleAdmin))
	assert.Equal(t, proto.ProtoRole_Admin, enum.MustTo[proto.ProtoRole](RoleAdmin))
}

func TestProtoStructNoNeedToHaveRepr(t *testing.T) {
	type role struct{}
	type Role = enum.WrapUintEnum[role]

	const (
		RoleUser Role = iota
		RoleAdmin
	)

	var (
		_ = enum.New[Role](proto.ProtoRole_User)
		_ = enum.New[Role]("admin")
	)
}

func TestProtoAnonymousStructNoNeedToHaveRepr(t *testing.T) {
	type Role = enum.WrapUintEnum[struct{}]

	const (
		RoleUser Role = iota
		RoleAdmin
	)

	var (
		_ = enum.New[Role](proto.ProtoRole_User)
		_ = enum.New[Role]("admin")
	)
}

func TestProtoNewMustImpl(t *testing.T) {
	type Role = enum.WrapUintEnum[proto.ProtoRole]

	const (
		RoleUser Role = iota
		RoleAdmin
	)

	var (
		_ = enum.New[Role](proto.ProtoRole_User)
	)

	assert.PanicsWithValue(t,
		"enum WrapUintEnum[ProtoRole] (1 (admin)): require a representation of proto.ProtoRole",
		func() { enum.New[Role]("admin") },
	)
}

func TestProtoMapMustImpl(t *testing.T) {
	type Role = enum.WrapEnum[proto.ProtoRole]

	const (
		RoleUser Role = iota
		RoleAdmin
	)

	var (
		_ = enum.Map(RoleUser, proto.ProtoRole_User)
	)

	assert.PanicsWithValue(t,
		"enum WrapEnum[ProtoRole] (1): representation User of proto.ProtoRole was already mapped to User",
		func() { enum.Map(RoleAdmin, proto.ProtoRole_User) },
	)

	assert.PanicsWithValue(t,
		"enum WrapEnum[ProtoRole] (1 (admin)): require a representation of proto.ProtoRole",
		func() { enum.Map(RoleAdmin, "admin") },
	)
}
