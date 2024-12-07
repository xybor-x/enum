// Package enum provides a type-safe and simple enum implementation for Go,
// offering easy conversion between numeric and string representations.
//
// It supports constant enums and integrates seamlessly with Go's `iota` enum pattern.
//
// Example usage:
//
//	type Role int
//
//	const (
//	  RoleUser Role = iota
//	  RoleAdmin
//	)
//
//	var (
//	  _ = enum.Map(RoleUser, "user")
//	  _ = enum.Map(RoleAdmin, "admin")
//	)
//
// Features:
//   - No code generation
//   - Supports constant enums
//   - Easy value conversions
//   - Support serialization
package enum

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/xybor-x/enum/internal"
	"github.com/xybor-x/enum/internal/mtkey"
	"github.com/xybor-x/enum/internal/mtmap"
)

const Undefined = "<<undefined>>"

func getAvailableEnumValue[T internal.Enumable]() T {
	id := T(0)
	for {
		if _, ok := mtmap.Get(mtkey.Enum2String(id)); !ok {
			break
		}
		id++
	}

	return T(id)
}

// New creates a dynamic enum value and maps it to a string representation.
//
// This function provides a convenient way to define and map an enum value
// without creating a constant explicitly. It is a shorthand for manually
// defining a constant and calling enum.Map().
//
// Note:
//   - Enums created with this function are variables, not constants. If you
//     need a constant enum, declare it explicitly and use enum.Map() instead.
//   - This method is not thread-safe and should only be called during
//     initialization or other safe execution points to avoid race conditions.
func New[T internal.Enumable](s string) T {
	e := getAvailableEnumValue[T]()
	return Map(e, s)
}

// Map associates an enum constant with a string representation.
//
// This function is used to map an enum value to its corresponding string,
// allowing easier handling of enums in contexts like serialization, logging,
// or user-facing output.
//
// In this example, RoleUser is mapped to "user" and RoleAdmin to "admin".
// Once mapped, these associations can be used to retrieve the string value
// or convert a string back to the corresponding enum.
//
// Note that this method is not thread-safe. Ensure mappings are set up during
// initialization or other safe execution points to avoid race conditions.
func Map[T internal.Enumable](value T, s string) T {
	if !strings.HasSuffix(ToString(value), Undefined) {
		panic("do not map a mapped enum")
	}

	if _, ok := FromString[T](s); ok {
		panic("do not map a mapped string")
	}

	mtmap.Set(mtkey.Enum2String(value), s)
	mtmap.Set(mtkey.String2Enum[T](s), value)

	allVals := mtmap.MustGet(mtkey.AllEnums[T]())
	allVals = append(allVals, value)
	mtmap.Set(mtkey.AllEnums[T](), allVals)

	return value
}

// FromInt returns the corresponding enum for a given int representation, and
// whether it is valid.
func FromInt[T internal.Enumable](i int) (T, bool) {
	t := T(i)
	return t, IsValid(t)
}

// MustFromInt returns the corresponding enum for a given int representation.
//
// It panics if the enum value is invalid.
func MustFromInt[T internal.Enumable](i int) T {
	t, ok := FromInt[T](i)
	if !ok {
		panic(fmt.Sprintf("enum %s: invalid", reflect.TypeOf(T(0)).Name()))
	}

	return t
}

// FromString returns the corresponding enum for a given string representation,
// and whether it is valid.
func FromString[T internal.Enumable](s string) (T, bool) {
	return mtmap.Get(mtkey.String2Enum[T](s))
}

// MustFromString returns the corresponding enum for a given string
// representation.
//
// It panics if the string does not correspond to a valid enum value.
func MustFromString[T internal.Enumable](s string) T {
	enum, ok := FromString[T](s)
	if !ok {
		panic(fmt.Sprintf("enum %s: invalid", reflect.TypeOf(T(0)).Name()))
	}

	return enum
}

// ToString returns the string representation of the given enum value.
//
// If the enum value is invalid, the returned string is formatted as:
// "EnumType::<<undefined>>".
func ToString[T internal.Enumable](value T) string {
	str, ok := mtmap.Get(mtkey.Enum2String(value))
	if !ok {
		return fmt.Sprintf("%s::%s", reflect.TypeOf(T(0)).Name(), Undefined)
	}
	return str
}

// MustToString returns the string representation of the given enum value.
//
// It panics if the enum value is invalid.
func MustToString[T internal.Enumable](value T) string {
	str := ToString(value)
	if strings.HasSuffix(str, Undefined) {
		panic(fmt.Sprintf("enum %s: invalid value %v", reflect.TypeOf(T(0)).Name(), value))
	}

	return str
}

// IsValid checks if an enum value is valid.
// It returns true if the enum value is valid, and false otherwise.
func IsValid[T internal.Enumable](value T) bool {
	_, ok := mtmap.Get(mtkey.Enum2String(value))
	return ok
}

// MarshalJSON serializes an enum value into its string representation.
func MarshalJSON[T internal.Enumable](value T) ([]byte, error) {
	if !IsValid(value) {
		return nil, fmt.Errorf("unknown %s: %v", reflect.TypeOf(T(0)).Name(), value)
	}

	return json.Marshal(ToString(value))
}

// UnmarshalJSON deserializes a string representation of an enum value from
// JSON.
func UnmarshalJSON[T internal.Enumable](data []byte, t *T) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	enum, ok := FromString[T](str)
	if !ok {
		return fmt.Errorf("unknown %s string: %s", reflect.TypeOf(T(0)).Name(), str)
	}

	*t = enum
	return nil
}

// ValueSQL serializes an enum into a database-compatible format.
func ValueSQL[T internal.Enumable](value T) (driver.Value, error) {
	if !IsValid(value) {
		return nil, fmt.Errorf("unknown %s: %v", reflect.TypeOf(T(0)).Name(), value)
	}

	return ToString(value), nil
}

// ScanSQL deserializes a database value into an enum type.
func ScanSQL[T internal.Enumable](a any, value *T) error {
	var data string
	switch t := a.(type) {
	case string:
		data = t
	case []byte:
		data = string(t)
	default:
		return fmt.Errorf("not support type %T", a)
	}

	enum, ok := FromString[T](data)
	if !ok {
		return fmt.Errorf("unknown %s string: %s", reflect.TypeOf(T(0)).Name(), data)
	}

	*value = enum

	return nil
}

// All returns a slice containing all enum values of a specific type.
func All[T internal.Enumable]() []T {
	return mtmap.MustGet(mtkey.AllEnums[T]())
}
