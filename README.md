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

*Please refer [Usage](./README.md#usage) for further details.*

In Go, `iota` is a special identifier used to create incrementing constants, making it perfect for defining enums.

```go
type role any
type Role = enum.RichEnum[role] // simply an alias of int rather than a struct

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
    enum.Map(RoleUser, "user")   // Maps RoleUser to "user"
    enum.Map(RoleAdmin, "admin") // Maps RoleAdmin to "admin"
}
```

### Utility functions

#### EnumOf

```go
role := enum.EnumOf[Role]("user")
if enum.IsValid(role) {
    fmt.Println("Enum representation:", role) // Output: RoleUser
} else {
    fmt.Println("Invalid enum")
}
```

#### StringOf

```go
fmt.Println("String:", enum.StringOf(RoleAdmin)) // Output: "admin"
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

Extend functionality with `enum.RichEnum`, which implements `fmt.Stringer`, `json.Marshaler`, and `json.Unmarshaler`.

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
