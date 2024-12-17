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
)

// IsSignedInt returns true if the value is one of signed integer types.
func IsSignedInt[T any]() bool {
	return slices.Contains(intKinds, reflect.TypeOf((*T)(nil)).Elem().Kind())
}

// IsUnsignedInt returns true if the value is one of unsigned integer types.
func IsUnsignedInt[T any]() bool {
	return slices.Contains(uintKinds, reflect.TypeOf((*T)(nil)).Elem().Kind())
}

// IsInt returns true if the value is one of any integer types.
func IsInt[T any]() bool {
	return IsSignedInt[T]() || IsUnsignedInt[T]()
}

// IsFloat32 returns true if the value is a float32.
func IsFloat32[T any]() bool {
	return reflect.TypeOf((*T)(nil)).Elem().Kind() == reflect.Float32
}

// IsFloat64 returns true if the value is a float64.
func IsFloat64[T any]() bool {
	return reflect.TypeOf((*T)(nil)).Elem().Kind() == reflect.Float64
}

// IsFloat returns true if the value is a float.
func IsFloat[T any]() bool {
	return IsFloat32[T]() || IsFloat64[T]()
}

// IsNumber returns true if the value is a number.
func IsNumber[T any]() bool {
	return IsFloat[T]() || IsInt[T]()
}

// IsString returns true if the value is a string.
func IsString[T any]() bool {
	return reflect.String == reflect.TypeOf((*T)(nil)).Elem().Kind()
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

// IsImplemented returns true if type T implements interface I.
func IsImplemented[T, I any]() bool {
	_, ok := any(Zero[T]()).(I)
	return ok
}

// ImplementZero converts the zero value of type T to interface I.
func ImplementZero[T, I any]() I {
	return any(Zero[T]()).(I)
}
