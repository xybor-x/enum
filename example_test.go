package enum_test

import (
	"encoding/json"
	"fmt"

	"github.com/xybor-x/enum"
)

func ExampleNew() {
	type Role int

	// Define enum values for Role.
	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
		_         = enum.Finalize[Role]()
	)

	// Print the string representation of enum values.
	fmt.Println("string(RoleUser):", enum.ToString(RoleUser))
	fmt.Println("string(RoleAdmin):", enum.ToString(RoleAdmin))

	// Demonstrate converting string to enum value.
	fmt.Println(`enum("user"):`, enum.MustFromString[Role]("user"))
	fmt.Println(`enum("admin"):`, enum.MustFromString[Role]("admin"))

	// Output:
	// string(RoleUser): user
	// string(RoleAdmin): admin
	// enum("user"): 0
	// enum("admin"): 1
}

func ExampleMap() {
	type Role int

	// Define enum values for Role using iota.
	const (
		RoleUser Role = iota
		RoleAdmin
	)

	// Map string representations to enum values.
	var (
		_ = enum.Map(RoleUser, "user")
		_ = enum.Map(RoleAdmin, "admin")
		_ = enum.Finalize[Role]()
	)

	// Print all enum values.
	for _, e := range enum.All[Role]() {
		fmt.Println(enum.ToString(e))
	}

	// Output:
	// user
	// admin
}

func ExampleWrapEnum() {
	type role any
	type Role = enum.WrapEnum[role]

	// Define enum values for Role using iota
	const (
		RoleUser Role = iota
		RoleAdmin
	)

	// Map string representations to enum values
	var (
		_ = enum.Map(RoleUser, "user")
		_ = enum.Map(RoleAdmin, "admin")
		_ = enum.Finalize[Role]()
	)

	// Define a struct that includes the Role enum.
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Role Role   `json:"role"`
	}

	// Serialize zero struct
	data, _ := json.Marshal(User{})
	fmt.Println(string(data))

	// Serialize the User struct to JSON.
	user1 := User{ID: 0, Name: "tester", Role: RoleAdmin}
	data, _ = json.Marshal(user1)
	fmt.Println(string(data))

	// Deserialize JSON back into a User struct and print the Role.
	user2 := User{}
	json.Unmarshal(data, &user2)
	fmt.Println(user2.Role)

	// Output:
	// {"id":0,"name":"","role":"user"}
	// {"id":0,"name":"tester","role":"admin"}
	// admin
}

func ExampleSafeEnum() {
	type role string
	type Role = enum.SafeEnum[role]

	var (
		RoleUser  = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
		_         = enum.Finalize[Role]()
	)

	// As Role is a SafeEnum, it can utilize methods from SafeEnum, including
	// utility functions and serde operations.
	fmt.Println(RoleUser.GoString()) // 0 (user)
	fmt.Println(RoleUser.IsValid())  // true

	// Define a struct that includes the Role enum.
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Role Role   `json:"role"`
	}

	// Serialize zero struct
	data, err := json.Marshal(User{})
	if err != nil {
		fmt.Println("cannot serialize zero struct")
	}

	// Serialize the User struct to JSON.
	user1 := User{ID: 0, Name: "tester", Role: RoleAdmin}
	data, _ = json.Marshal(user1)
	fmt.Println(string(data))

	// Deserialize JSON back into a User struct and print the Role.
	user2 := User{}
	json.Unmarshal(data, &user2)
	fmt.Println(user2.Role)

	// Output:
	// 0 (user)
	// true
	// cannot serialize zero struct
	// {"id":0,"name":"tester","role":"admin"}
	// admin
}

func ExampleNullable() {
	type Role int
	type NullRole = enum.Nullable[Role]

	var (
		_         = enum.New[Role]("user")
		RoleAdmin = enum.New[Role]("admin")
		_         = enum.Finalize[Role]()
	)

	// Define a struct that includes the Role enum.
	type User struct {
		ID   int      `json:"id"`
		Name string   `json:"name"`
		Role NullRole `json:"role"`
	}

	// Serialize zero struct
	data, _ := json.Marshal(User{})
	fmt.Println(string(data))

	// Serialize the User struct to JSON.
	user1 := User{ID: 0, Name: "tester", Role: NullRole{Enum: RoleAdmin, Valid: true}}
	data, _ = json.Marshal(user1)
	fmt.Println(string(data))

	// Output:
	// {"id":0,"name":"","role":null}
	// {"id":0,"name":"tester","role":"admin"}
}
