package core

import (
	"github.com/xybor-x/enum/internal/mtkey"
	"github.com/xybor-x/enum/internal/mtmap"
)

func GetAvailableEnumValue[T any]() int64 {
	id := int64(0)
	for {
		if _, ok := mtmap.Get(mtkey.Int2Enum[T](id)); !ok {
			break
		}
		id++
	}

	return id
}

func MapAny[T any](num int64, value T, s string) T {
	if num < 0 {
		panic("not support negative enum value")
	}

	if s == "" {
		panic("not allow empty string representation in enum definition")
	}

	if s == "<nil>" {
		panic("not allow \"<nil>\" string representation in enum definition")
	}

	if ok := mtmap.MustGet(mtkey.IsFinalized[T]()); ok {
		panic("enum is finalized")
	}

	if _, ok := mtmap.Get(mtkey.Int2Enum[T](num)); ok {
		panic("duplicate enum number is not allowed")
	}

	if _, ok := mtmap.Get(mtkey.String2Enum[T](s)); ok {
		panic("duplicate enum string is not allowed")
	}

	if _, ok := mtmap.Get(mtkey.Enum2Int(value)); ok {
		panic("duplicate enum is not allowed")
	}

	mtmap.Set(mtkey.Enum2String(value), s)
	mtmap.Set(mtkey.Enum2Int(value), num)
	mtmap.Set(mtkey.String2Enum[T](s), value)
	mtmap.Set(mtkey.Int2Enum[T](num), value)

	allVals := mtmap.MustGet(mtkey.AllEnums[T]())
	allVals = append(allVals, value)
	mtmap.Set(mtkey.AllEnums[T](), allVals)

	return value
}
