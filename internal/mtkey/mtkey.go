package mtkey

import "github.com/xybor-x/enum/internal/xreflect"

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
