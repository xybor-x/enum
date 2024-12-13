package exhaustive_test

import (
	"fmt"

	"github.com/xybor-x/enum"
	"github.com/xybor-x/enum/exhaustive"
)

type role any
type Role struct{ enum.SafeEnum[role] }

var (
	RoleUser  = enum.NewExtended[Role]("user")
	RoleAdmin = enum.NewExtended[Role]("admin")
	_         = exhaustive.CheckMethodOf[Role]()
)

type (
	CaseRoleUser  exhaustive.Case
	CaseRoleAdmin exhaustive.Case
)

func (r Role) Switch(c1 CaseRoleUser, c2 CaseRoleAdmin) exhaustive.SwitchDefault {
	return exhaustive.Switch(r, c1, c2)
}

func ExampleSwitch() {
	role := RoleAdmin

	role.Switch(
		CaseRoleUser{func() { fmt.Println("case user") }},
		CaseRoleAdmin{func() { fmt.Println("case admin") }},
	).ByDefault(func() {
		panic("invalid role")
	})

	// Output:
	// case admin
}
