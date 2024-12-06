package enum_test

import (
	"fmt"

	"github.com/xybor-x/enum"
)

func ExampleNew() {
	type Role int

	var (
		RoleUser  = enum.New[Role]("user")  // Dynamically creates and maps "user"
		RoleAdmin = enum.New[Role]("admin") // Dynamically creates and maps "admin"
	)

	fmt.Println("string repr of", RoleUser, "is", enum.StringOf(RoleUser))
	fmt.Println("string repr of", RoleAdmin, "is", enum.StringOf(RoleAdmin))

	fmt.Println("number repr of user is", enum.EnumOf[Role]("user"))
	fmt.Println("number repr of admin is", enum.EnumOf[Role]("admin"))

	// Output:
	// string repr of 0 is user
	// string repr of 1 is admin
	// number repr of user is 0
	// number repr of admin is 1
}

func ExampleMap() {
	type Role int

	const (
		RoleUser Role = iota
		RoleAdmin
	)

	var (
		_ = enum.Map(RoleUser, "user")   // Maps RoleUser to "user"
		_ = enum.Map(RoleAdmin, "admin") // Maps RoleAdmin to "admin"
	)

	fmt.Println("string repr of", RoleUser, "is", enum.StringOf(RoleUser))
	fmt.Println("string repr of", RoleAdmin, "is", enum.StringOf(RoleAdmin))

	fmt.Println("number repr of user is", enum.EnumOf[Role]("user"))
	fmt.Println("number repr of admin is", enum.EnumOf[Role]("admin"))

	// Output:
	// string repr of 0 is user
	// string repr of 1 is admin
	// number repr of user is 0
	// number repr of admin is 1
}
