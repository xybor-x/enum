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
// - No code generation
// - Supports constant enums
// - Easy value conversions
package enum

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const UndefinedString = "<<undefined>>"

var enums = &mtmap{}

func getAvailableEnumValue[T ~int]() T {
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
func New[T ~int](s string) T {
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
func Map[T ~int](value T, s string) T {
	if !strings.HasSuffix(StringOf(value), UndefinedString) {
		panic("do not map a mapped enum")
	}

	set(enums, key[T, string]{T(value)}, s)
	set(enums, key[string, T]{s}, value)
	allVals := get(enums, key[T, []T]{0})
	allVals = append(allVals, value)
	set(enums, key[T, []T]{0}, allVals)

	return value
}

// EnumOf returns the corresponding enum for a given string
// representation, and whether it is valid.
func EnumOf[T ~int](s string) (T, bool) {
	enum, ok := get2(enums, key[string, T]{s})
	if !ok {
		return 0, false
	}
	return enum, true
}

// MustEnumOf returns the corresponding enum for a given string representation.
// It panics if the string does not correspond to a valid enum value.
func MustEnumOf[T ~int](s string) T {
	enum, ok := get2(enums, key[string, T]{s})
	if !ok {
		panic(fmt.Sprintf("enum %s: invalid", reflect.TypeFor[T]().Name()))
	}

	return enum
}

// StringOf returns the string representation of an enum value.
func StringOf[T ~int](value T) string {
	enum, ok := get2(enums, key[T, string]{value})
	if !ok {
		return fmt.Sprintf("%s::%s", reflect.TypeFor[T]().Name(), UndefinedString)
	}
	return enum
}

// MustStringOf returns the string representation of an enum value.
// It panics if the enum value is invalid.
func MustStringOf[T ~int](value T) string {
	str := StringOf(value)
	if strings.HasSuffix(str, UndefinedString) {
		panic(fmt.Sprintf("enum %s: invalid value %d", reflect.TypeFor[T]().Name(), value))
	}

	return str
}

// IsValid checks if an enum value is valid.
// It returns true if the enum value is valid, and false otherwise.
func IsValid[T ~int](value T) bool {
	_, ok := get2(enums, key[T, string]{value})
	return ok
}

// MarshalJSON serializes an enum value into its string representation.
// This utility function takes an enum value and converts it into a JSON byte
// slice, representing the enum as a string instead of a numeric value.
//
// Example:
//
//	role := RoleAdmin
//	data, _ := MarshalJSON(role)  // Result: []byte(`"admin"`)
func MarshalJSON[T ~int](value T) ([]byte, error) {
	if !IsValid(value) {
		return nil, fmt.Errorf("unknown %s: %d", reflect.TypeFor[T]().Name(), value)
	}

	return json.Marshal(StringOf(value))
}

// UnmarshalJSON deserializes a string representation of an enum value from
// JSON. This utility function takes a byte slice of JSON data and converts it
// into the corresponding enum value.
//
// Example:
//
//	data := []byte(`"admin"`)
//	var role Role
//	_ := UnmarshalJSON(&role, data)  // role will be set to RoleAdmin
func UnmarshalJSON[T ~int](data []byte, t *T) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	enum, valid := EnumOf[T](str)
	if !valid {
		return fmt.Errorf("unknown %s string: %s", reflect.TypeFor[T]().Name(), str)
	}

	*t = enum
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
func All[T ~int]() []T {
	return get(enums, key[T, []T]{0})
}
