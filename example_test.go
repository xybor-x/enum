package enum_test

import (
	"encoding/json"
	"fmt"

	"github.com/xybor-x/enum"
)

func ExampleNew() {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
	)

	fmt.Println("string repr of RoleUser:", enum.StringOf(RoleUser))
	fmt.Println("string repr of RoleAdmin:", enum.StringOf(RoleAdmin))

	fmt.Println("number repr of \"user\":", enum.EnumOf[Role]("user"))
	fmt.Println("number repr of \"admin\":", enum.EnumOf[Role]("admin"))

	// Output:
	// string repr of RoleUser: user
	// string repr of RoleAdmin: admin
	// number repr of "user": 0
	// number repr of "admin": 1
}

func ExampleMap() {
	type Role int

	const (
		RoleUser Role = iota
		RoleAdmin
	)

	var (
		_ = enum.Map(RoleUser, "user")
		_ = enum.Map(RoleAdmin, "admin")
	)

	fmt.Println("string repr of RoleUser:", enum.StringOf(RoleUser))
	fmt.Println("string repr of RoleAdmin:", enum.StringOf(RoleAdmin))

	fmt.Println("number repr of \"user\":", enum.EnumOf[Role]("user"))
	fmt.Println("number repr of \"admin\":", enum.EnumOf[Role]("admin"))

	// Output:
	// string repr of RoleUser: user
	// string repr of RoleAdmin: admin
	// number repr of "user": 0
	// number repr of "admin": 1
}

func ExampleRichEnum() {
	type unsafeRole any
	type Role = enum.RichEnum[unsafeRole]

	const (
		RoleUser Role = iota
		RoleAdmin
	)

	var (
		_ = enum.Map(RoleUser, "user")   // Maps RoleUser to "user"
		_ = enum.Map(RoleAdmin, "admin") // Maps RoleAdmin to "admin"
	)

	fmt.Println("string repr of RoleUser:", RoleUser.String())
	fmt.Println("string repr of RoleAdmin:", RoleAdmin.String())

	fmt.Println("number repr of \"user\":", enum.EnumOf[Role]("user").Int())
	fmt.Println("number repr of \"admin\":", enum.EnumOf[Role]("admin").Int())

	data, err := json.Marshal(RoleUser)
	if err != nil {
		panic(err)
	}
	fmt.Println("marshal RoleUser:", string(data))

	var r Role
	if err := json.Unmarshal([]byte(`"user"`), &r); err != nil {
		panic(err)
	}
	fmt.Println("unmarshal \"user\":", r.Int())

	// Output:
	// string repr of RoleUser: user
	// string repr of RoleAdmin: admin
	// number repr of "user": 0
	// number repr of "admin": 1
	// marshal RoleUser: "user"
	// unmarshal "user": 0
}
