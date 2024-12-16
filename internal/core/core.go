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

// MapAny map the enum value to the enum system.
func MapAny[Enum any](num int64, value Enum, s string) Enum {
	if s == "" {
		panic("not allow empty string representation in enum definition")
	}

	if s == "<nil>" {
		panic("not allow \"<nil>\" string representation in enum definition")
	}

	if ok := mtmap.MustGet(mtkey.IsFinalized[Enum]()); ok {
		panic("enum is finalized")
	}

	if _, ok := mtmap.Get(mtkey.Int2Enum[Enum](num)); ok {
		panic("duplicate enum number is not allowed")
	}

	if _, ok := mtmap.Get(mtkey.String2Enum[Enum](s)); ok {
		panic("duplicate enum string is not allowed")
	}

	if _, ok := mtmap.Get(mtkey.Enum2Int(value)); ok {
		panic("duplicate enum is not allowed")
	}

	mtmap.Set(mtkey.Enum2String(value), s)
	mtmap.Set(mtkey.Enum2Int(value), num)
	mtmap.Set(mtkey.String2Enum[Enum](s), value)
	mtmap.Set(mtkey.Int2Enum[Enum](num), value)

	allVals := mtmap.MustGet(mtkey.AllEnums[Enum]())
	allVals = append(allVals, value)
	mtmap.Set(mtkey.AllEnums[Enum](), allVals)

	return value
}
