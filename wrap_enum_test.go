package enum_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xybor-x/enum"
)

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
