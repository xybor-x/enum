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

	num, _ := enum.EnumOf[Role]("user")
	fmt.Println("number repr of \"user\":", num)
	num, _ = enum.EnumOf[Role]("admin")
	fmt.Println("number repr of \"admin\":", num)

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

	num, _ := enum.EnumOf[Role]("user")
	fmt.Println("number repr of \"user\":", num)
	num, _ = enum.EnumOf[Role]("admin")
	fmt.Println("number repr of \"admin\":", num)

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
		_ = enum.Map(RoleUser, "user")
		_ = enum.Map(RoleAdmin, "admin")
	)

	data, _ := json.Marshal(RoleUser)
	fmt.Println(string(data))
	fmt.Printf("%d\n", RoleAdmin)
	fmt.Printf("%s\n", RoleAdmin)
	fmt.Println(RoleAdmin.IsValid())

	// Output:
	// "user"
	// 1
	// admin
	// true
}
