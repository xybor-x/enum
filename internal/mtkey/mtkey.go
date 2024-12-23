package mtkey

import (
	"reflect"

	"github.com/xybor-x/enum/internal/xreflect"
)

type enum2String[T any] struct{ key T }

func (enum2String[T]) InferValue() string { panic("not implemented") }

func Enum2String[T any](enum T) enum2String[T] {
	return enum2String[T]{key: enum}
}

type enum2Number[T any, N xreflect.Number] struct{ key T }

func (enum2Number[T, N]) InferValue() N { panic("not implemented") }

func Enum2Number[T any, N xreflect.Number](enum T) enum2Number[T, N] {
	return enum2Number[T, N]{key: enum}
}

type string2Enum[T any] struct{ key string }

func (string2Enum[T]) InferValue() T { panic("not implemented") }

func String2Enum[T any](s string) string2Enum[T] {
	return string2Enum[T]{key: s}
}

type number2Enum[T any] struct{ key any }

func (number2Enum[T]) InferValue() T { panic("not implemented") }

func Number2Enum[N xreflect.Number, T any](key N) number2Enum[T] {
	return number2Enum[T]{key: key}
}

func AnyNumber2Enum[T any](key any) number2Enum[T] {
	return number2Enum[T]{key: key}
}

type allEnums[T any] struct{}

func (allEnums[T]) InferValue() []T { panic("not implemented") }

func AllEnums[T any]() allEnums[T] {
	return allEnums[T]{}
}

type isFinalized[T any] struct{}

func (isFinalized[T]) InferValue() bool { panic("not implemented") }

func IsFinalized[T any]() isFinalized[T] {
	return isFinalized[T]{}
}

type nameOf[T any] struct{}

func (nameOf[T]) InferValue() string { panic("not implemented") }

func NameOf[T any]() nameOf[T] {
	return nameOf[T]{}
}

type trueNameOf[T any] struct{}

func (trueNameOf[T]) InferValue() string { panic("not implemented") }

func TrueNameOf[T any]() trueNameOf[T] {
	return trueNameOf[T]{}
}

type enum2JSON[T any] struct{ key T }

func (enum2JSON[T]) InferValue() string { panic("not implemented") }

func Enum2JSON[T any](key T) enum2JSON[T] {
	return enum2JSON[T]{key: key}
}

type enum2Extra[T any] struct {
	key T
	typ reflect.Type
}

func (enum2Extra[T]) InferValue() any { panic("not implemented") }

func Enum2Extra[T, P any](key T) enum2Extra[T] {
	return enum2Extra[T]{key: key, typ: reflect.TypeOf(xreflect.Zero[P]())}
}

func Enum2ExtraWith[T any](key T, extra any) enum2Extra[T] {
	return enum2Extra[T]{key: key, typ: reflect.TypeOf(extra)}
}

type extra2Enum[T any] struct{ key any }

func (extra2Enum[T]) InferValue() T { panic("not implemented") }

func Extra2Enum[T any](key any) extra2Enum[T] {
	return extra2Enum[T]{key: key}
}
