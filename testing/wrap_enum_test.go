package testing_test

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xybor-x/enum"
	"gopkg.in/yaml.v3"
)

func TestWrapEnumMarshalJSON(t *testing.T) {
	type role int
	type Role = enum.WrapEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := json.Marshal(RoleUser)
	assert.NoError(t, err)
	assert.Equal(t, `"user"`, string(data))

	_, err = json.Marshal(Role(1))
	assert.ErrorContains(t, err, "enum WrapEnum[role]: invalid value 1")
}

func TestWrapEnumUnmarshalJSON(t *testing.T) {
	type role int
	type Role = enum.WrapEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	var data Role

	err := json.Unmarshal([]byte(`"user"`), &data)
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	err = json.Unmarshal([]byte(`"admin"`), &data)
	assert.ErrorContains(t, err, "enum WrapEnum[role]: unknown string admin")
}

func TestWrapEnumValueSQL(t *testing.T) {
	type role int
	type Role = enum.WrapEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := RoleUser.Value()
	assert.NoError(t, err)
	assert.Equal(t, "user", data)

	_, err = Role(1).Value()
	assert.ErrorContains(t, err, "enum WrapEnum[role]: invalid value 1")
}

func TestWrapEnumScanSQL(t *testing.T) {
	type role int
	type Role = enum.WrapEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	var data Role

	// Scan bytes
	err := data.Scan([]byte(`user`))
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	// Scan string
	err = data.Scan("user")
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	// Invalid enum
	err = data.Scan("admin")
	assert.ErrorContains(t, err, "enum WrapEnum[role]: unknown string admin")
}

func TestWrapUintEnumMarshalJSON(t *testing.T) {
	type role int
	type Role = enum.WrapUintEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := json.Marshal(RoleUser)
	assert.NoError(t, err)
	assert.Equal(t, `"user"`, string(data))

	_, err = json.Marshal(Role(1))
	assert.ErrorContains(t, err, "enum WrapUintEnum[role]: invalid value 1")
}

func TestWrapUintEnumUnmarshalJSON(t *testing.T) {
	type role int
	type Role = enum.WrapUintEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	var data Role

	err := json.Unmarshal([]byte(`"user"`), &data)
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	err = json.Unmarshal([]byte(`"admin"`), &data)
	assert.ErrorContains(t, err, "enum WrapUintEnum[role]: unknown string admin")
}

func TestWrapUintEnumValueSQL(t *testing.T) {
	type role int
	type Role = enum.WrapUintEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := RoleUser.Value()
	assert.NoError(t, err)
	assert.Equal(t, "user", data)

	_, err = Role(1).Value()
	assert.ErrorContains(t, err, "enum WrapUintEnum[role]: invalid value 1")
}

func TestWrapUintEnumScanSQL(t *testing.T) {
	type role int
	type Role = enum.WrapUintEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	var data Role

	// Scan bytes
	err := data.Scan([]byte(`user`))
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	// Scan string
	err = data.Scan("user")
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	// Invalid enum
	err = data.Scan("admin")
	assert.ErrorContains(t, err, "enum WrapUintEnum[role]: unknown string admin")
}

func TestWrapFloatEnumMarshalJSON(t *testing.T) {
	type role int
	type Role = enum.WrapFloatEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := json.Marshal(RoleUser)
	assert.NoError(t, err)
	assert.Equal(t, `"user"`, string(data))

	_, err = json.Marshal(Role(1))
	assert.ErrorContains(t, err, "enum WrapFloatEnum[role]: invalid value 1")
}

func TestWrapFloatEnumUnmarshalJSON(t *testing.T) {
	type role int
	type Role = enum.WrapFloatEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	var data Role

	err := json.Unmarshal([]byte(`"user"`), &data)
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	err = json.Unmarshal([]byte(`"admin"`), &data)
	assert.ErrorContains(t, err, "enum WrapFloatEnum[role]: unknown string admin")
}

func TestWrapFloatEnumValueSQL(t *testing.T) {
	type role int
	type Role = enum.WrapFloatEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := RoleUser.Value()
	assert.NoError(t, err)
	assert.Equal(t, "user", data)

	_, err = Role(1).Value()
	assert.ErrorContains(t, err, "enum WrapFloatEnum[role]: invalid value 1")
}

func TestWrapFloatEnumScanSQL(t *testing.T) {
	type role int
	type Role = enum.WrapFloatEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	var data Role

	// Scan bytes
	err := data.Scan([]byte(`user`))
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	// Scan string
	err = data.Scan("user")
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	// Invalid enum
	err = data.Scan("admin")
	assert.ErrorContains(t, err, "enum WrapFloatEnum[role]: unknown string admin")
}

func TestSafeEnumMarshalJSON(t *testing.T) {
	type role int
	type Role = enum.SafeEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := json.Marshal(RoleUser)
	assert.NoError(t, err)
	assert.Equal(t, `"user"`, string(data))

	_, err = json.Marshal(Role{})
	assert.ErrorContains(t, err, "enum SafeEnum[role]: invalid value <nil>")
}

func TestSafeEnumUnmarshalJSON(t *testing.T) {
	type role int
	type Role = enum.SafeEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	var data Role

	err := json.Unmarshal([]byte(`"user"`), &data)
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	err = json.Unmarshal([]byte(`"admin"`), &data)
	assert.ErrorContains(t, err, "enum SafeEnum[role]: unknown string admin")
}

func TestSafeEnumValueSQL(t *testing.T) {
	type role int
	type Role = enum.SafeEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := RoleUser.Value()
	assert.NoError(t, err)
	assert.Equal(t, "user", data)

	_, err = Role{}.Value()
	assert.ErrorContains(t, err, "enum SafeEnum[role]: invalid value <nil>")
}

func TestSafeEnumScanSQL(t *testing.T) {
	type role int
	type Role = enum.SafeEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	var data Role

	// Scan bytes
	err := data.Scan([]byte(`user`))
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	// Scan string
	err = data.Scan("user")
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	// Invalid enum
	err = data.Scan("admin")
	assert.ErrorContains(t, err, "enum SafeEnum[role]: unknown string admin")
}

func TestWrapEnumMarshalXMLStruct(t *testing.T) {
	type role int
	type Role = enum.WrapEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	type Test1 struct {
		Role Role `xml:"CustomRole"`
	}

	data, err := xml.Marshal(Test1{Role: RoleUser})
	assert.NoError(t, err)
	assert.Equal(t, "<Test1><CustomRole>user</CustomRole></Test1>", string(data))

	type Test2 struct {
		Role Role
	}

	data, err = xml.Marshal(Test2{Role: RoleUser})
	assert.NoError(t, err)
	assert.Equal(t, "<Test2><Role>user</Role></Test2>", string(data))
}

func TestWrapEnumUnmarshalXML(t *testing.T) {
	type role int
	type Role = enum.WrapEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	var data Role

	err := xml.Unmarshal([]byte(`<Role>user</Role>`), &data)
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	err = xml.Unmarshal([]byte(`<Role>admin</Role>`), &data)
	assert.ErrorContains(t, err, "enum WrapEnum[role]: unknown string admin")
}

func TestWrapEnumMarshalYAML(t *testing.T) {
	type role int
	type Role = enum.WrapEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := yaml.Marshal(RoleUser)
	assert.NoError(t, err)
	assert.Equal(t, "user\n", string(data))

	_, err = yaml.Marshal(Role(1))
	assert.ErrorContains(t, err, "enum WrapEnum[role]: invalid value 1")
}

func TestWrapEnumUnmarshalYAML(t *testing.T) {
	type role int
	type Role = enum.WrapEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	type Test struct {
		Role Role `yaml:"role"`
	}
	var data Test

	err := yaml.Unmarshal([]byte(`role: user`), &data)
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data.Role)

	err = yaml.Unmarshal([]byte("role: admin"), &data)
	assert.ErrorContains(t, err, "enum WrapEnum[role]: unknown string admin")

	err = yaml.Unmarshal([]byte("role:\n- user\n"), &data)
	assert.ErrorContains(t, err, "enum WrapEnum[role]: only supports scalar in yaml enum")
}
