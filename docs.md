# ⚙️ Go Enum

**Elegant, powerful, and dependency-free enums for Go with zero code generation!**

[1]: #-basic-enum
[2]: #-wrapenum
[3]: #-safeenum
[4]: #-utility-functions
[5]: #-constant-support
[6]: #-serialization-and-deserialization
[7]: #-type-safety
[8]: #-nullable

## Table of contents

- [⚙️ Go Enum](#️-go-enum)
  - [Table of contents](#table-of-contents)
  - [🔧 Installation](#-installation)
  - [📋 Features](#-features)
  - [⭐ Basic enum](#-basic-enum)
  - [⭐ WrapEnum](#-wrapenum)
  - [⭐ SafeEnum](#-safeenum)
  - [💡 Utility functions](#-utility-functions)
  - [🔅 Constant support](#-constant-support)
  - [🔅 Serialization and deserialization](#-serialization-and-deserialization)
  - [🔅 Nullable](#-nullable)
  - [🔅 Type safety](#-type-safety)
  - [🔅 Extensible](#-extensible)

## 🔧 Installation

```sh
go get -u github.com/xybor-x/enum
```

## 📋 Features

> [!TIP]
> `xybor-x/enum` supports three enum types: **Basic enum** for simplicity, **Wrap enum** for enhanced functionality, and **Safe enum** for strict type safety.

|                                                | Basic enum ([#][1]) | Wrap enum ([#][2]) | Safe enum ([#][3]) |
| :--------------------------------------------- | ------------------- | ------------------ | ------------------ |
| **Underlying type required**                   | **No**              | Yes                | Yes                |
| **Built-in methods**                           | No                  | **Yes**            | **Yes**            |
| **Constant enum** ([#][5])                     | **Yes**             | **Yes**            | No                 |
| **Serialization and deserialization** ([#][6]) | No                  | **Yes**            | **Yes**            |
| **Type safety** ([#][7])                       | No                  | Basic              | **Strong**         |
| **Used with Nullable** ([#][8])                | **Yes**             | **Yes**            | **Yes**            |

## ⭐ Basic enum

The basic enum (a.k.a `iota` enum) is the most commonly used enum implementation in Go.

It is essentially a primitive type, which does not include any built-in methods. For handling this type of enum, please refer to the [utility functions][4].

**Pros 💪**
- Simple.
- Supports constant values ([#][5]).

**Cons 👎**
- No built-in methods.
- No type safety ([#][7]).
- Lacks serialization and deserialization support.

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

`WrapEnum` offers a set of built-in methods to simplify working with `int` enums.

> [!TIP]
> For other numeric types, use `WrapUintEnum` for `uint` and `WrapFloatEnum` for `float64`.

**Pros 💪**
- Supports constant values ([#][5]).
- Provides many useful built-in methods.
- Full serialization and deserialization support out of the box.

**Cons 👎**
- Provides only **basic type safety** ([#][7]).

```go
// Define enum's underlying type.
type role any

// Create a WrapEnum type for roles.
type Role = enum.WrapEnum[role] // NOTE: It must use type alias instead of type definition.

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
    // WrapEnum has many built-in methods for handling enum easier.
    data, _ := json.Marshal(RoleUser) // Output: "user"
    fmt.Println(RoleAdmin.IsValid())  // Output: true
}
```

## ⭐ SafeEnum

`SafeEnum` defines a strong type-safe enum. Like `WrapEnum`, it provides a set of built-in methods to simplify working with enums.

The `SafeEnum` enforces strict type safety, ensuring that only predefined enum values are allowed. It prevents the accidental creation of new enum types, providing a guaranteed set of valid values.

**Pros 💪**
- Provides **strong type safety** ([#][7]).
- Provides many useful built-in methods.
- Full serialization and deserialization support out of the box.

**Cons 👎**
- Does not support constant values ([#][5]).

```go
// Define enum's underlying type.
type role any

// Create a SafeEnum type for roles.
type Role = enum.SafeEnum[role] // NOTE: It must use type alias instead of type definition.

var (
    RoleUser  = enum.New[Role]("user")
    RoleAdmin = enum.New[Role]("admin")
)

func main() {
    // SafeEnum has many built-in methods for handling enum easier.
    data, _ := json.Marshal(RoleUser) // Output: "user"
    fmt.Println(RoleAdmin.IsValid())  // Output: true
}
```

## 💡 Utility functions

> [!NOTE]
> All of the following functions can be used with any type of enum.

**FromString**

`FromString` returns the corresponding `enum` for a given `string` representation, and whether it is valid.

```go
role, ok := enum.FromString[Role]("user")
if ok {
    fmt.Println(role) // Output: 0
} else {
    fmt.Println("Invalid enum")
}
```

**FromNumber**

`FromNumber` returns the corresponding `enum` for a given numeric representation, and whether it is valid.

```go
role, ok := enum.FromNumber[Role](42)
if ok {
    fmt.Println(role)
} else {
    fmt.Println("Invalid enum") // Output: Invalid enum
}
```

**IsValid**

`IsValid` checks if an enum value is valid or not.

```go
fmt.Println(enum.IsValid(RoleUser))  // true
fmt.Println(enum.IsValid(Role(0)))   // true
fmt.Println(enum.IsValid(Role(42)))  // false
```

**ToString**

`ToString` converts an `enum` to `string`. It returns `<nil>` for invalid enums.

```go
fmt.Println(enum.ToString(RoleAdmin))  // Output: "admin"
fmt.Println(enum.ToString(Role(42)))   // Output: "<nil>"
```

**All**

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

> [!WARNING] 
> Not all enum types support serde operations out of the box, please refer to the [features](#-features).

Currently supported:
- `JSON`: Implements `json.Marshaler` and `json.Unmarshaler`.
- `SQL`: Implements `driver.Valuer` and `sql.Scanner`.

## 🔅 Nullable

The `Nullable` transforms an enum type into a nullable enum, akin to `sql.NullXXX`, and is designed to handle nullable values in both JSON and SQL.

```go
type Role int
type NullRole = enum.Nullable[Role]

type User struct {
    ID   int      `json:"id"`
    Role NullRole `json:"role"`
}

func main() {
    fmt.Println(json.Marshal(User{})) // {"id": 0, "role": null}
}
```

## 🔅 Type safety

The [WrapEnum][2] prevents most invalid enum cases due to built-in methods for serialization and deserialization, offering **basic type safety**.

However, it is still possible to accidentally create an invalid enum value, like this:

```go
moderator := Role(42) // Invalid enum value
```

The [SafeEnum][3] provides **strong type safety**, ensuring that only predefined enum values are allowed. There is no way to create a new `SafeEnum` object without explicitly using the `New` function or zero initialization.

```go
moderator := Role(42)          // Compile-time error
moderator := Role("moderator") // Compile-time error
```

## 🔅 Extensible

**Extend basic enum**

Since this enum is just a primitive type and does not have built-in methods, you can easily extend it by directly adding additional methods.

```go
type Role int

const (
    RoleUser Role = iota
    RoleMod
    RoleAdmin
)

func init() {
    enum.Map(RoleUser, "user")
    enum.Map(RoleMod, "mod")
    enum.Map(RoleAdmin, "admin")
    enum.Finalize[Role]()
}

func (r Role) HasPermission() bool {
    return r == RoleMod || r == RoleAdmin
}
```

**Extend WrapEnum**

`WrapEnum` has many predefined methods. The only way to retain these methods while extending it is to wrap it as an embedded field in another struct.

However, this approach will break the constant-support property of the `WrapEnum` because Go does not support constants for structs.

You should consider extending [Basic enum](#extend-basic-enum) or [Safe enum](#extend-safeenum) instead.

**Extend SafeEnum**

`SafeEnum` has many predefined methods. The only way to retain these methods while extending it is to wrap it as an embedded field in another struct.

`xybor-x/enum` provides the `NewExtended` function to help create a wrapper of advanced enums.

```go
type role any
type Role struct { enum.SafeEnum[role] }

var (
    RoleUser  = enum.NewExtended[Role]("user")
    RoleMod   = enum.NewExtended[Role]("mod")
    RoleAdmin = enum.NewExtended[Role]("admin")
    _         = enum.Finalize[Role]()
)

func (r Role) HasPermission() bool {
    return r == RoleMod || r == RoleAdmin
}
```
