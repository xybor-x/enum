package mtkey

import "github.com/xybor-x/enum/internal"

type enum2String[T internal.Enumable] struct {
	key T
}

func Enum2String[T internal.Enumable](enum T) enum2String[T] {
	return enum2String[T]{key: enum}
}

func (enum2String[T]) InferValue() string { panic("not implemented") }

type string2Enum[T internal.Enumable] struct {
	key string
}

func String2Enum[T internal.Enumable](s string) string2Enum[T] {
	return string2Enum[T]{key: s}
}

func (string2Enum[T]) InferValue() T { panic("not implemented") }

type allEnums[T internal.Enumable] struct{}

func AllEnums[T internal.Enumable]() allEnums[T] {
	return allEnums[T]{}
}

func (allEnums[T]) InferValue() []T { panic("not implemented") }
