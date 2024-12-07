[![Go Reference](https://pkg.go.dev/badge/github.com/xybor-x/enum.svg)](https://pkg.go.dev/github.com/xybor-x/enum)

# Enum

**Elegant and powerful enums for Go with zero code generation!**

`xybor-x/enum` makes working with enums in Go intuitive and efficient by providing:
- Seamless enum-string mappings.
- Constant enums compatible with Go's `iota` conventions.
- Utility functions for JSON serialization and deserialization.
- Rich utilities for enhanced enum handling.


## Installation

Install the package via `go get`:
```sh
go get -u github.com/xybor-x/enum
```

## Quick start

*Please refer [Usage](#usage) for further details.*

In Go, `iota` is a special identifier used to create incrementing constants, making it perfect for defining enums.

`xybor-x/enum` is fully compatible with iota-based enums.

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

**Note**: Enum definitions are not thread-safe. Therefore, they should be finalized during initialization (at the global scope).

#### Short definition (variable enums)

The `enum.New` function allows dynamic initialization of enum values and maps them to a string representation.

```go
type Role int

var (
    RoleUser  = enum.New[Role]("user")
    RoleAdmin = enum.New[Role]("admin")
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
}
```

### Utility functions

#### EnumOf

```go
role, ok := enum.EnumOf[Role]("user")
if ok {
    fmt.Println("Enum representation:", role) // Output: 0
} else {
    fmt.Println("Invalid enum")
}
```

#### StringOf

```go
fmt.Println("String:", enum.StringOf(RoleAdmin)) // Output: "admin"
```

#### IsValid

```go
fmt.Println(enum.IsValid(RoleUser)) // true
fmt.Println(enum.IsValid(Role(0)))  // true
fmt.Println(enum.IsValid(Role(3)))  // false
```

#### All

```go
for _, role := range enum.All[Role]() {
    fmt.Println("Role:", enum.StringOf(role))
}
// Output:
// Role: user
// Role: admin
```

### Rich enum

Instead of defining your enum type as `int`, you can use `enum.RichEnum` (an `int` alias) to to leverage several convenient features:
- Interact with enums using methods instead of standalone functions.
- Built-in support for serialization and deserialization.

```go
type role any
type Role = enum.RichEnum[role] // NOTE: It must use type alias instead of type definition.

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
