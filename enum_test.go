package enum_test

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xybor-x/enum"

	_ "github.com/mattn/go-sqlite3"
)

func TestNew(t *testing.T) {
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

func TestMap(t *testing.T) {
	type Role int
	const (
		RoleUser Role = iota
		RoleAdmin
	)

	assert.Equal(t, enum.ToString(RoleUser), "Role::<<undefined>>")
	assert.Equal(t, enum.ToString(RoleAdmin), "Role::<<undefined>>")

	var (
		_ = enum.Map(RoleUser, "user")
		_ = enum.Map(RoleAdmin, "admin")
	)

	assert.Equal(t, enum.ToString(RoleUser), "user")
	assert.Equal(t, enum.ToString(RoleAdmin), "admin")
}

func TestMapDuplicated(t *testing.T) {
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

func TestToString(t *testing.T) {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, enum.ToString(RoleUser), "user")
	assert.Equal(t, enum.ToString(RoleAdmin), "admin")
	assert.Equal(t, enum.ToString(Role(2)), "Role::<<undefined>>")
}

func TestMustToString(t *testing.T) {
	type Role int

	assert.Panics(t, func() { enum.MustToString(Role(0)) })
	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	assert.Equal(t, enum.MustToString(RoleUser), "user")
	assert.Equal(t, enum.MustToString(RoleAdmin), "admin")
	assert.Panics(t, func() { enum.MustToString(Role(2)) })
}

func TestFromString(t *testing.T) {
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

func TestMustFromString(t *testing.T) {
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

func TestFromInt(t *testing.T) {
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

func TestMustFromInt(t *testing.T) {
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

func TestUndefined(t *testing.T) {
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

func TestUndefinedEnum(t *testing.T) {
	type Role int

	moderator, _ := enum.FromString[Role]("moderator")
	assert.False(t, enum.IsValid(moderator))
	assert.False(t, enum.IsValid(Role(0)))
}

func TestMarshalJSON(t *testing.T) {
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

func TestUnmarshalJSON(t *testing.T) {
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

func TestAll(t *testing.T) {
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

func TestNonIntEnum(t *testing.T) {
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

func TestValueSQL(t *testing.T) {
	type Role int

	var (
		RoleUser = enum.New[Role]("user")
	)

	data, err := enum.ValueSQL(RoleUser)
	assert.NoError(t, err)
	assert.Equal(t, "user", data)

	_, err = enum.ValueSQL(Role(1))
	assert.ErrorContains(t, err, "unknown Role: 1")
}

func TestScanSQL(t *testing.T) {
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
	assert.ErrorContains(t, err, "unknown Role string: admin")
}

func TestSQL(t *testing.T) {
	type role any
	type Role = enum.RichEnum[role]

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

func TestJSON(t *testing.T) {
	type role any
	type Role = enum.RichEnum[role]

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
