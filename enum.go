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
	"path"
	"reflect"
	"strings"

	"github.com/xybor-x/enum/internal/core"
	"github.com/xybor-x/enum/internal/mtkey"
	"github.com/xybor-x/enum/internal/mtmap"
	"github.com/xybor-x/enum/internal/xreflect"
)

// innerEnumable is an internal interface used for handling centralized
// initialization via New function.
type innerEnumable interface {
	// newEnum creates an enum value of the current type and map it into
	// the enum system.
	newEnum(id int64, s string) any
}

// Map associates an enum with its numeric and string representations. If the
// enum is an integer, its value will be used as the numeric representation.
// Otherwise, the library automatically assigns the smallest positive number
// available to the enum.
//
// Note that this function is not thread-safe and should only be called during
// initialization or other safe execution points to avoid race conditions.
func Map[Enum any](enum Enum, s string) Enum {
	var id int64
	switch {
	case xreflect.IsInt[Enum]():
		// In case the enum is an integer, it will be used directly as the
		// numeric representation.
		if xreflect.IsUnsignedInt[Enum]() {
			number := xreflect.Convert[uint64](enum)
			if number > math.MaxInt64 {
				// To ensure there is no overflow problem happens while
				// converting type, the library prevents cases of too large enum
				// value.
				//
				// If this causes a problem in your use cases, please raise an
				// issue on GitHub.
				panic("invalid enum value: numeric representation is too large")
			}

			id = int64(number)
		} else {
			id = int64(xreflect.Convert[int64](enum))
		}

	default:
		// Automatically assigns the smallest positive number available to the
		// numeric representation.
		id = core.GetAvailableEnumValue[Enum]()
	}

	return core.MapAny(id, enum, s)
}

// New creates a dynamic enum value. The Enum type must be an int, string, or
// supported enums (e.g WrapEnum, SafeEnum).
//
// The library automatically generates the smallest positive number available as
// the numeric representation of enum.
//
// If the enum is
//   - Integer: the numeric representation will be assigned to the enum value.
//   - String: the parameter s will be assigned to the enum value.
//   - Supported enum: the inner new function will be called to generate the
//     enum value.
//   - Other cases, panics.
//
// Note that this function is not thread-safe and should only be called during
// initialization or other safe execution points to avoid race conditions.
func New[Enum any](s string) Enum {
	id := core.GetAvailableEnumValue[Enum]()

	switch {
	case xreflect.IsInt[Enum]():
		// The numeric representation will be used as the the enum value.
		return core.MapAny(id, xreflect.Convert[Enum](id), s)
	case xreflect.IsString[Enum]():
		// The string representation will be used as the the enum value.
		return core.MapAny(id, xreflect.Convert[Enum](s), s)
	case xreflect.IsImplemented[Enum, innerEnumable]():
		return xreflect.ImplementZero[Enum, innerEnumable]().newEnum(id, s).(Enum)
	default:
		panic("invalid enum type: require integer, string, or innerEnumable, please use Map instead!")
	}
}

// NewExtended initializes an extended enum.
//
// An extended enum follows this structure (the embedded Enum must be an
// anonymous field to inherit its built-in methods):
//
//	type role any
//	type Role struct { enum.SafeEnum[role] }
//
// Note that this function is not thread-safe and should only be called during
// initialization or other safe execution points to avoid race conditions.
func NewExtended[T innerEnumable](s string) T {
	var extendEnum T

	extendEnumValue := reflect.ValueOf(&extendEnum).Elem()

	// Seek the embedded enumable field, then init that field.
	for i := 0; i < extendEnumValue.NumField(); i++ {
		fieldType := reflect.TypeOf(extendEnum).Field(i)

		// Ignore named fields.
		if !fieldType.Anonymous {
			continue
		}

		// Ignore non-enumable fields.
		if !fieldType.Type.Implements(reflect.TypeOf((*innerEnumable)(nil)).Elem()) {
			continue
		}

		id := core.GetAvailableEnumValue[T]()

		// Set value to the embedded enumable field.
		enumField := extendEnumValue.FieldByName(fieldType.Name)
		enumField.Set(reflect.ValueOf(enumField.Interface().(innerEnumable).newEnum(id, s)))

		// The newEnum method mapped the enum value to the system (see the
		// description of the newEnum method). Why is MapAny called again here?
		//
		// The mapping in the newEnum method only applies the enum value to the
		// embedded enum field type, not the extended enum type. To enable
		// utility functions to work with the extended enum type, we need to map
		// it again using MapAny.
		return core.MapAny(id, extendEnum, s)
	}

	panic("invalid enum type: NewExtended is only used to create an extended enum, please use New or Map instead!")
}

// Finalize prevents the creation of any new enum values for the current type.
func Finalize[Enum any]() bool {
	mtmap.Set(mtkey.IsFinalized[Enum](), true)
	return true
}

// FromInt returns the corresponding enum for a given int representation, and
// whether it is valid.
func FromInt[Enum any](i int) (Enum, bool) {
	return mtmap.Get(mtkey.Int2Enum[Enum](int64(i)))
}

// MustFromInt returns the corresponding enum for a given int representation.
//
// It panics if the enum value is invalid.
func MustFromInt[Enum any](i int) Enum {
	t, ok := FromInt[Enum](i)
	if !ok {
		panic(fmt.Sprintf("enum %s: invalid int %d", TrueNameOf[Enum](), i))
	}

	return t
}

// FromString returns the corresponding enum for a given string representation,
// and whether it is valid.
func FromString[Enum any](s string) (Enum, bool) {
	return mtmap.Get(mtkey.String2Enum[Enum](s))
}

// MustFromString returns the corresponding enum for a given string
// representation.
//
// It panics if the string does not correspond to a valid enum value.
func MustFromString[Enum any](s string) Enum {
	enum, ok := FromString[Enum](s)
	if !ok {
		panic(fmt.Sprintf("enum %s: invalid string %s", TrueNameOf[Enum](), s))
	}

	return enum
}

// ToString returns the string representation of the given enum value. It
// returns <nil> for invalid enums.
func ToString[Enum any](value Enum) string {
	str, ok := mtmap.Get(mtkey.Enum2String(value))
	if !ok {
		return "<nil>"
	}

	return str
}

// ToInt returns the int representation for the given enum value. It returns the
// smallest value of int (math.MinInt32) for invalid enums.
func ToInt[Enum any](enum Enum) int {
	value, ok := mtmap.Get(mtkey.Enum2Int(enum))
	if !ok {
		return math.MinInt32
	}

	return int(value)
}

// IsValid checks if an enum value is valid.
// It returns true if the enum value is valid, and false otherwise.
func IsValid[Enum any](value Enum) bool {
	_, ok := mtmap.Get(mtkey.Enum2String(value))
	return ok
}

// MarshalJSON serializes an enum value into its string representation.
func MarshalJSON[Enum any](value Enum) ([]byte, error) {
	if !IsValid(value) {
		return nil, fmt.Errorf("enum %s: invalid value %#v", TrueNameOf[Enum](), value)
	}

	return json.Marshal(ToString(value))
}

// UnmarshalJSON deserializes a string representation of an enum value from
// JSON.
func UnmarshalJSON[Enum any](data []byte, t *Enum) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	enum, ok := FromString[Enum](str)
	if !ok {
		return fmt.Errorf("enum %s: unknown string %s", TrueNameOf[Enum](), str)
	}

	*t = enum
	return nil
}

// ValueSQL serializes an enum into a database-compatible format.
func ValueSQL[Enum any](value Enum) (driver.Value, error) {
	if !IsValid(value) {
		return nil, fmt.Errorf("enum %s: invalid value %#v", TrueNameOf[Enum](), value)
	}

	return ToString(value), nil
}

// ScanSQL deserializes a database value into an enum type.
func ScanSQL[Enum any](a any, value *Enum) error {
	var data string
	switch t := a.(type) {
	case string:
		data = t
	case []byte:
		data = string(t)
	default:
		return fmt.Errorf("not support type %T", a)
	}

	enum, ok := FromString[Enum](data)
	if !ok {
		return fmt.Errorf("enum %s: unknown string %s", TrueNameOf[Enum](), data)
	}

	*value = enum

	return nil
}

// All returns a slice containing all enum values of a specific type.
func All[Enum any]() []Enum {
	return mtmap.MustGet(mtkey.AllEnums[Enum]())
}

var advancedEnumNames = []string{"WrapEnum", "SafeEnum"}

// NameOf returns the name of the enum type. In case of this is an advanced enum
// provided by this library, this function returns the only underlying enum
// name, which differs from TrueNameOf.
//
// For example:
//
//	NameOf[Role]()           = "Role"
//	NameOf[WrapEnum[role]]() = "Role"
func NameOf[T any]() string {
	if name, ok := mtmap.Get(mtkey.NameOf[T]()); ok {
		return name
	}

	name := reflect.TypeOf((*T)(nil)).Elem().Name()
	for _, prefix := range advancedEnumNames {
		if strings.HasPrefix(name, prefix) {
			name = capitalizeFirst(getUnderlyingName(name, prefix))
			break
		}
	}

	mtmap.Set(mtkey.NameOf[T](), name)
	return name
}

// TrueNameOf returns the name of the enum type. In case of this is an advanced
// enum provided by this library, this function returns the full name, which
// differs from NameOf.
//
// For example:
//
//	TrueNameOf[Role]()           = "Role"
//	TrueNameOf[WrapEnum[role]]() = "WrapEnum[role]"
func TrueNameOf[T any]() string {
	if name, ok := mtmap.Get(mtkey.TrueNameOf[T]()); ok {
		return name
	}

	name := reflect.TypeOf((*T)(nil)).Elem().Name()
	for _, prefix := range advancedEnumNames {
		if strings.HasPrefix(name, prefix) {
			name = fmt.Sprintf("%s[%s]", prefix, getUnderlyingName(name, prefix))
			break
		}
	}

	mtmap.Set(mtkey.TrueNameOf[T](), name)
	return name
}

func getUnderlyingName(name, prefix string) string {
	// name = prefix[path/to/module.underlying·id]
	inner := name[len(prefix)+1 : len(name)-1] // inner = path/to/module.underlying·id
	_, inner = path.Split(inner)               // inner = module.underlying·id

	parts := strings.Split(inner, ".")
	inner = parts[len(parts)-1] // inner = underlying·id

	parts = strings.Split(inner, string(rune(183))) // middle dot character.
	return parts[0]                                 // parts[0] = underlying
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s // Return empty string if input is empty
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}
