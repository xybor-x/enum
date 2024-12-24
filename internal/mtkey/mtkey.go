package mtkey

import (
	"reflect"

	"github.com/xybor-x/enum/internal/xreflect"
)

type allEnums[Enum any] struct{}

func (allEnums[Enum]) InferValue() []Enum { panic("not implemented") }

func AllEnums[Enum any]() allEnums[Enum] {
	return allEnums[Enum]{}
}

type isFinalized[Enum any] struct{}

func (isFinalized[Enum]) InferValue() bool { panic("not implemented") }

func IsFinalized[Enum any]() isFinalized[Enum] {
	return isFinalized[Enum]{}
}

type nameOf[Enum any] struct{}

func (nameOf[Enum]) InferValue() string { panic("not implemented") }

func NameOf[Enum any]() nameOf[Enum] {
	return nameOf[Enum]{}
}

type trueNameOf[Enum any] struct{}

func (trueNameOf[Enum]) InferValue() string { panic("not implemented") }

func TrueNameOf[Enum any]() trueNameOf[Enum] {
	return trueNameOf[Enum]{}
}

type enum2JSON[Enum any] struct{ key Enum }

func (enum2JSON[Enum]) InferValue() string { panic("not implemented") }

func Enum2JSON[Enum any](key Enum) enum2JSON[Enum] {
	return enum2JSON[Enum]{key: key}
}

type enum2Repr[Enum any] struct {
	key Enum
	typ reflect.Type
}

func (enum2Repr[Enum]) InferValue() any { panic("not implemented") }

func Enum2Repr[Enum, P any](key Enum) enum2Repr[Enum] {
	return enum2Repr[Enum]{key: key, typ: reflect.TypeOf(xreflect.Zero[P]())}
}

func Enum2ReprWith[Enum any](key Enum, extra any) enum2Repr[Enum] {
	return enum2Repr[Enum]{key: key, typ: reflect.TypeOf(extra)}
}

type repr2Enum[Enum any] struct{ key any }

func (repr2Enum[Enum]) InferValue() Enum { panic("not implemented") }

func Repr2Enum[Enum any](key any) repr2Enum[Enum] {
	return repr2Enum[Enum]{key: key}
}
