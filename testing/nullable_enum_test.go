package testing_test

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xybor-x/enum"
)

func TestNullableJSON(t *testing.T) {
	type role any
	type Role = enum.WrapEnum[role]
	type NullRole = enum.Nullable[Role]

	var (
		RoleUser = enum.New[Role]("user")
	)

	type TestJSON struct {
		ID   int      `json:"id"`
		Name string   `json:"name"`
		Role NullRole `json:"role"`
	}

	s := TestJSON{
		ID:   1,
		Name: "tester",
		Role: NullRole{Enum: RoleUser, Valid: true},
	}

	data, err := json.Marshal(s)
	assert.NoError(t, err)
	assert.Equal(t, "{\"id\":1,\"name\":\"tester\",\"role\":\"user\"}", string(data))

	err = json.Unmarshal([]byte("{\"id\":1,\"name\":\"tester\",\"role\":\"user\"}"), &s)
	assert.NoError(t, err)
	assert.True(t, s.Role.Valid)
	assert.Equal(t, RoleUser, s.Role.Enum)
}

func TestNullableJSONNull(t *testing.T) {
	type role any
	type Role = enum.WrapEnum[role]
	type NullRole = enum.Nullable[Role]

	var (
		_ = enum.New[Role]("user")
	)

	type TestJSON struct {
		ID   int      `json:"id"`
		Name string   `json:"name"`
		Role NullRole `json:"role"`
	}

	s := TestJSON{
		ID:   1,
		Name: "tester",
		Role: NullRole{},
	}

	data, err := json.Marshal(s)
	assert.NoError(t, err)
	assert.Equal(t, "{\"id\":1,\"name\":\"tester\",\"role\":null}", string(data))

	err = json.Unmarshal([]byte("{\"id\":1,\"name\":\"tester\",\"role\":null}"), &s)
	assert.NoError(t, err)
	assert.False(t, s.Role.Valid)
}

func TestNullableSQL(t *testing.T) {
	type role any
	type Role = enum.WrapEnum[role]
	type NullRole = enum.Nullable[Role]

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

	_, err = db.Exec(`INSERT INTO my_table (role) VALUES (?)`, NullRole{Enum: RoleUser, Valid: true})
	assert.NoError(t, err)

	// Retrieve the enum value from the table
	var retrievedRole NullRole
	err = db.QueryRow(`SELECT role FROM my_table WHERE id = 1`).Scan(&retrievedRole)
	assert.NoError(t, err)

	// Check if the deserialized value matches the expected value
	assert.True(t, retrievedRole.Valid)
	assert.Equal(t, retrievedRole.Enum, RoleUser)
}

func TestNullableSQLNull(t *testing.T) {
	type role any
	type Role = enum.WrapEnum[role]
	type NullRole = enum.Nullable[Role]

	var (
		_ = enum.New[Role]("user")
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

	_, err = db.Exec(`INSERT INTO my_table (role) VALUES (?)`, NullRole{})
	assert.NoError(t, err)

	// Retrieve the enum value from the table
	var retrievedRole NullRole
	err = db.QueryRow(`SELECT role FROM my_table WHERE id = 1`).Scan(&retrievedRole)
	assert.NoError(t, err)

	// Check if the deserialized value matches the expected value
	assert.False(t, retrievedRole.Valid)
}
