[![Go Reference](https://pkg.go.dev/badge/github.com/xybor-x/enum.svg)](https://pkg.go.dev/github.com/xybor-x/enum)

# Go Enum

**Elegant and powerful enums for Go with zero code generation!**

[1]: #basic-enum
[2]: #rich-enum
[3]: #safe-enum
[4]: #utility-functions
[5]: #constant-support
[6]: #serialization-and-deserialization
[7]: #type-safety

## Installation

Install the package via `go get`:
```sh
go get -u github.com/xybor-x/enum
```

## Features

|                           | [Basic enum][1]   | [Rich enum][2] | [Safe enum][3]     |
| ------------------------- | ----------------- | -------------- | ------------------ |
| [**Utility support**][4]  | Yes               | Yes            | Yes                |
| [**Constant-support**][5] | Yes               | Yes            | No                 |
| **Enum type**             | Any integer types | `int`          | `interface`        |
| **Enum value type**       | Any integer types | `int`          | `struct`           |
| [**Serde**][6]            | No                | Full           | Serialization only |
| [**Type safety**][7]      | Basic             | Basic          | Strong             |

**Note**: Enum definitions are ***NOT thread-safe***. Therefore, they should be finalized during initialization (at the global scope).

## Basic enum

### Dynamic style

```go
type Role int

// Dynamic style only supports variable enum value.
var (
    RoleUser  = enum.New[Role]("user")
    RoleAdmin = enum.New[Role]("admin")
    _         = enum.Finalize[Role]() // Optional: ensure no new enum values can be added to Role.
)
```

### Static style (constant support)

``` go
type Role int

// Static style supports constant enum value.
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


## Rich enum

```go
// Define enum's underlying type.
type role any

// Create a RichEnum type for roles.
type Role = enum.RichEnum[role] // NOTE: It must use type alias instead of type definition.

// Basic enum definition styles can also be used here. 
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
    // RichEnum has many utility methods for handling enum easier.
    data, _ := json.Marshal(RoleUser) // Output: "user"
    fmt.Println(RoleAdmin.IsValid())  // Output: true
}
```

## Safe enum

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

## Utility functions

*All following functions can be used with any style of enum.*

### FromString

`FromString` returns the corresponding `enum` for a given `string` representation, and whether it is valid.

```go
role, ok := enum.FromString[Role]("user")
if ok {
    fmt.Println("Enum representation:", role) // Output: 0
} else {
    fmt.Println("Invalid enum")
}
```

### FromInt

`FromInt` returns the corresponding `enum` for a given `int` representation, and whether it is valid.

```go
role, ok := enum.FromInt[Role](42)
if ok {
    fmt.Println("Enum representation:", role)
} else {
    fmt.Println("Invalid enum") // Output: Invalid enum
}
```

### IsValid

`IsValid` checks if an enum value is valid or not.

```go
fmt.Println(enum.IsValid(RoleUser))  // true
fmt.Println(enum.IsValid(Role(0)))   // true
fmt.Println(enum.IsValid(Role(42)))  // false
```

### ToString

`ToString` converts an `enum` to `string`. It panics if the `enum` is invalid.

```go
fmt.Println(enum.ToString(RoleAdmin))  // Output: "admin"

// Note that you should check if the enum is valid before calling ToString for
// an unsafe enum.
fmt.Println(enum.ToString(Role(42)))   // panic
```

### ToInt

`ToInt` converts an `enum` to `int`. It panics if the `enum` is invalid.

```go
fmt.Println(enum.ToInt(RoleAdmin))  // Output: 1

// Note that you should check if the enum is valid before calling ToInt for
// an unsafe enum.
fmt.Println(enum.ToInt(Role(42)))   // panic
```

### All

`All` returns a slice containing all enum values of a specific type.

```go
for _, role := range enum.All[Role]() {
    fmt.Println("Role:", enum.ToString(role))
}
// Output:
// Role: user
// Role: admin
```

## Constant support

Some static analysis tools support checking for exhaustive `switch` statements in constant enums:

- [golangci-lint](https://github.com/golangci/golangci-lint) – A meta-linter that integrates various linters, including exhaustive checking for enums.
- [exhaustive](https://github.com/nishanths/exhaustive) – A tool specifically designed to ensure switch statements handle all possible values of constant enums.
- [go-critic](https://github.com/go-critic/go-critic) – A comprehensive Go linter with checks, including validation for exhaustive switch statements.
- [nogo for Bazel](https://github.com/bazel-contrib/rules_go/blob/master/go/nogo.rst) – A Bazel-specific tool for enforcing static checks, including exhaustive switch checks for Go code.

## Serialization and deserialization

Serialization and deserialization are essential when working with enums, and our library provides seamless support for handling them out of the box.

Currently supported:
- `JSON`: Implements `json.Marshaler` and `json.Unmarshaler`.
- `SQL`: Implements `driver.Valuer` and `sql.Scanner`.

For enums that do not natively support serialization or deserialization, the `Serde` wrapper can be used to enable this functionality.

```go
// Basic enum doesn't support serialization and deserialization.
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

type User struct {
    ID   int              `json:"id"`
    Name string           `json:"name"`
    Role enum.Serde[Role] `json:"role"` // Without Serde, Role will be serialized as a normal int value.
}

func main() {
    user := User{ID: 1, Name: "serde", Role: enum.SerdeWrap(RoleAdmin)}
    data, _ := json.Marshal(user)
    fmt.Println(string(data)) // {"id": 1, "name": "serde", "role": "admin"}

    deuser := User{}
    json.Unmarshal(data, &deuser)
    fmt.Println(deuser.Role.Enum()) // 1
}
```

## Type safety

By default, `xybor-x/enum` provides functions to parse `enums` from [`string`](#fromstring) or [`int`](#fromint). These functions also help validate the enum values, offering a **basic type safety**.

However, it is still possible to accidentally create an invalid enum value, like so:

```go
moderator := Role(42) // Invalid enum value
```

The [`safe enum`][3] provides **strong type safety**, ensuring that only predefined enum values are allowed. There is no way to create a new `safe enum` without explicitly using the `safeenum.New` function.
