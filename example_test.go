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
		_         = enum.Finalize[Role]() // Optional: ensure no new enum values can be added to Role.
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
		_ = enum.Finalize[Role]() // Optional: ensure no new enum values can be added to Role.
	)

	// Print all enum values.
	for _, e := range enum.All[Role]() {
		fmt.Println(enum.ToString(e))
	}

	// Output:
	// user
	// admin
}

func ExampleSerde() {
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
		_ = enum.Finalize[Role]() // Optional: ensure no new enum values can be added to Role.
	)

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`

		// Serde provides functionality for serializing and deserializing enums
		// that cannot be directly serialized or deserialized.
		Role enum.Serde[Role] `json:"role"`
	}

	user := User{ID: 1, Name: "tester", Role: enum.SerdeWrap(RoleAdmin)}
	data, _ := json.Marshal(user)
	fmt.Println(string(data)) // role's value should be "admin" instead of 1.

	// Output:
	// {"id":1,"name":"tester","role":"admin"}
}

func ExampleIntEnum() {
	// Define a generic enum type
	type role any
	type Role = enum.IntEnum[role]

	// Define enum values for Role using iota
	const (
		RoleUser Role = iota
		RoleAdmin
	)

	// Map string representations to enum values
	var (
		_ = enum.Map(RoleUser, "user")
		_ = enum.Map(RoleAdmin, "admin")
		_ = enum.Finalize[Role]() // Optional: ensure no new enum values can be added to Role.
	)

	// As Role is a IntEnum, it can utilize methods from IntEnum, including
	// utility functions and serde operations.
	fmt.Println(RoleUser.GoString()) // 0 (user)
	fmt.Println(RoleUser.IsValid())  // true

	// Define a struct that includes the Role enum.
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Role Role   `json:"role"`
	}

	// Serialize the User struct to JSON.
	user1 := User{ID: 0, Name: "tester", Role: RoleAdmin}
	data, _ := json.Marshal(user1)
	fmt.Println(string(data))

	// Deserialize JSON back into a User struct and print the Role.
	user2 := User{}
	json.Unmarshal(data, &user2)
	fmt.Println(user2.Role)

	// Output:
	// 0 (user)
	// true
	// {"id":0,"name":"tester","role":"admin"}
	// admin
}

func ExampleStructEnum() {
	// Define a generic enum type
	type underlyingRole string
	type Role = enum.StructEnum[underlyingRole]

	// Define enum values for Role using iota
	var (
		RoleUser  = enum.NewStruct[underlyingRole]("user")
		RoleAdmin = enum.NewStruct[underlyingRole]("admin")
		_         = enum.Finalize[Role]() // Optional: ensure no new enum values can be added to Role.
	)

	// As Role is a StructEnum, it can utilize methods from StructEnum, including
	// utility functions and serde operations.
	fmt.Println(RoleUser.GoString()) // 0 (user)
	fmt.Println(RoleUser.IsValid())  // true

	// Define a struct that includes the Role enum.
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Role Role   `json:"role"`
	}

	// Serialize the User struct to JSON.
	user1 := User{ID: 0, Name: "tester", Role: RoleAdmin}
	data, _ := json.Marshal(user1)
	fmt.Println(string(data))

	// Deserialize JSON back into a User struct and print the Role.
	user2 := User{}
	json.Unmarshal(data, &user2)
	fmt.Println(user2.Role)

	// Output:
	// 0 (user)
	// true
	// {"id":0,"name":"tester","role":"admin"}
	// admin
}
