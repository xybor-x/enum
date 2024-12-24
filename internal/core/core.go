package core

import (
	"fmt"
	"math"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/xybor-x/enum/internal/mtkey"
	"github.com/xybor-x/enum/internal/mtmap"
	"github.com/xybor-x/enum/internal/xmath"
	"github.com/xybor-x/enum/internal/xreflect"
)

func GetAvailableEnumValue[Enum any]() int64 {
	id := int64(0)
	for {
		if _, ok := mtmap.Get2(mtkey.Repr2Enum[Enum](id)); !ok {
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
		panic(fmt.Sprintf("enum %s: the enum was already finalized", TrueNameOf[Enum]()))
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
				panic(fmt.Sprintf("enum %s (%#v): multiple primitive numerics are provided (%v, %v)",
					TrueNameOf[Enum](), enum, numericRepr, repr))
			}

			numericRepr = repr
			hasPrimitiveNumeric = true

		case xreflect.IsPrimitiveString(repr):
			if hasPrimitiveStr {
				panic(fmt.Sprintf("enum %s (%#v): multiple primitive strings are provided (%v, %v)",
					TrueNameOf[Enum](), enum, strRepr, repr))
			}

			strRepr = xreflect.Convert[string](repr)
			hasStrRepr = true
			hasPrimitiveStr = true

		default:
			if v, ok := mtmap.Get2(mtkey.Repr2Enum[Enum](repr)); ok {
				panic(fmt.Sprintf("enum %s (%#v): representation %v of %T was already mapped to %v",
					TrueNameOf[Enum](), enum, repr, repr, v))
			}

			if _, ok := mtmap.Get2(mtkey.Enum2ReprWith(enum, repr)); ok {
				panic(fmt.Sprintf("enum %s (%#v): do not map type %s twice",
					TrueNameOf[Enum](), enum, reflect.TypeOf(repr).Name()))
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

			mtmap.Set(mtkey.Enum2ReprWith(enum, repr), repr)
			mtmap.Set(mtkey.Repr2Enum[Enum](repr), enum)
		}
	}

	if !hasStrRepr {
		panic(fmt.Sprintf("enum %s (%#v): not found any string representation", TrueNameOf[Enum](), enum))
	}

	if numericRepr == nil {
		numericRepr = GetAvailableEnumValue[Enum]()
	}

	mapEnumNumber(enum, numericRepr)

	if v, ok := mtmap.Get2(mtkey.Repr2Enum[Enum](strRepr)); ok {
		panic(fmt.Sprintf("enum %s (%#v): string %s was already mapped to %v",
			TrueNameOf[Enum](), enum, strRepr, v))
	}

	if _, ok := mtmap.Get2(mtkey.Enum2Repr[Enum, string](enum)); ok {
		panic(fmt.Sprintf("enum %s (%#v): do not map string twice", TrueNameOf[Enum](), enum))
	}

	mtmap.Set(mtkey.Enum2JSON(enum), strconv.Quote(strRepr))
	mtmap.Set(mtkey.Enum2Repr[Enum, string](enum), any(strRepr))
	mtmap.Set(mtkey.Repr2Enum[Enum](strRepr), enum)

	allVals := mtmap.Get(mtkey.AllEnums[Enum]())
	allVals = append(allVals, enum)
	mtmap.Set(mtkey.AllEnums[Enum](), allVals)

	return enum
}

var advancedEnumNames = []string{"WrapEnum", "WrapUintEnum", "WrapFloatEnum", "SafeEnum"}

func NameOf[T any]() string {
	if name, ok := mtmap.Get2(mtkey.NameOf[T]()); ok {
		return name
	}

	name := reflect.TypeOf((*T)(nil)).Elem().Name()
	for _, prefix := range advancedEnumNames {
		if strings.HasPrefix(name, prefix) {
			name = capitalizeFirst(getUnderlyingName(name, prefix))
			break
		}
	}

	mtmap.Set(mtkey.NameOf[T](), name)
	return name
}

func TrueNameOf[T any]() string {
	if name, ok := mtmap.Get2(mtkey.TrueNameOf[T]()); ok {
		return name
	}

	name := reflect.TypeOf((*T)(nil)).Elem().Name()
	for _, prefix := range advancedEnumNames {
		if strings.HasPrefix(name, prefix) {
			name = fmt.Sprintf("%s[%s]", prefix, getUnderlyingName(name, prefix))
			break
		}
	}

	mtmap.Set(mtkey.TrueNameOf[T](), name)
	return name
}

func getUnderlyingName(name, prefix string) string {
	// name = prefix[path/to/module.underlying·id]
	inner := name[len(prefix)+1 : len(name)-1] // inner = path/to/module.underlying·id
	_, inner = path.Split(inner)               // inner = module.underlying·id

	parts := strings.Split(inner, ".")
	inner = parts[len(parts)-1] // inner = underlying·id

	parts = strings.Split(inner, string(rune(183))) // middle dot character.
	return parts[0]                                 // parts[0] = underlying
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s // Return empty string if input is empty
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

// mapEnumNumber maps the enum to all its number representations (including
// signed and unsigned integers, floating-point numbers) and vice versa.
func mapEnumNumber[Enum any](enum Enum, n any) {
	if !xreflect.IsNumber(n) {
		panic(fmt.Sprintf("require a number, got %v", n))
	}

	if v, ok := mtmap.Get2(mtkey.Repr2Enum[Enum](n)); ok {
		panic(fmt.Sprintf("enum %s (%v): number %v was already mapped to %v",
			reflect.TypeOf(enum).Name(), enum, n, v))
	}

	// The mapping to float32 always exists in all cases.
	if _, ok := mtmap.Get2(mtkey.Enum2Repr[Enum, float32](enum)); ok {
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
		mtmap.Set(mtkey.Enum2Repr[Enum, int](enum), any(xreflect.Convert[int](n)))
		mtmap.Set(mtkey.Enum2Repr[Enum, int8](enum), any(xreflect.Convert[int8](n)))
		mtmap.Set(mtkey.Enum2Repr[Enum, int16](enum), any(xreflect.Convert[int16](n)))
		mtmap.Set(mtkey.Enum2Repr[Enum, int32](enum), any(xreflect.Convert[int32](n)))
		mtmap.Set(mtkey.Enum2Repr[Enum, int64](enum), any(xreflect.Convert[int64](n)))

		// Map enum to all unsigned integers.
		mtmap.Set(mtkey.Enum2Repr[Enum, uint](enum), any(xreflect.Convert[uint](n)))
		mtmap.Set(mtkey.Enum2Repr[Enum, uint8](enum), any(xreflect.Convert[uint8](n)))
		mtmap.Set(mtkey.Enum2Repr[Enum, uint16](enum), any(xreflect.Convert[uint16](n)))
		mtmap.Set(mtkey.Enum2Repr[Enum, uint32](enum), any(xreflect.Convert[uint32](n)))
		mtmap.Set(mtkey.Enum2Repr[Enum, uint64](enum), any(xreflect.Convert[uint64](n)))

		// Map all signed integers to enum.
		mtmap.Set(mtkey.Repr2Enum[Enum](xreflect.Convert[int](n)), enum)
		mtmap.Set(mtkey.Repr2Enum[Enum](xreflect.Convert[int8](n)), enum)
		mtmap.Set(mtkey.Repr2Enum[Enum](xreflect.Convert[int16](n)), enum)
		mtmap.Set(mtkey.Repr2Enum[Enum](xreflect.Convert[int32](n)), enum)
		mtmap.Set(mtkey.Repr2Enum[Enum](xreflect.Convert[int64](n)), enum)

		// Map all unsigned integers to enum.
		mtmap.Set(mtkey.Repr2Enum[Enum](xreflect.Convert[uint](n)), enum)
		mtmap.Set(mtkey.Repr2Enum[Enum](xreflect.Convert[uint8](n)), enum)
		mtmap.Set(mtkey.Repr2Enum[Enum](xreflect.Convert[uint16](n)), enum)
		mtmap.Set(mtkey.Repr2Enum[Enum](xreflect.Convert[uint32](n)), enum)
		mtmap.Set(mtkey.Repr2Enum[Enum](xreflect.Convert[uint64](n)), enum)
	}

	// Map enum to all floats.
	mtmap.Set(mtkey.Enum2Repr[Enum, float32](enum), any(xreflect.Convert[float32](n)))
	mtmap.Set(mtkey.Enum2Repr[Enum, float64](enum), any(xreflect.Convert[float64](n)))

	// Map all floats to enum.
	mtmap.Set(mtkey.Repr2Enum[Enum](xreflect.Convert[float32](n)), enum)
	mtmap.Set(mtkey.Repr2Enum[Enum](xreflect.Convert[float64](n)), enum)
}
