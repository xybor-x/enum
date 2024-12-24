package xreflect

import (
	"reflect"
)

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

var (
	IntKinds = map[reflect.Kind]bool{
		reflect.Int:   true,
		reflect.Int8:  true,
		reflect.Int16: true,
		reflect.Int32: true,
		reflect.Int64: true,
	}
	UintKinds = map[reflect.Kind]bool{
		reflect.Uint:   true,
		reflect.Uint8:  true,
		reflect.Uint16: true,
		reflect.Uint32: true,
		reflect.Uint64: true,
	}

	intTypes = map[reflect.Type]bool{
		reflect.TypeOf(int(0)):   true,
		reflect.TypeOf(int8(0)):  true,
		reflect.TypeOf(int16(0)): true,
		reflect.TypeOf(int32(0)): true,
		reflect.TypeOf(int64(0)): true,
	}
	uintTypes = map[reflect.Type]bool{
		reflect.TypeOf(uint(0)):   true,
		reflect.TypeOf(uint8(0)):  true,
		reflect.TypeOf(uint16(0)): true,
		reflect.TypeOf(uint32(0)): true,
		reflect.TypeOf(uint64(0)): true,
	}
)

// IsSignedInt returns true if the value is one of signed integer types.
func IsSignedInt(v any) bool {
	kind, ok := v.(reflect.Kind)
	if !ok {
		if typ := reflect.TypeOf(v); typ != nil {
			kind = typ.Kind()
		}
	}
	return IntKinds[kind]
}

// IsUnsignedInt returns true if the value is one of unsigned integer types.
func IsUnsignedInt(v any) bool {
	kind, ok := v.(reflect.Kind)
	if !ok {
		if typ := reflect.TypeOf(v); typ != nil {
			kind = typ.Kind()
		}
	}
	return UintKinds[kind]
}

// IsInt returns true if the value is one of any integer types.
func IsInt(v any) bool {
	kind, ok := v.(reflect.Kind)
	if !ok {
		if typ := reflect.TypeOf(v); typ != nil {
			kind = typ.Kind()
		}
	}
	return IsSignedInt(kind) || IsUnsignedInt(kind)
}

// IsFloat32 returns true if the value is a float32.
func IsFloat32(v any) bool {
	kind, ok := v.(reflect.Kind)
	if !ok {
		if typ := reflect.TypeOf(v); typ != nil {
			kind = typ.Kind()
		}
	}
	return kind == reflect.Float32
}

// IsFloat64 returns true if the value is a float64.
func IsFloat64(v any) bool {
	kind, ok := v.(reflect.Kind)
	if !ok {
		if typ := reflect.TypeOf(v); typ != nil {
			kind = typ.Kind()
		}
	}
	return kind == reflect.Float64
}

// IsFloat returns true if the value is a float.
func IsFloat(v any) bool {
	kind, ok := v.(reflect.Kind)
	if !ok {
		if typ := reflect.TypeOf(v); typ != nil {
			kind = typ.Kind()
		}
	}
	return IsFloat32(kind) || IsFloat64(kind)
}

// IsNumber returns true if the value is a number.
func IsNumber(v any) bool {
	kind, ok := v.(reflect.Kind)
	if !ok {
		if typ := reflect.TypeOf(v); typ != nil {
			kind = typ.Kind()
		}
	}
	return IsFloat(kind) || IsInt(kind)
}

// IsString returns true if the value is a string.
func IsString(v any) bool {
	kind, ok := v.(reflect.Kind)
	if !ok {
		if typ := reflect.TypeOf(v); typ != nil {
			kind = typ.Kind()
		}
	}
	return reflect.String == kind
}

// IsPrimitiveSignedInt returns true if the value is one of signed integer types.
func IsPrimitiveSignedInt(v any) bool {
	typ, ok := v.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(v)
	}
	return intTypes[typ]
}

// IsPrimitiveUnsignedInt returns true if the value is one of unsigned integer types.
func IsPrimitiveUnsignedInt(v any) bool {
	typ, ok := v.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(v)
	}
	return uintTypes[typ]
}

// IsPrimitiveInt returns true if the value is one of any integer types.
func IsPrimitiveInt(v any) bool {
	typ, ok := v.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(v)
	}
	return IsPrimitiveSignedInt(typ) || IsPrimitiveUnsignedInt(typ)
}

// IsPrimitiveFloat32 returns true if the value is a float32.
func IsPrimitiveFloat32(v any) bool {
	typ, ok := v.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(v)
	}
	return typ == reflect.TypeOf(float32(0))
}

// IsPrimitiveFloat64 returns true if the value is a float64.
func IsPrimitiveFloat64(v any) bool {
	typ, ok := v.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(v)
	}
	return typ == reflect.TypeOf(float64(0))
}

// IsPrimitiveFloat returns true if the value is a float.
func IsPrimitiveFloat(v any) bool {
	typ, ok := v.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(v)
	}
	return IsPrimitiveFloat32(typ) || IsPrimitiveFloat64(typ)
}

// IsPrimitiveNumber returns true if the value is a number.
func IsPrimitiveNumber(v any) bool {
	typ, ok := v.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(v)
	}
	return IsPrimitiveFloat(typ) || IsPrimitiveInt(v)
}

// IsPrimitiveString returns true if the value is a string.
func IsPrimitiveString(v any) bool {
	typ, ok := v.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(v)
	}
	return typ == reflect.TypeOf("")
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
	return len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z'
}
