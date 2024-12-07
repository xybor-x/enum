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
)

const UndefinedString = "<<undefined>>"

var enums = &mtmap{}

type enumable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func getAvailableEnumValue[T enumable]() T {
	id := 0
	for {
		if _, ok := get2(enums, key[T, string]{T(id)}); !ok {
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
// Example:
//
//	type Role int
//
//	var (
//		RoleUser = enum.New("user")   // Dynamically creates and maps "user"
//		RoleAdmin = enum.New("admin") // Dynamically creates and maps "admin"
//	)
//
// Note:
//   - Enums created with this function are variables, not constants. If you
//     need a constant enum, declare it explicitly and use enum.Map() instead.
//   - This method is not thread-safe and should only be called during
//     initialization or other safe execution points to avoid race conditions.
func New[T enumable](s string) T {
	e := getAvailableEnumValue[T]()
	return Map(e, s)
}

// Map associates an enum constant with a string representation.
//
// This function is used to map an enum value to its corresponding string,
// allowing easier handling of enums in contexts like serialization, logging,
// or user-facing output.
//
// Example:
//
//	type Role int
//	const (
//		RoleUser Role = iota
//		RoleAdmin
//	)
//
//	var (
//		_ = enum.Map(RoleUser, "user")
//		_ = enum.Map(RoleAdmin, "admin")
//	)
//
// In this example, RoleUser is mapped to "user" and RoleAdmin to "admin".
// Once mapped, these associations can be used to retrieve the string value
// or convert a string back to the corresponding enum.
//
// Note that this method is not thread-safe. Ensure mappings are set up during
// initialization or other safe execution points to avoid race conditions.
func Map[T enumable](value T, s string) T {
	if !strings.HasSuffix(StringOf(value), UndefinedString) {
		panic("do not map a mapped enum")
	}

	if _, ok := EnumOf[T](s); ok {
		panic("do not map a mapped string")
	}

	set(enums, key[T, string]{T(value)}, s)
	set(enums, key[string, T]{s}, value)
	allVals := get(enums, key[T, []T]{T(0)})
	allVals = append(allVals, value)
	set(enums, key[T, []T]{T(0)}, allVals)

	return value
}

// EnumOf returns the corresponding enum for a given string
// representation, and whether it is valid.
func EnumOf[T enumable](s string) (T, bool) {
	enum, ok := get2(enums, key[string, T]{s})
	if !ok {
		return T(0), false
	}
	return enum, true
}

// MustEnumOf returns the corresponding enum for a given string representation.
// It panics if the string does not correspond to a valid enum value.
func MustEnumOf[T enumable](s string) T {
	enum, ok := get2(enums, key[string, T]{s})
	if !ok {
		panic(fmt.Sprintf("enum %s: invalid", reflect.TypeOf(T(0)).Name()))
	}

	return enum
}

// StringOf returns the string representation of an enum value.
func StringOf[T enumable](value T) string {
	enum, ok := get2(enums, key[T, string]{value})
	if !ok {
		return fmt.Sprintf("%s::%s", reflect.TypeOf(T(0)).Name(), UndefinedString)
	}
	return enum
}

// MustStringOf returns the string representation of an enum value.
// It panics if the enum value is invalid.
func MustStringOf[T enumable](value T) string {
	str := StringOf(value)
	if strings.HasSuffix(str, UndefinedString) {
		panic(fmt.Sprintf("enum %s: invalid value %v", reflect.TypeOf(T(0)).Name(), value))
	}

	return str
}

// IsValid checks if an enum value is valid.
// It returns true if the enum value is valid, and false otherwise.
func IsValid[T enumable](value T) bool {
	_, ok := get2(enums, key[T, string]{value})
	return ok
}

// MarshalJSON serializes an enum value into its string representation.
func MarshalJSON[T enumable](value T) ([]byte, error) {
	if !IsValid(value) {
		return nil, fmt.Errorf("unknown %s: %v", reflect.TypeOf(T(0)).Name(), value)
	}

	return json.Marshal(StringOf(value))
}

// UnmarshalJSON deserializes a string representation of an enum value from
// JSON.
func UnmarshalJSON[T enumable](data []byte, t *T) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	enum, valid := EnumOf[T](str)
	if !valid {
		return fmt.Errorf("unknown %s string: %s", reflect.TypeOf(T(0)).Name(), str)
	}

	*t = enum
	return nil
}

// ValueSQL serializes an enum into a database-compatible format.
func ValueSQL[T enumable](value T) (driver.Value, error) {
	if !IsValid(value) {
		return nil, fmt.Errorf("unknown %s: %v", reflect.TypeOf(T(0)).Name(), value)
	}

	return StringOf(value), nil
}

// ScanSQL deserializes a database value into an enum type.
func ScanSQL[T enumable](a any, value *T) error {
	var data string
	switch t := a.(type) {
	case string:
		data = t
	case []byte:
		data = string(t)
	default:
		return fmt.Errorf("not support type %T", a)
	}

	enum, ok := EnumOf[T](data)
	if !ok {
		return fmt.Errorf("unknown %s string: %s", reflect.TypeOf(T(0)).Name(), data)
	}

	*value = enum

	return nil
}

// All returns a slice containing all enum values of a specific type.
//
// This function iterates over all defined constants of the given enum type and
// returns them as a slice, enabling easy access to the complete set of enum values.
//
// Example usage:
//
//	roles := All[Role]() // roles = []Role{RoleAdmin, RoleUser}
func All[T enumable]() []T {
	return get(enums, key[T, []T]{T(0)})
}
