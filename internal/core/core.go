package core

import (
	"math"
	"strconv"

	"github.com/xybor-x/enum/internal/mtkey"
	"github.com/xybor-x/enum/internal/mtmap"
	"github.com/xybor-x/enum/internal/xmath"
	"github.com/xybor-x/enum/internal/xreflect"
)

func GetAvailableEnumValue[Enum any]() int64 {
	id := int64(0)
	for {
		if _, ok := mtmap.Get2(mtkey.Number2Enum[int64, Enum](id)); !ok {
			break
		}
		id++
	}

	return id
}

// MapAny map the enum value to the enum system.
func MapAny[N xreflect.Number, Enum any](id N, enum Enum, s string) Enum {
	if mtmap.Get(mtkey.IsFinalized[Enum]()) {
		panic("enum is finalized")
	}

	if _, ok := mtmap.Get2(mtkey.Number2Enum[N, Enum](id)); ok {
		panic("duplicate enum number is not allowed")
	}

	if _, ok := mtmap.Get2(mtkey.String2Enum[Enum](s)); ok {
		panic("duplicate enum string is not allowed")
	}

	if _, ok := mtmap.Get2(mtkey.Enum2Number[Enum, N](enum)); ok {
		panic("duplicate enum is not allowed")
	}

	mtmap.Set(mtkey.EnumToJSON(enum), strconv.Quote(s))
	mtmap.Set(mtkey.Enum2String(enum), s)
	mtmap.Set(mtkey.String2Enum[Enum](s), enum)
	mapEnumNumber(enum, id)

	allVals := mtmap.Get(mtkey.AllEnums[Enum]())
	allVals = append(allVals, enum)
	mtmap.Set(mtkey.AllEnums[Enum](), allVals)

	return enum
}

// mapEnumNumber maps the enum to all its number representations (including
// signed and unsigned integers, floating-point numbers) and vice versa.
func mapEnumNumber[Enum any, N xreflect.Number](enum Enum, n N) {
	// Only map the enum to integers if the enum is represented by integer
	// values, where the integer corresponds to the actual numeric value,
	// regardless of the underlying type.
	//
	// For example: float32(3) is also an integer, whereas float32(0.7) is not.
	mapInteger := true
	if xreflect.IsFloat32[N]() {
		mapInteger = xreflect.Convert[float32](n) == xmath.Trunc32(xreflect.Convert[float32](n))
	} else if xreflect.IsFloat64[N]() {
		mapInteger = xreflect.Convert[float64](n) == math.Trunc(xreflect.Convert[float64](n))
	}

	if mapInteger {
		// Map enum to all signed integers.
		mtmap.Set(mtkey.Enum2Number[Enum, int](enum), xreflect.Convert[int](n))
		mtmap.Set(mtkey.Enum2Number[Enum, int8](enum), xreflect.Convert[int8](n))
		mtmap.Set(mtkey.Enum2Number[Enum, int16](enum), xreflect.Convert[int16](n))
		mtmap.Set(mtkey.Enum2Number[Enum, int32](enum), xreflect.Convert[int32](n))
		mtmap.Set(mtkey.Enum2Number[Enum, int64](enum), xreflect.Convert[int64](n))

		// Map enum to all unsigned integers.
		mtmap.Set(mtkey.Enum2Number[Enum, uint](enum), xreflect.Convert[uint](n))
		mtmap.Set(mtkey.Enum2Number[Enum, uint8](enum), xreflect.Convert[uint8](n))
		mtmap.Set(mtkey.Enum2Number[Enum, uint16](enum), xreflect.Convert[uint16](n))
		mtmap.Set(mtkey.Enum2Number[Enum, uint32](enum), xreflect.Convert[uint32](n))
		mtmap.Set(mtkey.Enum2Number[Enum, uint64](enum), xreflect.Convert[uint64](n))

		// Map all signed integers to enum.
		mtmap.Set(mtkey.Number2Enum[int, Enum](xreflect.Convert[int](n)), enum)
		mtmap.Set(mtkey.Number2Enum[int8, Enum](xreflect.Convert[int8](n)), enum)
		mtmap.Set(mtkey.Number2Enum[int16, Enum](xreflect.Convert[int16](n)), enum)
		mtmap.Set(mtkey.Number2Enum[int32, Enum](xreflect.Convert[int32](n)), enum)
		mtmap.Set(mtkey.Number2Enum[int64, Enum](xreflect.Convert[int64](n)), enum)

		// Map all unsigned integers to enum.
		mtmap.Set(mtkey.Number2Enum[uint, Enum](xreflect.Convert[uint](n)), enum)
		mtmap.Set(mtkey.Number2Enum[uint8, Enum](xreflect.Convert[uint8](n)), enum)
		mtmap.Set(mtkey.Number2Enum[uint16, Enum](xreflect.Convert[uint16](n)), enum)
		mtmap.Set(mtkey.Number2Enum[uint32, Enum](xreflect.Convert[uint32](n)), enum)
		mtmap.Set(mtkey.Number2Enum[uint64, Enum](xreflect.Convert[uint64](n)), enum)
	}

	// Map enum to all floats.
	mtmap.Set(mtkey.Enum2Number[Enum, float32](enum), xreflect.Convert[float32](n))
	mtmap.Set(mtkey.Enum2Number[Enum, float64](enum), xreflect.Convert[float64](n))

	// Map all floats to enum.
	mtmap.Set(mtkey.Number2Enum[float32, Enum](xreflect.Convert[float32](n)), enum)
	mtmap.Set(mtkey.Number2Enum[float64, Enum](xreflect.Convert[float64](n)), enum)
}
