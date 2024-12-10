[![Go Reference](https://pkg.go.dev/badge/github.com/xybor-x/enum.svg)](https://pkg.go.dev/github.com/xybor-x/enum)
[![GitHub top language](https://img.shields.io/github/languages/top/xybor-x/enum?color=lightblue)](https://go.dev/)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/xybor-x/enum)](https://go.dev/blog/go1.21)
[![GitHub release (release name instead of tag name)](https://img.shields.io/github/v/release/xybor-x/enum?include_prereleases)](https://github.com/xybor-x/enum/releases/latest)
[![GitHub Repo stars](https://img.shields.io/github/stars/xybor-x/enum?color=blue)](https://github.com/xybor-x/enum)


![Golang](./.github/go-enum.png)

# ⚙️ Go Enum

**Elegant and powerful enums for Go with zero code generation!**

[1]: #-basic-enum
[2]: #-wrapenum
[3]: #-safeenum
[4]: #-utility-functions
[5]: #-constant-support
[6]: #-serialization-and-deserialization
[7]: #-type-safety

> [!WARNING]
> Please keep in mind that `xybor-x/enum` is still under active development
> and therefore full backward compatibility is not guaranteed before reaching v1.0.0.

## 🔧 Installation

```sh
go get -u github.com/xybor-x/enum
```

## 📋 Features

All of the following enum types are compatible with the APIs provided by `xybor-x/enum`.

|                                                | Basic enum ([#][1]) | Wrap enum ([#][2]) | Safe enum ([#][3]) |
| ---------------------------------------------- | ------------------- | ------------------ | ------------------ |
| **Built-in methods**                           | No                  | **Yes**            | **Yes**            |
| **Constant enum** ([#][5])                     | **Yes**             | **Yes**            | No                 |
| **Serialization and deserialization** ([#][6]) | No                  | **Yes**            | **Yes**            |
| **Type safety** ([#][7])                       | No                  | Basic              | **Strong**         |

> [!CAUTION]
> Enum definitions are not thread-safe.
> Therefore, they should be finalized during initialization (at the global scope).


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

`WrapEnum` offers a set of built-in methods to simplify working with enums.

**Pros 💪**
- Supports constant values ([#][5]).
- Provides many useful built-in methods.
- Full serialization and deserialization support out of the box.

**Cons 👎**
- Provides only **basic type safety** ([#][7]).

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
- Provides **strong type safety** ([#][7]).
- Provides many useful built-in methods.
- Full serialization and deserialization support out of the box.

**Cons 👎**
- Does not support constant values ([#][5]).

```go
// Define enum's underlying type.
type underlyingRole any

// Create a SafeEnum type for roles.
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

> [!NOTE]
> All of the following functions can be used with any type of enum.

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

> [!NOTE] 
> Not all enum types support serde operations, please refer to the [features](#-features).

## 🔅 Type safety

The [WrapEnum][2] prevents most invalid enum cases due to built-in methods for serialization and deserialization, offering **basic type safety**.

However, it is still possible to accidentally create an invalid enum value, like this:

```go
moderator := Role(42) // Invalid enum value
```

The [SafeEnum][3] provides **strong type safety**, ensuring that only predefined enum values are allowed. There is no way to create a new `SafeEnum` object without explicitly using the `NewSafe` function or zero initialization.

```go
moderator := Role(42)          // Compile-time error
moderator := Role("moderator") // Compile-time error
```

## 🔅 Extensible

> [!TIP]
> There are two ways to extend an enum type:
> - **For type definition enum**: Extend it directly.
> - **For type alias enum**: Embed it as a field in a struct, then add methods to the struct.
> 
> [What's in an (Alias) Name?](https://go.dev/blog/alias-names)

### Extend basic enum

Since this is a type definition enum, you can easily extend it by directly adding additional methods.

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

### Extend WrapEnum

> [!TIP]
> You should consider extending [Basic enum](#extend-basic-enum) or [Safe enum](#extend-safeenum) instead.

Since this is a type alias enum, the only way to maintain its built-in methods while still extending the enum is to wrap it as an embedded field in a struct.

However, this approach will break the constant-support property of the `WrapEnum`.

### Extend SafeEnum

Since this is a type alias enum, embedding it inside a struct allows you to maintain its predefined built-in methods while still being able to extend it.

```go
type underlyingRole any

type Role struct {
    enum.SafeEnum[underlyingRole]
}

var (
    RoleUser  = Role{enum.NewSafe[underlyingRole]("user")}
    RoleMod   = Role{enum.NewSafe[underlyingRole]("mod")}
    RoleAdmin = Role{enum.NewSafe[underlyingRole]("admin")}
    enum.Finalize[enum.SafeEnum[underlyingRole]]()
)

func (r Role) HasPermission() bool {
    return r == RoleMod || r == RoleAdmin
}

```
