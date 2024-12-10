[![Go Reference](https://pkg.go.dev/badge/github.com/xybor-x/enum.svg)](https://pkg.go.dev/github.com/xybor-x/enum)
[![GitHub top language](https://img.shields.io/github/languages/top/xybor-x/enum?color=lightblue)](https://go.dev/)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/xybor-x/enum)](https://go.dev/blog/go1.21)
[![GitHub release (release name instead of tag name)](https://img.shields.io/github/v/release/xybor-x/enum?include_prereleases)](https://github.com/xybor-x/enum/releases/latest)
[![GitHub Repo stars](https://img.shields.io/github/stars/xybor-x/enum?color=blue)](https://github.com/xybor-x/enum)


![Golang](./.github/go-enum.png)

# ⚙️ Go Enum

**Elegant and powerful enums for Go with zero code generation!**

[1]: #-iota-enum
[2]: #-wrapenum
[3]: #-safeenum
[4]: #-utility-functions
[5]: #-constant-support
[6]: #-serialization-and-deserialization
[7]: #-type-safety

## 🔧 Installation

```sh
go get -u github.com/xybor-x/enum
```

## 📋 Features

All enum types behave nearly consistently, so you can choose the style that best fits your use case without worrying about differences in functionality. You can refer to the [recommendations](#-recommendations).

|                            | Basic enum ([#][1]) | Wrap enum ([#][2]) | Safe enum ([#][3]) |
| -------------------------- | ------------------- | ------------------ | ------------------ |
| **Built-in methods**       | No                  | Yes                | Yes                |
| **Constant enum** ([#][5]) | Yes                 | Yes                | No                 |
| **Enum type**              | Any integer types   | `int`              | `struct`           |
| **Enum value type**        | Any integer types   | `int`              | `struct`           |
| **Serde** ([#][6])         | No                  | Yes                | Yes                |
| **Type safety** ([#][7])   | No                  | Basic              | Strong             |

❗ **Note**: Enum definitions are ***NOT thread-safe***. Therefore, they should be finalized during initialization (at the global scope).


## 🔍 Recommendations

|                               | Basic enum | Wrap enum | Safe enum |
| ----------------------------- | ---------- | --------- | --------- |
| **Simplified use**            | Yes        | Yes       | Yes       |
| **Exhaustive check required** | Yes        | Yes       | No        |
| **Type safety required**      | No         | Maybe     | Yes       |


## ⭐ Basic enum

The basic enum (`iota` approach) is the most commonly used enum implementation in Go.

It is essentially a primitive type, which does not include any built-in methods. For handling this type of enum, please refer to the [utility functions][4].

**Pros 💪**
- Simple.
- Supports constant values.

**Cons 👎**
- No built-in methods.
- Lacks serialization and deserialization support.
- No type safety.

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

## ⭐ WrapEnum
`WrapEnum` offers a set of built-in methods to simplify working with enums.

**Pros 💪**
- Supports constant values.
- Provides many useful built-in methods.
- Full serialization and deserialization support out of the box.

**Cons 👎**
- Provides only **basic type safety**.

```go
// Define enum's underlying type.
type underlyingRole any

// Create a WrapEnum type for roles.
type Role = enum.WrapEnum[underlyingRole] // NOTE: It must use type alias instead of type definition.

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
    ID   int  `json:"id"`
    Role Role `json:"role"`
}

func main() {
    // WrapEnum has many built-in methods for handling enum easier.
    data, _ := json.Marshal(RoleUser) // Output: "user"
    fmt.Println(RoleAdmin.IsValid())  // Output: true
}
```

## ⭐ SafeEnum

`SafeEnum` defines a strong type-safe enum. Like `WrapEnum`, it provides a set of built-in methods to simplify working with enums.

The `SafeEnum` enforces strict type safety, ensuring that only predefined enum values are allowed. It prevents the accidental creation of new enum types, providing a guaranteed set of valid values.

**Pros 💪**
- Provides **strong type safety**.
- Provides many useful built-in methods.
- Full serialization and deserialization support out of the box.

**Cons 👎**
- Does not support constant values.

```go
// Define enum's underlying type.
type underlyingRole string

// Create a StructEnum type for roles.
type Role = enum.SafeEnum[underlyingRole] // NOTE: It must use type alias instead of type definition.

var (
    RoleUser  = enum.NewSafe[underlyingRole]("user")
    RoleAdmin = enum.NewSafe[underlyingRole]("admin")
)

func main() {
    // SafeEnum has many built-in methods for handling enum easier.
    data, _ := json.Marshal(RoleUser) // Output: "user"
    fmt.Println(RoleAdmin.IsValid())  // Output: true
}
```

## 💡 Utility functions

*All of the following functions can be used with any style of enum. Note that this differs from the built-in methods, which are tied to the enum object rather than being standalone functions.*

### FromString

`FromString` returns the corresponding `enum` for a given `string` representation, and whether it is valid.

```go
role, ok := enum.FromString[Role]("user")
if ok {
    fmt.Println(role) // Output: 0
} else {
    fmt.Println("Invalid enum")
}
```

### FromInt

`FromInt` returns the corresponding `enum` for a given `int` representation, and whether it is valid.

```go
role, ok := enum.FromInt[Role](42)
if ok {
    fmt.Println(role)
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

`ToString` converts an `enum` to `string`. It returns `<nil>` for invalid enums.

```go
fmt.Println(enum.ToString(RoleAdmin))  // Output: "admin"
fmt.Println(enum.ToString(Role(42)))   // Output: "<nil>"
```

### ToInt

`ToInt` converts an `enum` to `int`.  It returns the smallest number of `int` for invalid enums.

```go
fmt.Println(enum.ToInt(RoleAdmin))  // Output: 1
fmt.Println(enum.ToInt(Role(42)))   // Output: -2147483648
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

Currently supported:
- `JSON`: Implements `json.Marshaler` and `json.Unmarshaler`.
- `SQL`: Implements `driver.Valuer` and `sql.Scanner`.

❗ *Note that NOT ALL enum styles support serde operations, please refer to the [features/serde](#-features).*

## 🔅 Type safety

`WrapEnum` includes built-in methods for serialization and deserialization, offering **basic type safety** and preventing most invalid enum cases.

However, it is still possible to accidentally create an invalid enum value, like this:

```go
moderator := Role(42) // Invalid enum value
```

The [`SafeEnum`][4] provides **strong type safety**, ensuring that only predefined enum values are allowed. There is no way to create a new `SafeEnum` object without explicitly using the `NewSafe` function or zero initialization.

```go
moderator := Role(42)          // Compile-time error
moderator := Role("moderator") // Compile-time error
```
