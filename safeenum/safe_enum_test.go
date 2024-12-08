package safeenum_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xybor-x/enum"
	"github.com/xybor-x/enum/safeenum"
)

type role any
type Role = safeenum.SafeEnum[role]

var (
	RoleUser  = safeenum.New[role, safeenum.P0]("user")
	RoleAdmin = safeenum.New[role, safeenum.P1]("admin")
)

func TestString(t *testing.T) {
	assert.Equal(t, "user", RoleUser.String())
	assert.Equal(t, "admin", RoleAdmin.String())
}

func TestAll(t *testing.T) {
	assert.Equal(t, enum.All[Role](), []Role{RoleUser, RoleAdmin})
}

func TestMarshal(t *testing.T) {
	data, err := json.Marshal(RoleUser)
	assert.NoError(t, err)
	assert.Equal(t, `"user"`, string(data))
}

func TestValue(t *testing.T) {
	data, err := RoleUser.Value()
	assert.NoError(t, err)
	assert.Equal(t, "user", data)
}

func TestInt(t *testing.T) {
	assert.Equal(t, 0, RoleUser.Int())
	assert.Equal(t, 1, RoleAdmin.Int())
}

func TestGoString(t *testing.T) {
	assert.Equal(t, "0 (user)", RoleUser.GoString())
}

func TestUnmarshal(t *testing.T) {
	var r enum.Serde[Role]
	err := json.Unmarshal([]byte(`"user"`), &r)
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, r.Enum())
}

func TestScan(t *testing.T) {
	var r enum.Serde[Role]

	assert.NoError(t, r.Scan("user"))
	assert.Equal(t, RoleUser, r.Enum())

	assert.Error(t, r.Scan("nothing"))
}

func TestUnmarshalJSONStruct(t *testing.T) {
	data := "{\"id\":1,\"role\":\"user\"}"

	type User1 struct {
		ID   int  `json:"id"`
		Role Role `json:"role"`
	}

	var user1 User1
	assert.ErrorContains(t, json.Unmarshal([]byte(data), &user1),
		"json: cannot unmarshal string into Go struct field User1.role")

	type User2 struct {
		ID   int              `json:"id"`
		Role enum.Serde[Role] `json:"role"`
	}

	var user2 User2
	assert.NoError(t, json.Unmarshal([]byte(data), &user2))
	assert.Equal(t, RoleUser, user2.Role.Enum())
}
