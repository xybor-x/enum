package core

import (
	"fmt"
	"math"
	"reflect"
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

func GetNumericRepresentation(reprs []any) any {
	var numericRepr any

	for _, repr := range reprs {
		switch {
		case xreflect.IsPrimitiveNumber(repr):
			numericRepr = repr

		default:
			if numericRepr == nil && xreflect.IsNumber(repr) {
				numericRepr = repr
			}
		}
	}

	return numericRepr
}

func GetStringRepresentation(reprs []any) (string, bool) {
	var strRepr string
	var hasStrRepr bool

	for _, repr := range reprs {
		switch {
		case xreflect.IsPrimitiveString(repr):
			strRepr = xreflect.Convert[string](repr)
			hasStrRepr = true

		default:
			if !hasStrRepr {
				if xreflect.IsImplement[fmt.Stringer](repr) {
					strRepr = repr.(fmt.Stringer).String()
					hasStrRepr = true
				} else if xreflect.IsString(repr) {
					strRepr = xreflect.Convert[string](repr)
					hasStrRepr = true
				}
			}
		}
	}

	return strRepr, hasStrRepr
}

func RemoveStringRepresentation(reprs []any) []any {
	strReprIdx := -1

	for i, repr := range reprs {
		switch {
		case xreflect.IsPrimitiveString(repr):
			strReprIdx = i
		}
	}

	if strReprIdx == -1 {
		return reprs
	}

	return append(reprs[:strReprIdx], reprs[strReprIdx+1:]...)
}

func RemoveNumericRepresentation(reprs []any) []any {
	strReprIdx := -1

	for i, repr := range reprs {
		switch {
		case xreflect.IsPrimitiveNumber(repr):
			strReprIdx = i
		}
	}

	if strReprIdx == -1 {
		return reprs
	}

	return append(reprs[:strReprIdx], reprs[strReprIdx+1:]...)
}

// MapAny maps the enum value to its representations.
func MapAny[Enum any](enum Enum, reprs []any) Enum {
	if mtmap.Get(mtkey.IsFinalized[Enum]()) {
		panic("enum is finalized")
	}

	var strRepr string
	var hasStrRepr bool
	var hasPrimitiveStr bool

	var numericRepr any
	var hasPrimitiveNumeric bool

	if xreflect.IsNumber(enum) {
		numericRepr = enum
		hasPrimitiveNumeric = true
	}

	if xreflect.IsString(enum) {
		strRepr = xreflect.Convert[string](enum)
		hasStrRepr = true
		hasPrimitiveStr = true
	}

	for _, repr := range reprs {
		switch {
		case xreflect.IsPrimitiveNumber(repr):
			if hasPrimitiveNumeric {
				panic(fmt.Sprintf("enum %s (%v): multiple primitive numerics are provided (%v, %v)",
					reflect.TypeOf(enum).Name(), enum, numericRepr, repr))
			}

			numericRepr = repr
			hasPrimitiveNumeric = true

		case xreflect.IsPrimitiveString(repr):
			if hasPrimitiveStr {
				panic(fmt.Sprintf("enum %s (%v): multiple primitive strings are provided (%v, %v)",
					reflect.TypeOf(enum).Name(), enum, strRepr, repr))
			}

			strRepr = xreflect.Convert[string](repr)
			hasStrRepr = true
			hasPrimitiveStr = true

		default:
			if v, ok := mtmap.Get2(mtkey.Extra2Enum[Enum](repr)); ok {
				panic(fmt.Sprintf("enum %s (%v): representation %v was already mapped to %v",
					reflect.TypeOf(enum).Name(), enum, repr, v))
			}

			if _, ok := mtmap.Get2(mtkey.Enum2ExtraWith(enum, repr)); ok {
				panic(fmt.Sprintf("enum %s (%v): do not map type %s twice",
					reflect.TypeOf(enum).Name(), enum, reflect.TypeOf(repr).Name()))
			}

			if !hasStrRepr {
				if xreflect.IsImplement[fmt.Stringer](repr) {
					strRepr = repr.(fmt.Stringer).String()
					hasStrRepr = true
				} else if xreflect.IsString(repr) {
					strRepr = xreflect.Convert[string](repr)
					hasStrRepr = true
				}
			}

			if numericRepr == nil && xreflect.IsNumber(repr) {
				numericRepr = repr
			}

			mtmap.Set(mtkey.Enum2ExtraWith(enum, repr), repr)
			mtmap.Set(mtkey.Extra2Enum[Enum](repr), enum)
		}
	}

	if !hasStrRepr {
		panic(fmt.Sprintf("enum %s (%v): not found any string representation", reflect.TypeOf(enum).Name(), enum))
	}

	if numericRepr == nil {
		numericRepr = GetAvailableEnumValue[Enum]()
	}

	mapEnumNumber(enum, numericRepr)

	if v, ok := mtmap.Get2(mtkey.String2Enum[Enum](strRepr)); ok {
		panic(fmt.Sprintf("enum %s (%v): string %s was already mapped to %v",
			reflect.TypeOf(enum).Name(), enum, strRepr, v))
	}

	if _, ok := mtmap.Get2(mtkey.Enum2String(enum)); ok {
		panic(fmt.Sprintf("enum %s (%v): do not map string twice", reflect.TypeOf(enum).Name(), enum))
	}

	mtmap.Set(mtkey.Enum2JSON(enum), strconv.Quote(strRepr))
	mtmap.Set(mtkey.Enum2String(enum), strRepr)
	mtmap.Set(mtkey.String2Enum[Enum](strRepr), enum)

	allVals := mtmap.Get(mtkey.AllEnums[Enum]())
	allVals = append(allVals, enum)
	mtmap.Set(mtkey.AllEnums[Enum](), allVals)

	return enum
}

// mapEnumNumber maps the enum to all its number representations (including
// signed and unsigned integers, floating-point numbers) and vice versa.
func mapEnumNumber[Enum any](enum Enum, n any) {
	if !xreflect.IsNumber(n) {
		panic(fmt.Sprintf("require a number, got %v", n))
	}

	if v, ok := mtmap.Get2(mtkey.AnyNumber2Enum[Enum](n)); ok {
		panic(fmt.Sprintf("enum %s (%v): number %v was already mapped to %v",
			reflect.TypeOf(enum).Name(), enum, n, v))
	}

	// The mapping to float32 always exists in all cases.
	if _, ok := mtmap.Get2(mtkey.Enum2Number[Enum, float32](enum)); ok {
		panic(fmt.Sprintf("enum %s (%v): do not map number twice", reflect.TypeOf(enum).Name(), enum))
	}

	// Only map the enum to integers if the enum is represented by integer
	// values, where the integer corresponds to the actual numeric value,
	// regardless of the underlying type.
	//
	// For example: float32(3) is also an integer, whereas float32(0.7) is not.
	mapInteger := true
	if xreflect.IsFloat32(n) {
		mapInteger = xreflect.Convert[float32](n) == xmath.Trunc32(xreflect.Convert[float32](n))
	} else if xreflect.IsFloat64(n) {
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
