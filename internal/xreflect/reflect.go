package xreflect

import (
	"reflect"
	"slices"
)

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

var (
	intKinds  = []reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64}
	uintKinds = []reflect.Kind{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64}

	intTypes = []reflect.Type{
		reflect.TypeOf(int(0)),
		reflect.TypeOf(int8(0)),
		reflect.TypeOf(int16(0)),
		reflect.TypeOf(int32(0)),
		reflect.TypeOf(int64(0)),
	}
	uintTypes = []reflect.Type{
		reflect.TypeOf(uint(0)),
		reflect.TypeOf(uint8(0)),
		reflect.TypeOf(uint16(0)),
		reflect.TypeOf(uint32(0)),
		reflect.TypeOf(uint64(0)),
	}
)

// IsSignedInt returns true if the value is one of signed integer types.
func IsSignedInt(v any) bool {
	return slices.Contains(intKinds, reflect.TypeOf(v).Kind())
}

// IsUnsignedInt returns true if the value is one of unsigned integer types.
func IsUnsignedInt(v any) bool {
	return slices.Contains(uintKinds, reflect.TypeOf(v).Kind())
}

// IsInt returns true if the value is one of any integer types.
func IsInt(v any) bool {
	return IsSignedInt(v) || IsUnsignedInt(v)
}

// IsFloat32 returns true if the value is a float32.
func IsFloat32(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Float32
}

// IsFloat64 returns true if the value is a float64.
func IsFloat64(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Float64
}

// IsFloat returns true if the value is a float.
func IsFloat(v any) bool {
	return IsFloat32(v) || IsFloat64(v)
}

// IsNumber returns true if the value is a number.
func IsNumber(v any) bool {
	return IsFloat(v) || IsInt(v)
}

// IsString returns true if the value is a string.
func IsString(v any) bool {
	return reflect.String == reflect.TypeOf(v).Kind()
}

// IsPrimitiveSignedInt returns true if the value is one of signed integer types.
func IsPrimitiveSignedInt(v any) bool {
	return slices.Contains(intTypes, reflect.TypeOf(v))
}

// IsPrimitiveUnsignedInt returns true if the value is one of unsigned integer types.
func IsPrimitiveUnsignedInt(v any) bool {
	return slices.Contains(uintTypes, reflect.TypeOf(v))
}

// IsPrimitiveInt returns true if the value is one of any integer types.
func IsPrimitiveInt(v any) bool {
	return IsPrimitiveSignedInt(v) || IsPrimitiveUnsignedInt(v)
}

// IsPrimitiveFloat32 returns true if the value is a float32.
func IsPrimitiveFloat32(v any) bool {
	return reflect.TypeOf(v) == reflect.TypeOf(float32(0))
}

// IsPrimitiveFloat64 returns true if the value is a float64.
func IsPrimitiveFloat64(v any) bool {
	return reflect.TypeOf(v) == reflect.TypeOf(float64(0))
}

// IsPrimitiveFloat returns true if the value is a float.
func IsPrimitiveFloat(v any) bool {
	return IsPrimitiveFloat32(v) || IsPrimitiveFloat64(v)
}

// IsPrimitiveNumber returns true if the value is a number.
func IsPrimitiveNumber(v any) bool {
	return IsPrimitiveFloat(v) || IsPrimitiveInt(v)
}

// IsPrimitiveString returns true if the value is a string.
func IsPrimitiveString(v any) bool {
	return reflect.TypeOf(v) == reflect.TypeOf("")
}

// Convert returns the value converted to type T.
func Convert[T any](value any) T {
	return reflect.ValueOf(value).Convert(reflect.TypeOf((*T)(nil)).Elem()).Interface().(T)
}

// Zero gets the zero value of type T.
func Zero[T any]() T {
	var t T
	return t
}

// IsZeroImplement returns true if the zero value of type T implements interface
// I.
func IsZeroImplement[T, I any]() bool {
	_, ok := any(Zero[T]()).(I)
	return ok
}

// IsImplement returns true if value v implements interface I.
func IsImplement[I any](v any) bool {
	_, ok := v.(I)
	return ok
}

// ImplementZero converts the zero value of type T to interface I.
func ImplementZero[T, I any]() I {
	return any(Zero[T]()).(I)
}

func IsExported[T any]() bool {
	name := reflect.TypeOf((*T)(nil)).Elem().Name()
	if len(name) == 0 {
		return false
	}

	return name[0] >= 'A' && name[0] <= 'Z'
}
