[![Go Reference](https://pkg.go.dev/badge/github.com/xybor-x/enum.svg)](https://pkg.go.dev/github.com/xybor-x/enum)

# Enum
A simple enum utility provided by Todennus.

## Features

- **No Code Generation**: Simplifies usage by eliminating the need for additional tools or build steps.
- **Supports Constant Enums**: Enables defining immutable enum values for safer and more predictable behavior.
- **Compatible with Standard Enum Definitions**: Seamlessly integrates with Go's conventional enum patterns, including `iota`.
- **Easy Conversion**: Effortlessly convert between numeric and string representations for better usability and flexibility.

## Usage

### Using `enum.New`

The `enum.New` function allows dynamic initialization of enum values and maps them to a string representation.

```go
type Role int

var (
    RoleUser  = enum.New[Role]("user")  // Dynamically creates and maps "user"
    RoleAdmin = enum.New[Role]("admin") // Dynamically creates and maps "admin"
)
```

### Using `enum.Map`

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
