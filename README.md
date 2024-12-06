[![Go Reference](https://pkg.go.dev/badge/github.com/xybor-x/enum.svg)](https://pkg.go.dev/github.com/xybor-x/enum)

# Enum

An easy-to-use Go library for working with enums.

It simplifies serialization, deserialization, and supports constant enums, all without the need for code generation.

## Features

- **No Code Generation**: Simplifies usage by eliminating the need for additional tools or build steps.
- **Supports Constant Enums**: Enables defining immutable enum values for safer and more predictable behavior.
- **Compatible with Standard Enum Definitions**: Easily integrates with Go's conventional enum patterns, including `iota`.
- **Easy Conversion**: Effortlessly convert between numeric and string representations for better usability and flexibility.
- **Serializable**: Provides JSON serialization and deserialization with `enum.RichEnum`, making it easier to work with enums in JSON.

## Usage

### Define enum using `enum.New`

The `enum.New` function allows dynamic initialization of enum values and maps them to a string representation.

```go
type Role int

var (
    RoleUser  = enum.New[Role]("user")  // Dynamically creates and maps "user"
    RoleAdmin = enum.New[Role]("admin") // Dynamically creates and maps "admin"
)
```

### Define enum using `enum.Map`

The `enum.Map` function supports defining constant enums and is fully compatible with Go's standard enum patterns.

``` go
type Role int

const (
    RoleUser Role = iota
    RoleAdmin
)

var (
    _ = enum.Map(RoleUser, "user")   // Maps RoleUser to "user"
    _ = enum.Map(RoleAdmin, "admin") // Maps RoleAdmin to "admin"
)
```

With `enum.Map`, you can associate string values with constants, making them easier to work with in serialization or user-facing output.

### Rich Enum

`enum.RichEnum` is a struct with set of utility methods to simplify working with enums.

It includes various helper functions for operations like serialization, deserialization, string conversion, and validation, making it easier to manage and manipulate enum values across your codebase.

Please refer [example_test.go](./example_test.go) for detailed examples.

```go
type unsafeRole any
type Role = enum.RichEnum[unsafeRole] // NOTE: It must use type alias instead of type definition.

const (
    RoleUser Role = iota
    RoleAdmin
)

var (
    _ = enum.Map(RoleUser, "user")   // Maps RoleUser to "user"
    _ = enum.Map(RoleAdmin, "admin") // Maps RoleAdmin to "admin"
)

data, err := json.Marshal(RoleUser) // data should be "user"
fmt.Println(RoleAdmin)              // Output: "admin" (instead of 1)
```
