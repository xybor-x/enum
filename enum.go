// Package enum provides a type-safe and powerful enum implementation for Go,
// offering easy conversion between numeric and string representations.
//
// It supports constant enums and integrates seamlessly with Go's `iota` enum pattern.
//
// Features:
//   - No code generation
//   - Supports constant enums
//   - Easy value conversions
//   - Out of the box serialization
package enum

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"

	"github.com/xybor-x/enum/internal/common"
	"github.com/xybor-x/enum/internal/core"
	"github.com/xybor-x/enum/internal/mtkey"
	"github.com/xybor-x/enum/internal/mtmap"
)

// New creates a dynamic enum value and maps it to a string representation.
//
// This function provides a convenient way to define and map an enum value
// without creating a constant explicitly. It is a shorthand for manually
// defining a constant and calling enum.Map().
//
// Note:
//   - Enums created with this function are variables, not constants. If you
//     need a constant enum, declare it explicitly and use enum.Map() instead.
//   - This function is not thread-safe and should only be called during
//     initialization or other safe execution points to avoid race conditions.
func New[T common.Integer](s string) T {
	id := core.GetAvailableEnumValue[T]()
	return core.MapAny(id, T(id), s)
}

// Map associates an enum with its numeric and string representations. If the
// enum is a number, its value will be used as the numeric representation.
// Otherwise, the library automatically assigns the smallest available number
// to the enum.
//
// Note that this function is not thread-safe. Ensure mappings are set up during
// initialization or other safe execution points to avoid race conditions.
func Map[T any](enum T, s string) T {
	var num int64

	var anyt any = enum
	switch v := anyt.(type) {
	case int:
		num = int64(v)
	case int8:
		num = int64(v)
	case int16:
		num = int64(v)
	case int32:
		num = int64(v)
	case int64:
		num = int64(v)
	case uint:
		num = int64(v)
	case uint8:
		num = int64(v)
	case uint16:
		num = int64(v)
	case uint32:
		num = int64(v)
	case uint64:
		num = int64(v)
	default:
		num = core.GetAvailableEnumValue[T]()
	}

	return core.MapAny(num, enum, s)
}

// Finalize prevents any further creation of new enum values.
func Finalize[T any]() bool {
	mtmap.Set(mtkey.IsFinalized[T](), true)
	return true
}

// FromInt returns the corresponding enum for a given int representation, and
// whether it is valid.
func FromInt[T any](i int) (T, bool) {
	return mtmap.Get(mtkey.Int2Enum[T](int64(i)))
}

// MustFromInt returns the corresponding enum for a given int representation.
//
// It panics if the enum value is invalid.
func MustFromInt[T any](i int) T {
	t, ok := FromInt[T](i)
	if !ok {
		panic(fmt.Sprintf("enum %s: invalid int %d", common.NameOf[T](), i))
	}

	return t
}

// FromString returns the corresponding enum for a given string representation,
// and whether it is valid.
func FromString[T any](s string) (T, bool) {
	return mtmap.Get(mtkey.String2Enum[T](s))
}

// MustFromString returns the corresponding enum for a given string
// representation.
//
// It panics if the string does not correspond to a valid enum value.
func MustFromString[T any](s string) T {
	enum, ok := FromString[T](s)
	if !ok {
		panic(fmt.Sprintf("enum %s: invalid string %s", common.NameOf[T](), s))
	}

	return enum
}

// ToString returns the string representation of the given enum value. It
// returns <nil> for invalid enums.
func ToString[T any](value T) string {
	str, ok := mtmap.Get(mtkey.Enum2String(value))
	if !ok {
		return "<nil>"
	}

	return str
}

// ToInt returns the int representation for the given enum value. It returns the
// smallest value of int (math.MinInt32) for invalid enums.
func ToInt[T any](enum T) int {
	value, ok := mtmap.Get(mtkey.Enum2Int(enum))
	if !ok {
		return math.MinInt32
	}

	return int(value)
}

// IsValid checks if an enum value is valid.
// It returns true if the enum value is valid, and false otherwise.
func IsValid[T any](value T) bool {
	_, ok := mtmap.Get(mtkey.Enum2String(value))
	return ok
}

// MarshalJSON serializes an enum value into its string representation.
func MarshalJSON[T any](value T) ([]byte, error) {
	if !IsValid(value) {
		return nil, fmt.Errorf("enum %s: invalid %#v", common.NameOf[T](), value)
	}

	return json.Marshal(ToString(value))
}

// UnmarshalJSON deserializes a string representation of an enum value from
// JSON.
func UnmarshalJSON[T any](data []byte, t *T) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	enum, ok := FromString[T](str)
	if !ok {
		return fmt.Errorf("enum %s: unknown string %s", common.NameOf[T](), str)
	}

	*t = enum
	return nil
}

// ValueSQL serializes an enum into a database-compatible format.
func ValueSQL[T any](value T) (driver.Value, error) {
	if !IsValid(value) {
		return nil, fmt.Errorf("enum %s: invalid %#v", common.NameOf[T](), value)
	}

	return ToString(value), nil
}

// ScanSQL deserializes a database value into an enum type.
func ScanSQL[T any](a any, value *T) error {
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
		return fmt.Errorf("enum %s: unknown string %s", common.NameOf[T](), data)
	}

	*value = enum

	return nil
}

// All returns a slice containing all enum values of a specific type.
func All[T any]() []T {
	return mtmap.MustGet(mtkey.AllEnums[T]())
}
