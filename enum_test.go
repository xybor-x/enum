package enum_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
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

	type File int
	var (
		FileImage = enum.New[File]("image")
		FilePDF   = enum.New[File]("pdf")
	)

	assert.Equal(t, FileImage, File(0))
	assert.Equal(t, FilePDF, File(1))

	assert.NotEqual(t, File(0), Role(0))
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

	assert.Panics(t, func() { enum.New[Role]("moderator") })
}

func TestEnumMapDuplicated(t *testing.T) {
	type Role int

	var (
		RoleUser = enum.New[Role]("user")
	)

	const (
		RoleAdmin Role = iota
	)

	assert.Equal(t, enum.ToString(RoleUser), "user")
	assert.Panics(t, func() { enum.Map(RoleAdmin, "admin") })
	assert.Panics(t, func() { enum.Map(Role(1), "user") })
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

	assert.Panics(t, func() { enum.MustFromString[Role]("moderator") })

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, enum.MustFromString[Role]("user"), RoleUser)
	assert.Equal(t, enum.MustFromString[Role]("admin"), RoleAdmin)
	assert.Panics(t, func() { enum.MustFromString[Role]("moderator") })
}

func TestEnumFromInt(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	userRole, ok := enum.FromInt[Role](0)
	assert.True(t, ok)
	assert.Equal(t, userRole, RoleUser)

	adminRole, ok := enum.FromInt[Role](1)
	assert.True(t, ok)
	assert.Equal(t, adminRole, RoleAdmin)

	_, ok = enum.FromString[Role]("moderator")
	assert.False(t, ok)
}

func TestEnumMustFromInt(t *testing.T) {
	type Role int

	assert.Panics(t, func() { enum.MustFromInt[Role](0) })

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, enum.MustFromInt[Role](0), RoleUser)
	assert.Equal(t, enum.MustFromInt[Role](1), RoleAdmin)
	assert.Panics(t, func() { enum.MustFromInt[Role](2) })
}

func TestEnumMustToInt(t *testing.T) {
	type Role int

	assert.Equal(t, enum.ToInt(Role(42)), math.MinInt32)

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, enum.ToInt(RoleUser), 0)
	assert.Equal(t, enum.ToInt(RoleAdmin), 1)
	assert.Equal(t, enum.ToInt(Role(42)), math.MinInt32)
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
	assert.ErrorContains(t, err, "invalid character")

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

func TestEnumNonIntEnum(t *testing.T) {
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

func TestEnumValueSQL(t *testing.T) {
	type Role int

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := enum.ValueSQL(RoleUser)
	assert.NoError(t, err)
	assert.Equal(t, "user", data)

	_, err = enum.ValueSQL(Role(1))
	assert.ErrorContains(t, err, "enum Role: invalid 1")
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

func TestBasicEnumPrintZeroStruct(t *testing.T) {
	type Role int

	var (
		_ = enum.New[Role]("user")
	)

	type User struct {
		Role Role
	}

	assert.Equal(t, "{0}", fmt.Sprint(User{}))
}

func TestWrapEnumPrintZeroStruct(t *testing.T) {
	type underlyingRole any
	type Role = enum.WrapEnum[underlyingRole]

	var (
		_ = enum.New[Role]("user")
	)

	type User struct {
		Role Role
	}

	assert.Equal(t, "{user}", fmt.Sprint(User{}))
}

func TestStructEnumPrintZeroStruct(t *testing.T) {
	type underlyingRole any
	type Role = enum.SafeEnum[underlyingRole]

	var (
		_ = enum.NewSafe[underlyingRole]("user")
	)

	type User struct {
		Role Role
	}

	assert.Equal(t, "{<nil>}", fmt.Sprint(User{}))
}
