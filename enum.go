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
	"fmt"
	"math"
	"reflect"
)

var id2StrMap = make(map[any]map[int]string)
var str2EnumMap = make(map[any]map[string]int)

func getAvaiableID[T ~int]() T {
	undefined := UndefinedOf[T]()
	if _, ok := id2StrMap[undefined]; !ok {
		return 0
	}

	id := 0
	for {
		if _, ok := id2StrMap[undefined][id]; !ok {
			break
		}
		id++
	}

	return T(id)
}

func UndefinedOf[T ~int]() T {
	return math.MaxInt32
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
func New[T ~int](val string) T {
	e := getAvaiableID[T]()
	return Map(e, val)
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
func Map[T ~int](enum T, val string) T {
	undefined := UndefinedOf[T]()
	if _, ok := id2StrMap[undefined]; !ok {
		id2StrMap[undefined] = make(map[int]string)
		str2EnumMap[undefined] = make(map[string]int)
	}

	if StringOf(enum) != StringOf(undefined) {
		panic("do not map a mapped enum")
	}

	id2StrMap[undefined][int(enum)] = val
	str2EnumMap[undefined][val] = int(enum)

	return enum
}

// EnumOf returns the corresponding enum for a given string representation.
func EnumOf[T ~int](s string) T {
	undefined := UndefinedOf[T]()
	enum, ok := str2EnumMap[undefined]
	if !ok {
		return undefined
	}

	enumValue, ok := enum[s]
	if !ok {
		return undefined
	}

	return T(enumValue)
}

// MustEnumOf returns the corresponding enum for a given string representation.
// It panics if the string does not correspond to a valid enum value.
func MustEnumOf[T ~int](s string) T {
	undefined := UndefinedOf[T]()
	enum, ok := str2EnumMap[undefined]
	if !ok {
		panic(fmt.Sprintf("enum %s: invalid", reflect.TypeFor[T]().Name()))
	}

	enumValue, ok := enum[s]
	if !ok {
		panic(fmt.Sprintf("enum %s: invalid string representation %s", reflect.TypeFor[T]().Name(), s))
	}

	return T(enumValue)
}

// StringOf returns the string representation of an enum value.
func StringOf[T ~int](id T) string {
	t := UndefinedOf[T]()

	enum, ok := id2StrMap[t]
	if !ok {
		return fmt.Sprintf("%s::<<undefined>>", reflect.TypeOf(t).Name())
	}

	enumNumber, ok := enum[int(id)]
	if !ok {
		return fmt.Sprintf("%s::<<undefined>>", reflect.TypeOf(t).Name())
	}

	return enumNumber
}

// MustStringOf returns the string representation of an enum value.
// It panics if the enum value is invalid.
func MustStringOf[T ~int](id T) string {
	t := UndefinedOf[T]()

	enum, ok := id2StrMap[t]
	if !ok {
		panic(fmt.Sprintf("enum %s: invalid", reflect.TypeFor[T]().Name()))
	}

	enumNumber, ok := enum[int(id)]
	if !ok {
		panic(fmt.Sprintf("enum %s: invalid %d", reflect.TypeFor[T]().Name(), id))
	}

	return enumNumber
}
