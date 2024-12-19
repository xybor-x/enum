[![Go Reference](https://pkg.go.dev/badge/github.com/xybor-x/enum.svg)](https://pkg.go.dev/github.com/xybor-x/enum)
[![GitHub top language](https://img.shields.io/github/languages/top/xybor-x/enum?color=lightblue)](https://go.dev/)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/xybor-x/enum)](https://go.dev/blog/go1.21)
[![GitHub release (release name instead of tag name)](https://img.shields.io/github/v/release/xybor-x/enum?include_prereleases)](https://github.com/xybor-x/enum/releases/latest)
[![GitHub Repo stars](https://img.shields.io/github/stars/xybor-x/enum?color=blue)](https://github.com/xybor-x/enum)

![Golang](./.github/go-enum.png)

# ⚙️ Go Enum

**Elegant, powerful, and dependency-free enums for Go with zero code generation!**

## 🔧 Installation

```sh
go get -u github.com/xybor-x/enum
```

## ⚡ Quick start

> [!TIP]
> This is just a ⚡ quick definition for general use cases.
> See more advanced features at [documentations](./docs.md).

> [!CAUTION]
> Enum definitions are not thread-safe.
> Therefore, they should be finalized during initialization (at the global scope).

**Definition**

```go
package main

import "github.com/xybor-x/enum"

type role any
type Role = enum.WrapEnum[role]

// Optional: NullRole is similar to sql.NullXXX structs, designed to handle 
// nullable SQL and JSON fields.
type NullRole = enum.Nullable[Role]

const (
    RoleUser Role = iota
    RoleAdmin
)

var (
    _ = enum.Map(RoleUser, "user")
    _ = enum.Map(RoleAdmin, "admin")
    _ = enum.Finalize[Role]()
)
```

**Usage**

```go
package main

import (
    "encoding/json"
    "fmt"
)

type User struct {
    Role Role `json:"role"`
}

type NullUser struct {
    Role NullRole `json:"role"`
}

func main() {
    // Print out the string representation of enum.
    fmt.Println(RoleAdmin) // Output: admin

    // Serialize a user to json.
    user := User{Role: RoleUser}
    data, _ := json.Marshal(user) 
    fmt.Println(string(data)) // Output: {"role": "user"}

    // Serialize a nullable role with a non-null value.
    nulluser := NullUser{Role: NullRole{Enum: RoleUser, Valid: true}}
    data, _ := json.Marshal(nulluser) 
    fmt.Println(string(data)) // Output: {"role": "user"}

    // Serialize a nullable role with a null value.
    data, _ = json.Marshal(NullUser{})
    fmt.Println(string(data)) // Output: {"role": null}
}
```
