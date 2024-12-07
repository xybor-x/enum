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

	fmt.Println("string repr of RoleUser:", enum.ToString(RoleUser))
	fmt.Println("string repr of RoleAdmin:", enum.ToString(RoleAdmin))

	fmt.Println("number repr of user:", enum.MustFromString[Role]("user"))
	fmt.Println("number repr of admin:", enum.MustFromString[Role]("admin"))

	// Output:
	// string repr of RoleUser: user
	// string repr of RoleAdmin: admin
	// number repr of user: 0
	// number repr of admin: 1
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

	for _, e := range enum.All[Role]() {
		fmt.Println(enum.ToString(e))
	}

	// Output:
	// user
	// admin
}

func ExampleRichEnum() {
	type role any
	type Role = enum.RichEnum[role]

	const (
		RoleUser Role = iota
		RoleAdmin
	)

	var (
		_ = enum.Map(RoleUser, "user")
		_ = enum.Map(RoleAdmin, "admin")
	)

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Role Role   `json:"role"`
	}

	user1 := User{ID: 0, Name: "tester", Role: RoleAdmin}
	data, _ := json.Marshal(user1)
	fmt.Println(string(data))

	user2 := User{}
	json.Unmarshal(data, &user2)
	fmt.Printf("%#v", user2.Role)

	// Output:
	// {"id":0,"name":"tester","role":"admin"}
	// 1 (admin)
}
