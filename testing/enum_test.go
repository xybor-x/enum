package testing_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xybor-x/enum"

	_ "github.com/mattn/go-sqlite3"
)

func TestEnumNew(t *testing.T) {
	type Role int
	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, RoleUser, Role(0))
	assert.Equal(t, RoleAdmin, Role(1))
	assert.Equal(t, "admin", enum.ToString(RoleAdmin))

	type File int
	var (
		FileImage = enum.New[File]("image")
		FilePDF   = enum.New[File]("pdf")
	)

	assert.Equal(t, FileImage, File(0))
	assert.Equal(t, FilePDF, File(1))
	assert.Equal(t, "pdf", enum.ToString(FilePDF))

	assert.NotEqual(t, File(0), Role(0))
}

func TestEnumNewString(t *testing.T) {
	type Role string
	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, RoleUser, Role("user"))
	assert.Equal(t, RoleAdmin, Role("admin"))
	assert.Equal(t, "admin", enum.ToString(RoleAdmin))
	assert.Equal(t, 1, enum.ToInt(RoleAdmin))
}

func TestEnumMap(t *testing.T) {
	type Role int
	const (
		RoleUser Role = iota
		RoleAdmin
	)

	var (
		_ = enum.Map(RoleUser, "user")
		_ = enum.Map(RoleAdmin, "admin")
	)

	assert.Equal(t, enum.ToString(RoleUser), "user")
	assert.Equal(t, enum.ToString(RoleAdmin), "admin")
}

func TestEnumMapDiffEnumber(t *testing.T) {
	type Role int
	const (
		RoleUser Role = iota + 1
		RoleAdmin
	)

	var (
		_ = enum.Map(RoleUser, "user")
		_ = enum.Map(RoleAdmin, "admin")
	)

	assert.Equal(t, int(RoleUser), 1)
	assert.Equal(t, int(RoleAdmin), 2)
}

func TestEnumFinalize(t *testing.T) {
	type Role int
	const (
		RoleUser Role = iota
		RoleAdmin
	)

	var (
		_ = enum.Map(RoleUser, "user")
		_ = enum.Map(RoleAdmin, "admin")
		_ = enum.Finalize[Role]()
	)

	assert.PanicsWithValue(t, "enum Role: the enum was already finalized", func() { enum.New[Role]("moderator") })
}

func TestEnumMapDuplicated(t *testing.T) {
	type Role int

	var (
		RoleUser = enum.New[Role]("user")
	)

	assert.Equal(t, enum.ToString(RoleUser), "user")
	assert.PanicsWithValue(t,
		"enum Role (0): do not map number twice",
		func() { enum.Map(RoleUser, "admin") },
	)
	assert.PanicsWithValue(t,
		"enum Role (1): string user was already mapped to 0",
		func() { enum.Map(Role(1), "user") },
	)
}

func TestEnumMustToString(t *testing.T) {
	type Role int

	assert.Equal(t, enum.ToString(Role(42)), "<nil>")

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, enum.ToString(RoleUser), "user")
	assert.Equal(t, enum.ToString(RoleAdmin), "admin")
	assert.Equal(t, enum.ToString(Role(42)), "<nil>")
}

func TestEnumFromString(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	userRole, _ := enum.FromString[Role]("user")
	assert.Equal(t, userRole, RoleUser)
	adminRole, _ := enum.FromString[Role]("admin")
	assert.Equal(t, adminRole, RoleAdmin)
	_, valid := enum.FromString[Role]("moderator")
	assert.False(t, valid)
}

func TestEnumMustFromString(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, enum.MustFromString[Role]("user"), RoleUser)
	assert.Equal(t, enum.MustFromString[Role]("admin"), RoleAdmin)
}

func TestEnumFromNumber(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	userRole, ok := enum.FromNumber[Role](0)
	assert.True(t, ok)
	assert.Equal(t, userRole, RoleUser)

	adminRole, ok := enum.FromNumber[Role](1)
	assert.True(t, ok)
	assert.Equal(t, adminRole, RoleAdmin)

	_, ok = enum.FromString[Role]("moderator")
	assert.False(t, ok)
}

func TestEnumMustFromNumber(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, enum.MustFromNumber[Role](0), RoleUser)
	assert.Equal(t, enum.MustFromNumber[Role](1), RoleAdmin)
}

func TestEnumUndefined(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.True(t, enum.IsValid(RoleUser))
	assert.True(t, enum.IsValid(RoleAdmin))

	_, ok := enum.FromString[Role]("moderator")
	assert.False(t, ok)
}

func TestEnumUndefinedEnum(t *testing.T) {
	type Role int

	moderator, _ := enum.FromString[Role]("moderator")
	assert.False(t, enum.IsValid(moderator))
	assert.False(t, enum.IsValid(Role(0)))
}

func TestEnumMarshalJSON(t *testing.T) {
	type Role int

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := enum.MarshalJSON(RoleUser)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"user"`), data)

	_, err = enum.MarshalJSON(Role(1))
	assert.ErrorContains(t, err, "enum Role: invalid")
}

func TestEnumUnmarshalJSON(t *testing.T) {
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
	assert.ErrorContains(t, err, "invalid string")

	// Invalid enum
	err = enum.UnmarshalJSON([]byte(`"admin"`), &data)
	assert.ErrorContains(t, err, "enum Role: unknown string admin")
}

func TestEnumAll(t *testing.T) {
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

func TestEnumByte(t *testing.T) {
	type Role byte

	assert.Nil(t, enum.All[Role]())

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, []Role{RoleUser, RoleAdmin}, enum.All[Role]())
}

func TestEnumFloat32(t *testing.T) {
	type Role float32

	assert.Nil(t, enum.All[Role]())

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, []Role{RoleUser, RoleAdmin}, enum.All[Role]())
}

func TestEnumFloat64(t *testing.T) {
	type Role float64

	assert.Nil(t, enum.All[Role]())

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, []Role{RoleUser, RoleAdmin}, enum.All[Role]())
}

func TestEnumFloat32Map(t *testing.T) {
	type Role float32

	const (
		RoleUser  Role = 1.13
		RoleAdmin Role = 3.14
	)

	var (
		_ = enum.Map(RoleUser, "user")
		_ = enum.Map(RoleAdmin, "admin")
	)

	role, ok := enum.FromNumber[Role](float32(1.13))
	assert.True(t, ok)
	assert.Equal(t, RoleUser, role)

	role, ok = enum.FromNumber[Role](float32(3.14))
	assert.True(t, ok)
	assert.Equal(t, RoleAdmin, role)
}

func TestEnumFloat64Map(t *testing.T) {
	type Role float64

	const (
		RoleUser  Role = 1.13
		RoleAdmin Role = 3.14
	)

	var (
		_ = enum.Map(RoleUser, "user")
		_ = enum.Map(RoleAdmin, "admin")
	)

	role, ok := enum.FromNumber[Role](float64(1.13))
	assert.True(t, ok)
	assert.Equal(t, RoleUser, role)

	role, ok = enum.FromNumber[Role](float64(3.14))
	assert.True(t, ok)
	assert.Equal(t, RoleAdmin, role)
}

func TestEnumFloat32LikeInt(t *testing.T) {
	type Role float32

	assert.Nil(t, enum.All[Role]())

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, []Role{RoleUser, RoleAdmin}, enum.All[Role]())

	role, ok := enum.FromNumber[Role](0)
	assert.True(t, ok)
	assert.Equal(t, RoleUser, role)

	role, ok = enum.FromNumber[Role](1)
	assert.True(t, ok)
	assert.Equal(t, RoleAdmin, role)
}

func TestEnumFloat64LikeInt(t *testing.T) {
	type Role float32

	assert.Nil(t, enum.All[Role]())

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, []Role{RoleUser, RoleAdmin}, enum.All[Role]())

	role, ok := enum.FromNumber[Role](0)
	assert.True(t, ok)
	assert.Equal(t, RoleUser, role)

	role, ok = enum.FromNumber[Role](1)
	assert.True(t, ok)
	assert.Equal(t, RoleAdmin, role)
}

func TestEnumIntFromFloat64(t *testing.T) {
	type Role int

	assert.Nil(t, enum.All[Role]())

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, []Role{RoleUser, RoleAdmin}, enum.All[Role]())

	role, ok := enum.FromNumber[Role](float64(0))
	assert.True(t, ok)
	assert.Equal(t, RoleUser, role)

	role, ok = enum.FromNumber[Role](float64(1))
	assert.True(t, ok)
	assert.Equal(t, RoleAdmin, role)
}

func TestEnumNameOf(t *testing.T) {
	type Role int
	assert.Equal(t, "Role", enum.NameOf[Role]())
	assert.Equal(t, "Role", enum.NameOf[Role]())

	type weekday any
	type Weekday = enum.WrapEnum[weekday]
	assert.Equal(t, "Weekday", enum.NameOf[Weekday]())

	type someURL any
	type SomeURL = enum.WrapEnum[someURL]
	assert.Equal(t, "SomeURL", enum.NameOf[SomeURL]())
}

func TestEnumValueSQL(t *testing.T) {
	type Role int

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := enum.ValueSQL(RoleUser)
	assert.NoError(t, err)
	assert.Equal(t, "user", data)

	_, err = enum.ValueSQL(Role(1))
	assert.ErrorContains(t, err, "enum Role: invalid value 1")
}

func TestEnumScanSQL(t *testing.T) {
	type Role int

	var (
		RoleUser = enum.New[Role]("user")
	)

	var data Role

	// Scan bytes
	err := enum.ScanSQL([]byte(`user`), &data)
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	// Scan string
	err = enum.ScanSQL("user", &data)
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, data)

	// Invalid enum
	err = enum.ScanSQL("admin", &data)
	assert.ErrorContains(t, err, "enum Role: unknown string admin")
}

func TestEnumSQL(t *testing.T) {
	type role any
	type Role = enum.WrapEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	// Open an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	// Create a table for storing enum values
	_, err = db.Exec(`CREATE TABLE my_table (
		id INTEGER PRIMARY KEY,
		role TEXT
	);`)
	assert.NoError(t, err)

	_, err = db.Exec(`INSERT INTO my_table (role) VALUES (?)`, RoleUser)
	assert.NoError(t, err)

	// Retrieve the enum value from the table
	var retrievedRole Role
	err = db.QueryRow(`SELECT role FROM my_table WHERE id = 1`).Scan(&retrievedRole)
	assert.NoError(t, err)

	// Check if the deserialized value matches the expected value
	assert.Equal(t, retrievedRole, RoleUser)
}

func TestEnumJSON(t *testing.T) {
	type role any
	type Role = enum.WrapEnum[role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	type TestJSON struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Role Role   `json:"role"`
	}

	s := TestJSON{
		ID:   1,
		Name: "tester",
		Role: RoleUser,
	}

	data, err := json.Marshal(s)
	assert.NoError(t, err)
	assert.Equal(t, "{\"id\":1,\"name\":\"tester\",\"role\":\"user\"}", string(data))

	err = json.Unmarshal([]byte("{\"id\":1,\"name\":\"tester\",\"role\":\"user\"}"), &s)
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, s.Role)
}

func TestEnumPrintZeroStruct(t *testing.T) {
	type Role int

	var (
		_ = enum.New[Role]("user")
	)

	type User struct {
		Role Role
	}

	assert.Equal(t, "{0}", fmt.Sprint(User{}))
}

func TestNewExtended(t *testing.T) {
	type Role struct{ enum.SafeEnum[int] }

	var (
		RoleUser  = enum.NewExtended[Role]("user")
		RoleAdmin = enum.NewExtended[Role]("admin")
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
		_ = enum.New[Role]("user")
	)

	type User struct {
		Role Role
	}

	assert.Equal(t, "{<nil>}", fmt.Sprint(User{}))
}

func TestWrapEnumPrintZeroStruct(t *testing.T) {
	type role any
	type Role = enum.WrapEnum[role]

	var (
		_ = enum.New[Role]("user")
	)

	type User struct {
		Role Role
	}

	assert.Equal(t, "{user}", fmt.Sprint(User{}))
}

func TestEnumMarshalJSONInvalid(t *testing.T) {
	type Role int

	_, err := enum.MarshalJSON(Role(0))
	assert.ErrorContains(t, err, "enum Role: invalid value 0")
}

func TestWrapEnumMarshalJSONInvalid(t *testing.T) {
	type role int
	type Role = enum.WrapEnum[role]

	_, err := json.Marshal(Role(0))
	assert.ErrorContains(t, err, "enum WrapEnum[role]: invalid value 0")
}

func TestSafeEnumMarshalJSONInvalid(t *testing.T) {
	type role int
	type Role = enum.SafeEnum[role]

	_, err := json.Marshal(Role{})
	assert.ErrorContains(t, err, "enum SafeEnum[role]: invalid value <nil>")
}

func TestEnumUnmarshalJSONInvalid(t *testing.T) {
	type Role int

	var r Role
	err := enum.UnmarshalJSON([]byte(`"invalid"`), &r)
	assert.ErrorContains(t, err, "enum Role: unknown string invalid")
}

func TestWrapEnumUnmarshalJSONInvalid(t *testing.T) {
	type role int
	type Role = enum.WrapEnum[role]

	var r Role
	err := json.Unmarshal([]byte(`"invalid"`), &r)
	assert.ErrorContains(t, err, "enum WrapEnum[role]: unknown string invalid")
}

func TestSafeEnumUnmarshalJSONInvalid(t *testing.T) {
	type role int
	type Role = enum.SafeEnum[role]

	var r Role
	err := json.Unmarshal([]byte(`"invalid"`), &r)
	assert.ErrorContains(t, err, "enum SafeEnum[role]: unknown string invalid")
}

func TestEnumValueSQLInvalid(t *testing.T) {
	type Role int

	_, err := enum.ValueSQL(Role(0))
	assert.ErrorContains(t, err, "enum Role: invalid value 0")
}

func TestWrapEnumValueSQLInvalid(t *testing.T) {
	type role int
	type Role = enum.WrapEnum[role]

	_, err := enum.ValueSQL(Role(0))
	assert.ErrorContains(t, err, "enum WrapEnum[role]: invalid value 0")
}

func TestSafeEnumValueSQLInvalid(t *testing.T) {
	type role int
	type Role = enum.SafeEnum[role]

	_, err := enum.ValueSQL(Role{})
	assert.ErrorContains(t, err, "enum SafeEnum[role]: invalid value <nil>")
}

func TestEnumScanSQLInvalid(t *testing.T) {
	type Role int

	var r Role
	err := enum.ScanSQL([]byte("invalid"), &r)
	assert.ErrorContains(t, err, "enum Role: unknown string invalid")
}

func TestWrapEnumScanSQLInvalid(t *testing.T) {
	type role int
	type Role = enum.WrapEnum[role]

	var r Role
	err := enum.ScanSQL([]byte("invalid"), &r)
	assert.ErrorContains(t, err, "enum WrapEnum[role]: unknown string invalid")
}

func TestSafeEnumScanSQLInvalid(t *testing.T) {
	type role int
	type Role = enum.SafeEnum[role]

	var r Role
	err := enum.ScanSQL([]byte("invalid"), &r)
	assert.ErrorContains(t, err, "enum SafeEnum[role]: unknown string invalid")
}
