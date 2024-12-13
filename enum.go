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
)

// innerEnumable is an internal interface used for handling centralized
// initialization via New function.
type innerEnumable interface {
	newInnerEnum(s string) any
}

// Map associates an enum with its numeric and string representations. If the
// enum is an integer, its value will be used as the numeric representation.
// Otherwise, the library automatically assigns the smallest positive number
// available to the enum.
//
// Note that this function is not thread-safe. Ensure mappings are set up during
// initialization or other safe execution points to avoid race conditions.
func Map[Enum any](enum Enum, s string) Enum {
	var id int64

	switch reflect.TypeOf(enum).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		id = reflect.ValueOf(enum).Convert(reflect.TypeOf(int64(0))).Interface().(int64)
	default:
		id = core.GetAvailableEnumValue[Enum]()
	}

	return core.MapAny(id, enum, s)
}

// New creates a dynamic enum value. The Enum type must be an int, string, or
// supported enumable (e.g WrapEnum, SafeEnum).
//
// This function provides a convenient way to define and map an enum value
// without creating a constant explicitly.
//
// Note:
//   - Enums created with this function are variables, not constants. If you
//     need a constant enum, declare it explicitly and use enum.Map() instead.
//   - This function is not thread-safe and should only be called during
//     initialization or other safe execution points to avoid race conditions.
func New[Enum any](s string) Enum {
	var enum Enum
	if enumer, ok := any(enum).(innerEnumable); ok {
		return enumer.newInnerEnum(s).(Enum)
	}

	id := core.GetAvailableEnumValue[Enum]()

	switch reflect.TypeOf(enum).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		enum = reflect.ValueOf(id).Convert(reflect.TypeOf((*Enum)(nil)).Elem()).Interface().(Enum)
	case reflect.String:
		enum = reflect.ValueOf(s).Convert(reflect.TypeOf((*Enum)(nil)).Elem()).Interface().(Enum)
	default:
		panic("invalid enum type: require integer, string, or innerEnumable")
	}

	return core.MapAny(id, enum, s)
}

// NewExtended helps to initialize the extended enum.
//
// Note that this function is not thread-safe and should only be called during
// initialization or other safe execution points to avoid race conditions.
func NewExtended[T innerEnumable](s string) T {
	var extendEnum T

	extendEnumValue := reflect.ValueOf(&extendEnum).Elem()

	for i := 0; i < extendEnumValue.NumField(); i++ {
		fieldType := reflect.TypeOf(extendEnum).Field(i)
		if !fieldType.Anonymous {
			continue
		}

		if !fieldType.Type.Implements(reflect.TypeOf((*innerEnumable)(nil)).Elem()) {
			continue
		}

		fieldValue := extendEnumValue.FieldByName(fieldType.Name)

		inner := fieldValue.Interface().(innerEnumable).newInnerEnum(s)
		fieldValue.Set(reflect.ValueOf(inner))

		return core.MapAny(core.GetAvailableEnumValue[T](), extendEnum, s)
	}

	panic("NewExtended is only used to dynamically create an extended enum, please use New instead")
}

// Finalize prevents any further creation of new enum values.
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
