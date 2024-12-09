package safeenum_test

import (
	"encoding/json"
	"fmt"

	"github.com/xybor-x/enum"
	"github.com/xybor-x/enum/safeenum"
)

func ExampleSafeEnum() {
	// Define enum's underlying type.
	type underlyingRole any

	// Create a SafeEnum type for roles.
	type Role = safeenum.SafeEnum[underlyingRole]

	// Define specific enum values for the Role type.
	// The second type parameter is known as the positioner. Note that each enum
	// must have a unique positioner; no two enums can share the same positioner.
	var (
		RoleUser  = safeenum.New[underlyingRole, safeenum.P0]("user")
		RoleAdmin = safeenum.New[underlyingRole, safeenum.P1]("admin")
		_         = enum.Finalize[Role]() // Optional: ensure no new enum values can be added to Role.
	)

	// Utility functions of enum can still be used with SafeEnum.
	fmt.Println(enum.All[Role]())                      // [user admin]
	fmt.Println(enum.MustFromInt[Role](0) == RoleUser) // true

	// Define a User struct that includes a Role field.
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`

		// Use enum.Serde due to the designation of SafeEnum, as it cannot be directly deserialized.
		Role enum.Serde[Role] `json:"role"`
	}

	// Create a new User with the RoleAdmin and serialize it to JSON.
	user1 := User{ID: 0, Name: "tester", Role: enum.SerdeWrap(RoleAdmin)}
	data, _ := json.Marshal(user1) // Marshal user1 into JSON
	fmt.Println(string(data))      // Print the JSON string representation of user1

	// Unmarshal the JSON data back into a User object.
	user2 := User{}
	json.Unmarshal(data, &user2)   // Deserialize the JSON data into user2
	fmt.Println(user2.Role.Enum()) // Print the role of user2, which should be "admin"

	// Output:
	// [user admin]
	// true
	// {"id":0,"name":"tester","role":"admin"}
	// admin
}
