package mtkey

type enum2String[T any] struct{ key T }

func (enum2String[T]) InferValue() string { panic("not implemented") }

func Enum2String[T any](enum T) enum2String[T] {
	return enum2String[T]{key: enum}
}

type enum2Int[T any] struct{ key T }

func (enum2Int[T]) InferValue() int64 { panic("not implemented") }

func Enum2Int[T any](enum T) enum2Int[T] {
	return enum2Int[T]{key: enum}
}

type string2Enum[T any] struct{ key string }

func (string2Enum[T]) InferValue() T { panic("not implemented") }

func String2Enum[T any](s string) string2Enum[T] {
	return string2Enum[T]{key: s}
}

type num2Enum[T any] struct{ key int64 }

func (num2Enum[T]) InferValue() T { panic("not implemented") }

func Int2Enum[T any](key int64) num2Enum[T] {
	return num2Enum[T]{key: key}
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
