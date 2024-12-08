[![Go Reference](https://pkg.go.dev/badge/github.com/xybor-x/enum.svg)](https://pkg.go.dev/github.com/xybor-x/enum)

# Enum

**Elegant and powerful enums for Go with zero code generation!**

`xybor-x/enum` makes working with enums in Go intuitive and efficient by providing:
- Seamless enum-string mappings.
- Constant enums compatible with Go's `iota` conventions.
- Out of the box JSON/SQL serialization and deserialization.
- Rich utilities for enhanced enum handling.


## Installation

Install the package via `go get`:
```sh
go get -u github.com/xybor-x/enum
```


## Quick start

*Please refer [Usage](#usage) for further details.*

In Go, `iota` is a special identifier used to create incrementing constants, making it perfect for defining enums.

`xybor-x/enum` is fully compatible with `iota`-based enums.

```go
type role any
type Role = enum.RichEnum[role]

const (
	RoleUser Role = iota
	RoleAdmin
)

func init() {
	enum.Map(RoleUser, "user")
	enum.Map(RoleAdmin, "admin")
}

func main() {
    data, _ := json.Marshal(RoleUser) // Output: "user"
    fmt.Println(RoleAdmin.IsValid())  // Output: true
}
```


## Usage

### Define enum

**Note**: Enum definitions are *NOT thread-safe*. Therefore, they should be finalized during initialization (at the global scope).

#### Short definition (variable enums)

```go
type Role int

var (
    RoleUser  = enum.New[Role]("user")
    RoleAdmin = enum.New[Role]("admin")
	_         = enum.Finalize[Role]() // Optional: ensure no new enum values can be added to Role.
)
```

#### Full definition (constant enums)

``` go
type Role int

const (
    RoleUser Role = iota
    RoleAdmin
)

func init() {
    enum.Map(RoleUser, "user")
    enum.Map(RoleAdmin, "admin")
	enum.Finalize[Role]() // Optional: ensure no new enum values can be added to Role.
}
```

### Utility functions

#### FromString

Convert a `string` value to an `enum` type and check if it's valid.

```go
role, ok := enum.FromString[Role]("user")
if ok {
    fmt.Println("Enum representation:", role) // Output: 0
} else {
    fmt.Println("Invalid enum")
}
```

#### FromInt

Convert an `int` value to an `enum` type and check if it's valid.

```go
role, ok := enum.FromInt[Role](42)
if ok {
    fmt.Println("Enum representation:", role)
} else {
    fmt.Println("Invalid enum") // Output: Invalid enum
}
```

#### ToString

Convert an `enum` to string.

```go
fmt.Println("String:", enum.ToString(RoleAdmin))  // Output: "admin"

// Note that you should check if enum is valid before calling ToString for
// an unsafe enum.
fmt.Println("String:", enum.ToString(Role(42)))   // panic
```

#### IsValid

```go
fmt.Println(enum.IsValid(RoleUser)) // true
fmt.Println(enum.IsValid(Role(0)))  // true
fmt.Println(enum.IsValid(Role(42)))  // false
```

#### All

```go
for _, role := range enum.All[Role]() {
    fmt.Println("Role:", enum.ToString(role))
}
// Output:
// Role: user
// Role: admin
```

### Rich enum

Instead of defining your enum type as `int`, you can use `enum.RichEnum` (also an `int` alias) to leverage several convenient features:
- Interact with enums using methods instead of standalone functions.
- Built-in support for serialization and deserialization (JSON and SQL).

```go
// Define enum's underlying type.
type role any

// Create a RichEnum type for roles.
type Role = enum.RichEnum[role] // NOTE: It must use type alias instead of type definition.

const (
    RoleUser Role = iota
    RoleAdmin
)

func init() {
    enum.Map(RoleUser, "user")
    enum.Map(RoleAdmin, "admin")
	enum.Finalize[Role]() // Optional: ensure no new enum values can be added to Role.
}

func main() {
    data, _ := json.Marshal(RoleUser) // Output: "user"
    fmt.Println(RoleAdmin.IsValid())  // Output: true
}
```

### Safe enum

`SafeEnum` defines a type-safe enum.

The `SafeEnum` enforces strict type safety, ensuring that only predefined enum values are allowed. It prevents the accidental creation of new enum types, providing a guaranteed set of valid values.

```go
// Define enum's underlying type.
type unsafeRole any

// Create a SafeEnum type for roles.
type Role = safeenum.SafeEnum[unsafeRole]

// Define specific enum values for the Role type.
// The second type parameter is known as the positioner. Note that each enum
// must have a unique positioner; no two enums can share the same positioner.
var (
    RoleUser  = safeenum.New[unsafeRole, safeenum.P0]("user")
    RoleAdmin = safeenum.New[unsafeRole, safeenum.P1]("admin")
    _         = enum.Finalize[Role]() // Optional: ensure no new enum values can be added to Role.
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    
    // Use enum.Serde due to the designation of SafeEnum, as it cannot be directly deserialized.
    Role enum.Serde[Role] `json:"role"`
}
```
