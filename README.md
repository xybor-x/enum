[![Go Reference](https://pkg.go.dev/badge/github.com/xybor-x/enum.svg)](https://pkg.go.dev/github.com/xybor-x/enum)
[![GitHub top language](https://img.shields.io/github/languages/top/xybor-x/enum?color=lightblue)](https://go.dev/)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/xybor-x/enum)](https://go.dev/blog/go1.21)
[![GitHub release (release name instead of tag name)](https://img.shields.io/github/v/release/xybor-x/enum?include_prereleases)](https://github.com/xybor-x/enum/releases/latest)
[![GitHub Repo stars](https://img.shields.io/github/stars/xybor-x/enum?color=blue)](https://github.com/xybor-x/enum)


![Golang](./.github/go-enum.png)

# ⚙️ Go Enum

**Elegant and powerful enums for Go with zero code generation!**

[1]: #-basic-enum
[2]: #-int-enum
[3]: #-struct-enum
[4]: #-safe-enum
[5]: #-utility-functions
[6]: #-constant-support
[7]: #-serialization-and-deserialization
[8]: #-type-safety

## 🔧 Installation

Install the package via `go get`:
```sh
go get -u github.com/xybor-x/enum
```

## 📋 Features

All enum types behave nearly consistently, so you can choose the style that best fits your use case without worrying about differences in functionality.

|                            | Basic enum ([#][1]) | Int enum ([#][2]) | Struct Enum [#][3] | Safe enum ([#][4]) |
| -------------------------- | ------------------- | ----------------- | ------------------ | ------------------ |
| **Built-in methods**       | No                  | Yes               | Yes                | Yes                |
| **Constant enum** ([#][6]) | Yes                 | Yes               | No                 | No                 |
| **Enum type**              | Any integer types   | `int`             | `struct`           | `interface`        |
| **Enum value type**        | Any integer types   | `int`             | `struct`           | `struct`           |
| **Serde** ([#][7])         | No                  | Full              | Full               | Serialization only |
| **Type safety** ([#][8])   | Basic               | Basic             | Good               | Strong             |

❗ **Note**: Enum definitions are ***NOT thread-safe***. Therefore, they should be finalized during initialization (at the global scope).

## ⭐ Basic enum

Basic enum is the simplest type, but since it has no built-in methods, please refer to the [utility functions][5] for handling this enum.

**Pros 💪**
- Simplest.
- Supports constant values.

**Cons 👎**
- No built-in methods.
- Lacks serialization and deserialization support.
- Provides only **basic type safety**.

### Dynamic style

```go
type Role int

// Dynamic style doesn't support constant enum value.
var (
    RoleUser  = enum.New[Role]("user")
    RoleAdmin = enum.New[Role]("admin")
    _         = enum.Finalize[Role]() // Optional: ensure no new enum values can be added to Role.
)
```

### Static style

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

## ⭐ Int enum
`IntEnum` offers a set of built-in methods to simplify working with enums.

**Pros 💪**
- Supports constant values.
- Provides many useful built-in methods.
- Full serialization and deserialization support out of the box.

**Cons 👎**
- Provides only **basic type safety**.

```go
// Define enum's underlying type.
type underlyingRole any

// Create a IntEnum type for roles.
type Role = enum.IntEnum[underlyingRole] // NOTE: It must use type alias instead of type definition.

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
    // IntEnum has many built-in methods for handling enum easier.
    data, _ := json.Marshal(RoleUser) // Output: "user"
    fmt.Println(RoleAdmin.IsValid())  // Output: true
}
```

## ⭐ Struct enum

`StructEnum` provides a good type-safe enum, which is better than `IntEnum`, but not as good as `SafeEnum`. Like `IntEnum`, it provides a set of built-in methods to simplify working with enums.

**Pros 💪**
- Provides many useful built-in methods.
- Full serialization and deserialization support out of the box.
- Provides **good type safety**.

**Cons 👎**
- Does not support constant values.

```go
// Define enum's underlying type.
type underlyingRole any

// Create a StructEnum type for roles.
type Role = enum.StructEnum[underlyingRole] // NOTE: It must use type alias instead of type definition.

var (
    RoleUser  = enum.NewStruct[underlyingRole]("user")
    RoleAdmin = enum.NewStruct[underlyingRole]("admin")
)

func main() {
    // StructEnum has many built-in methods for handling enum easier.
    data, _ := json.Marshal(RoleUser) // Output: "user"
    fmt.Println(RoleAdmin.IsValid())  // Output: true
}
```

## ⭐ Safe enum

`SafeEnum` defines a strong type-safe enum. Like `IntEnum`, it provides a set of built-in methods to simplify working with enums.

The `SafeEnum` enforces strict type safety, ensuring that only predefined enum values are allowed. It prevents the accidental creation of new enum types, providing a guaranteed set of valid values.

**Pros 💪**
- Provides **strong type safety**.
- Provides many useful built-in methods.
- Serialization support out of the box.

**Cons 👎**
- Does not support constant values.
- Lacks deserialization support.

```go
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

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    
    // Use enum.Serde due to the designation of SafeEnum, as it cannot be directly deserialized.
    Role enum.Serde[Role] `json:"role"`
}
```

## 💡 Utility functions

*All of the following functions can be used with any style of enum. Note that this differs from the built-in methods, which are tied to the enum object rather than being standalone functions.*

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
// unsafe enums.
fmt.Println(enum.ToString(Role(42)))   // panic
```

### ToInt

`ToInt` converts an `enum` to `int`. It panics if the `enum` is invalid.

```go
fmt.Println(enum.ToInt(RoleAdmin))  // Output: 1

// Note that you should check if the enum is valid before calling ToInt for
// unsafe enums.
fmt.Println(enum.ToInt(Role(42)))   // panic
```

### All

`All` returns a slice containing all enum values of a specific enum type.

```go
for _, role := range enum.All[Role]() {
    fmt.Println("Role:", enum.ToString(role))
}
// Output:
// Role: user
// Role: admin
```

## 🔅 Constant support

Some static analysis tools support checking for exhaustive `switch` statements in constant enums. By choosing an `enum` with constant support, you can enable this functionality in these tools.

## 🔅 Serialization and deserialization

Serialization and deserialization are essential when working with enums, and our library provides seamless support for handling them out of the box.

❗ *Note that NOT ALL enum styles support both serialization and deserialization, please refer to the [features/serde](#-features).*

Currently supported:
- `JSON`: Implements `json.Marshaler` and `json.Unmarshaler`.
- `SQL`: Implements `driver.Valuer` and `sql.Scanner`.

For enum styles that do not natively support serialization or deserialization, the `Serde` wrapper can be used to enable this functionality.

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

## 🔅 Type safety

By default, `xybor-x/enum` provides [functions][5] to parse and validate an `enum`, offering **basic type safety**.

However, it is still possible to accidentally create an invalid enum value, like so:

```go
moderator := Role(42) // Invalid enum value
```

The [`StructEnum`][3] offers **good type safety**, which prevents most accidental creation of invalid enum values, but still allows for invalid values due to zero initialization:

```go
role := Role{} // Normally, enums are not created in this way, except for unmarshalling purposes.
```

The [`SafeEnum`][4] provides **strong type safety**, ensuring that only predefined enum values are allowed. There is no way to create a new `SafeEnum` without explicitly using the `safeenum.New` function.
