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
	"encoding/xml"
	"fmt"
	"math"
	"reflect"

	"github.com/xybor-x/enum/internal/core"
	"github.com/xybor-x/enum/internal/mtkey"
	"github.com/xybor-x/enum/internal/mtmap"
	"github.com/xybor-x/enum/internal/xreflect"
	"gopkg.in/yaml.v3"
)

// newableEnum is an internal interface used for handling centralized
// initialization via New function.
type newableEnum interface {
	// newEnum creates a dynamic enum value of the current type and map it into
	// the enum system.
	newEnum(reprs []any) any
}

// hookAfterEnum calls hookAfter() method after the enum is created.
type hookAfterEnum interface {
	hookAfter()
}

// Map associates an enum with its representations under strict rules:
//   - String enums map to themselves as the string representation; Stringer is
//     also treated as a string representation if no string repr is found.
//   - Numeric enums map to themselves as the numeric representation; all
//     primitive numeric types (ints, uints, floats) are treated as a single
//     type.
//   - An enum cannot be mapped to multiple representations of the same type.
//
// Note that this function is not thread-safe and should only be called during
// initialization or other safe execution points to avoid race conditions.
func Map[Enum any](enum Enum, reprs ...any) Enum {
	defer func() {
		if hook, ok := any(enum).(hookAfterEnum); ok {
			hook.hookAfter()
		}
	}()

	return core.MapAny(enum, reprs)
}

// New creates a dynamic enum value then mapped to its representations. The Enum
// type must be a number, string, or supported enums (e.g WrapEnum, SafeEnum).
//
// If the enum is
//   - Supported enum: the inner new function will be called to generate the
//     enum value.
//   - Number: the numeric representation will be assigned to the enum value.
//   - String: the string representation will be assigned to the enum value.
//   - Other cases, panics.
//
// Note that this function is not thread-safe and should only be called during
// initialization or other safe execution points to avoid race conditions.
func New[Enum any](reprs ...any) (enum Enum) {
	defer func() {
		if hook, ok := any(enum).(hookAfterEnum); ok {
			hook.hookAfter()
		}
	}()

	switch {
	case xreflect.IsZeroImplement[Enum, newableEnum]():
		return xreflect.ImplementZero[Enum, newableEnum]().newEnum(reprs).(Enum)

	case xreflect.IsNumber(xreflect.Zero[Enum]()):
		// The numeric representation will be used as the the enum value.
		numericRepr := core.GetNumericRepresentation(reprs)
		if numericRepr == nil {
			numericRepr = core.GetAvailableEnumValue[Enum]()
		}

		return core.MapAny(xreflect.Convert[Enum](numericRepr), core.RemoveNumericRepresentation(reprs))

	case xreflect.IsString(xreflect.Zero[Enum]()):
		// The string representation will be used as the the enum value.
		strRepr, ok := core.GetStringRepresentation(reprs)
		if !ok {
			panic(fmt.Sprintf("enum %s: new a string enum must provide its string representation", TrueNameOf[Enum]()))
		}

		return core.MapAny(xreflect.Convert[Enum](strRepr), core.RemoveStringRepresentation(reprs))

	default:
		// TODO: For the Enum type, I want to use type constraints to allow only
		// numbers, strings, and innerEnumable. However, type constraints
		// currently prevent combining unions with interfaces.
		panic("invalid enum type: require integer, string, or innerEnumable, otherwise use Map instead!")
	}
}

// NewExtended initializes an extended enum then mapped to its representations.
//
// An extended enum follows this structure (the embedded Enum must be an
// anonymous field to inherit its built-in methods):
//
//	type role any
//	type Role struct { enum.SafeEnum[role] }
//
// Note that this function is not thread-safe and should only be called during
// initialization or other safe execution points to avoid race conditions.
func NewExtended[T newableEnum](reprs ...any) (enum T) {
	defer func() {
		if hook, ok := any(enum).(hookAfterEnum); ok {
			hook.hookAfter()
		}
	}()

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
		if !fieldType.Type.Implements(reflect.TypeOf((*newableEnum)(nil)).Elem()) {
			continue
		}

		if core.GetNumericRepresentation(reprs) == nil {
			reprs = append(reprs, core.GetAvailableEnumValue[T]())
		}

		// Set value to the embedded enumable field.
		enumField := extendEnumValue.FieldByName(fieldType.Name)
		enumField.Set(reflect.ValueOf(enumField.Interface().(newableEnum).newEnum(reprs)))

		// The newEnum method mapped the enum value to the system (see the
		// description of the newEnum method). Why is MapAny called again here?
		//
		// The mapping in the newEnum method only applies the enum value to the
		// embedded enum field type, not the extended enum type. To enable
		// utility functions to work with the extended enum type, we need to map
		// it again using MapAny.
		return core.MapAny(extendEnum, reprs)
	}

	panic("invalid enum type: NewExtended is only used to create an extended enum, otherwise use New or Map instead!")
}

// Finalize prevents the creation of any new enum values for the current type.
func Finalize[Enum any]() bool {
	mtmap.Set(mtkey.IsFinalized[Enum](), true)
	return true
}

// FromInt returns the corresponding enum for a given int representation, and
// whether it is valid.
//
// DEPRECATED: Use FromNumber instead.
func FromInt[Enum any](i int) (Enum, bool) {
	return FromNumber[Enum](i)
}

// FromNumber returns the corresponding enum for a given number representation,
// and whether it is valid.
func FromNumber[Enum any, N xreflect.Number](n N) (Enum, bool) {
	return From[Enum](n)
}

// MustFromInt returns the corresponding enum for a given int representation.
//
// It returns zero value if the enum value is invalid.
//
// DEPRECATED: Use MustFromNumber instead.
func MustFromInt[Enum any](i int) Enum {
	t, _ := FromInt[Enum](i)
	return t
}

// MustFromNumber returns the corresponding enum for a given number
// representation.
//
// It returns the zero value if the enum value is invalid.
func MustFromNumber[Enum any, N xreflect.Number](n N) Enum {
	t, _ := FromNumber[Enum](n)
	return t
}

// FromString returns the corresponding enum for a given string representation,
// and whether it is valid.
func FromString[Enum any](s string) (Enum, bool) {
	return From[Enum](s)
}

// MustFromString returns the corresponding enum for a given string
// representation.
//
// It returns zero value if the string does not correspond to a valid enum
// value.
func MustFromString[Enum any](s string) Enum {
	enum, _ := FromString[Enum](s)
	return enum
}

// ToString returns the string representation of the given enum value. It
// returns <nil> for invalid enums.
func ToString[Enum any](value Enum) string {
	str, ok := To[string](value)
	if !ok {
		return "<nil>"
	}

	return str
}

// ToInt returns the int representation for the given enum value. It returns the
// smallest value of int (math.MinInt32) for invalid enums.
//
// DEPRECATED: This function returns math.MinInt32 for invalid enums, which may
// cause unexpected behavior. Use To() or MustTo() instead.
func ToInt[Enum any](enum Enum) int {
	value, ok := To[int](enum)
	if !ok {
		return math.MinInt32
	}

	return value
}

// From returns the corresponding enum for a given representation, and whether
// it is valid.
func From[Enum any, P any](a P) (Enum, bool) {
	return mtmap.Get2(mtkey.Repr2Enum[Enum](a))
}

// MustFrom returns the corresponding enum for a given representation. It
// returns the zero value of enum in case the representation is unknown.
func MustFrom[Enum any, P any](a P) Enum {
	e, _ := mtmap.Get2(mtkey.Repr2Enum[Enum](a))
	return e
}

// To returns the representation (the type is relied on P type parameter) for
// the given enum value. The latter returned value is false if the enum is
// invalid or the enum doesn't have any representation of type P.
func To[P, Enum any](enum Enum) (P, bool) {
	ret, ok := mtmap.Get2(mtkey.Enum2Repr[Enum, P](enum))
	if !ok {
		return xreflect.Zero[P](), false
	}

	return ret.(P), true
}

// MustTo returns the representation (the type is relied on P type parameter)
// for the given enum value. It returns zero value if the enum is invalid or the
// enum doesn't have any representation of type P.
func MustTo[P, Enum any](enum Enum) P {
	val, _ := To[P](enum)
	return val
}

// IsValid checks if an enum value is valid. It returns true if the enum value
// is valid, and false otherwise.
func IsValid[Enum any](value Enum) bool {
	_, ok := mtmap.Get2(mtkey.Enum2Repr[Enum, string](value))
	return ok
}

// MarshalJSON serializes an enum value into its string representation.
func MarshalJSON[Enum any](value Enum) ([]byte, error) {
	s, ok := mtmap.Get2(mtkey.Enum2JSON(value))
	if !ok {
		return nil, fmt.Errorf("enum %s: invalid value %#v", TrueNameOf[Enum](), value)
	}

	return []byte(s), nil
}

// UnmarshalJSON deserializes a string representation of an enum value from
// JSON.
func UnmarshalJSON[Enum any](data []byte, t *Enum) (err error) {
	n := len(data)
	if n < 2 || data[0] != '"' || data[n-1] != '"' {
		return fmt.Errorf("enum %s: invalid string %s", TrueNameOf[Enum](), string(data))
	}

	enum, ok := mtmap.Get2(mtkey.Repr2Enum[Enum](string(data[1 : n-1])))
	if !ok {
		return fmt.Errorf("enum %s: unknown string %s", TrueNameOf[Enum](), string(data[1:n-1]))
	}

	*t = enum
	return nil
}

// MarshalYAML serializes an enum value into its string representation.
func MarshalYAML[Enum any](value Enum) (any, error) {
	s, ok := mtmap.Get2(mtkey.Enum2Repr[Enum, string](value))
	if !ok {
		return nil, fmt.Errorf("enum %s: invalid value %#v", TrueNameOf[Enum](), value)
	}

	return s, nil
}

// UnmarshalYAML deserializes a string representation of an enum value from
// YAML.
func UnmarshalYAML[Enum any](value *yaml.Node, t *Enum) error {
	// Check if the value is a scalar (string in this case)
	if value.Kind != yaml.ScalarNode {
		return fmt.Errorf("enum %s: only supports scalar in yaml enum", TrueNameOf[Enum]())
	}

	// Assign the string value directly
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}

	var ok bool
	*t, ok = From[Enum](s)
	if !ok {
		return fmt.Errorf("enum %s: unknown string %s", TrueNameOf[Enum](), s)
	}

	return nil
}

// MarshalXML converts enum to its string representation.
func MarshalXML[Enum any](encoder *xml.Encoder, start xml.StartElement, enum Enum) error {
	str, ok := To[string](enum)
	if !ok {
		return fmt.Errorf("enum %s: invalid value %#v", TrueNameOf[Enum](), enum)
	}

	if start.Name.Local == "" {
		start.Name.Local = NameOf[Enum]()
	}

	return encoder.EncodeElement(str, start)
}

// UnmarshalXML parses the string representation back into an enum.
func UnmarshalXML[Enum any](decoder *xml.Decoder, start xml.StartElement, enum *Enum) error {
	var str string
	if err := decoder.DecodeElement(&str, &start); err != nil {
		return err
	}

	val, ok := FromString[Enum](str)
	if !ok {
		return fmt.Errorf("enum %s: unknown string %s", TrueNameOf[Enum](), str)
	}

	*enum = val
	return nil
}

// ValueSQL serializes an enum into a database-compatible format.
func ValueSQL[Enum any](value Enum) (driver.Value, error) {
	str, ok := mtmap.Get2(mtkey.Enum2Repr[Enum, string](value))
	if !ok {
		return nil, fmt.Errorf("enum %s: invalid value %#v", TrueNameOf[Enum](), value)
	}

	return str, nil
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
		return fmt.Errorf("enum %s: not support type %s", TrueNameOf[Enum](), reflect.TypeOf(a))
	}

	enum, ok := mtmap.Get2(mtkey.Repr2Enum[Enum](data))
	if !ok {
		return fmt.Errorf("enum %s: unknown string %s", TrueNameOf[Enum](), data)
	}

	*value = enum
	return nil
}

// All returns a slice containing all enum values of a specific type.
func All[Enum any]() []Enum {
	return mtmap.Get(mtkey.AllEnums[Enum]())
}

// NameOf returns the name of the enum type. In case of this is an advanced enum
// provided by this library, this function returns the only underlying enum
// name, which differs from TrueNameOf.
//
// For example:
//
//	NameOf[Role]()           = "Role"
//	NameOf[WrapEnum[role]]() = "Role"
func NameOf[T any]() string {
	return core.NameOf[T]()
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
	return core.TrueNameOf[T]()
}

// mustHaveUnderlyingRepr ensures an enum has a representation of its underlying
// type.
func mustHaveUnderlyingRepr[underlyingEnum, Enum any](e Enum) {
	if !IsValid(e) {
		return
	}

	mapUnderlying[underlyingEnum](e)

	if _, ok := To[underlyingEnum](e); !ok {
		panic(fmt.Sprintf("enum %s (%#v): require a representation of %T",
			TrueNameOf[Enum](), e, xreflect.Zero[underlyingEnum]()))
	}
}

// mapUnderlying maps the enum to underlying enum in case the underlying enum
// is a string or numeric type. It ignores cases where the underlying type is
// exported and define at least one method.
func mapUnderlying[underlyingEnum, Enum any](enum Enum) {
	if reflect.TypeOf((*underlyingEnum)(nil)).Elem().NumMethod() > 0 || xreflect.IsExported[underlyingEnum]() {
		return
	}

	var repr underlyingEnum
	switch {
	case xreflect.IsSignedInt(repr):
		repr = xreflect.Convert[underlyingEnum](MustTo[int64](enum))
	case xreflect.IsUnsignedInt(repr):
		repr = xreflect.Convert[underlyingEnum](MustTo[uint64](enum))
	case xreflect.IsFloat32(repr):
		repr = xreflect.Convert[underlyingEnum](MustTo[float32](enum))
	case xreflect.IsFloat64(repr):
		repr = xreflect.Convert[underlyingEnum](MustTo[float64](enum))
	case xreflect.IsString(repr):
		repr = xreflect.Convert[underlyingEnum](MustTo[string](enum))
	default:
		str := MustTo[string](enum)
		if reflect.TypeOf(str).ConvertibleTo(reflect.TypeOf((*underlyingEnum)(nil)).Elem()) {
			repr = xreflect.Convert[underlyingEnum](str)
		} else {
			// Ignore if the underlying enum is not a string or numeric type.
			return
		}
	}

	mtmap.Set(mtkey.Repr2Enum[Enum](repr), enum)
	mtmap.Set(mtkey.Enum2Repr[Enum, underlyingEnum](enum), any(repr))
}
