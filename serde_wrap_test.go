package enum_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xybor-x/enum"
)

func TestSerdeInt(t *testing.T) {
	type Role int

	var (
		RoleUser = enum.New[Role]("user")
	)

	nonSerdeJSON := `{"id":1,"role":0}`
	serdeJSON := `{"id":1,"role":"user"}`

	type NonSerdeUser struct {
		ID   int  `json:"id"`
		Role Role `json:"role"`
	}

	user1 := NonSerdeUser{ID: 1, Role: RoleUser}
	data, err := json.Marshal(user1)
	assert.NoError(t, err)
	assert.Equal(t, nonSerdeJSON, string(data))

	type SerdeUser struct {
		ID   int              `json:"id"`
		Role enum.Serde[Role] `json:"role"`
	}

	user2 := SerdeUser{ID: 1, Role: enum.SerdeWrap(RoleUser)}
	data, err = json.Marshal(user2)
	assert.NoError(t, err)
	assert.Equal(t, serdeJSON, string(data))

	user3 := SerdeUser{}
	assert.ErrorContains(t, json.Unmarshal([]byte(nonSerdeJSON), &user3), "json: cannot unmarshal")

	assert.NoError(t, json.Unmarshal([]byte(serdeJSON), &user3))
	assert.Equal(t, RoleUser, user3.Role.Enum())
}
