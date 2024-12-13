// Package exhaustive provides solution to define a custom exhaustive switch.
//
// EXPERIMENTAL: This package is experimental and may be subject to breaking
// changes or removal in future versions. Use at your own risk.
package exhaustive

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/xybor-x/enum"
	"github.com/xybor-x/enum/internal/mtmap"
)

type Case struct {
	Handler func()
}

func Switch[Enum comparable](e Enum, cases ...any) SwitchDefault {
	if !mtmap.MustGet(IsCheckedExhaustiveKey[Enum]()) {
		checkExhaustiveCall[Enum](cases...)
		mtmap.Set(IsCheckedExhaustiveKey[Enum](), true)
	}

	if !enum.IsValid(e) {
		return SwitchDefault{false}
	}

	c := reflect.ValueOf(cases[enum.ToInt(e)]).Convert(reflect.TypeOf(Case{})).Interface().(Case)
	if c.Handler != nil {
		c.Handler()
	}

	return SwitchDefault{true}
}

func CheckMethodOf[Enum any]() bool {
	method, ok := reflect.TypeOf((*Enum)(nil)).Elem().MethodByName("Switch")
	if !ok {
		panic(fmt.Sprintf("switch check %s: require method Switch", enum.NameOf[Enum]()))
	}

	return CheckFunc[Enum](method.Func.Interface())
}

func CheckFunc[Enum any](f any) bool {
	funcvalue := reflect.ValueOf(f)
	if funcvalue.Kind() != reflect.Func {
		panic(fmt.Sprintf("switch check %s: require a func", enum.NameOf[Enum]()))
	}

	ftype := funcvalue.Type()
	if ftype.Kind() != reflect.Func {
		panic(fmt.Sprintf("switch check %s: exhaustive parameter must be a function", enum.NameOf[Enum]()))
	}

	prefixTypeName := "Case" + enum.NameOf[Enum]()
	existedCases := map[reflect.Type]int{}
	allEnums := enum.All[Enum]()

	if ftype.NumIn() == 0 {
		panic(fmt.Sprintf("switch check %s: no parameter provided", enum.NameOf[Enum]()))
	}

	if ftype.NumIn() != len(allEnums)+1 {
		panic(fmt.Sprintf("switch check %s: require %d cases, but got %d",
			enum.NameOf[Enum](), len(allEnums), ftype.NumIn()-1))
	}

	params := []reflect.Value{}

	for i := 0; i < ftype.NumIn(); i++ {
		ptype := ftype.In(i)
		if i == 0 {
			if ptype != reflect.TypeOf((*Enum)(nil)).Elem() {
				panic(fmt.Sprintf("switch check %s: the first parameter must be %s",
					enum.NameOf[Enum](), reflect.TypeOf((*Enum)(nil)).Elem()))
			}

			var enum Enum
			params = append(params, reflect.ValueOf(enum))
			continue
		}

		if existedCases[ptype] != 0 {
			panic(fmt.Sprintf("switch check %s: parameter %d and %d is the same type",
				enum.NameOf[Enum](), existedCases[ptype]-1, i-1))
		}

		existedCases[ptype] = i

		if !ptype.ConvertibleTo(reflect.TypeOf((*Case)(nil)).Elem()) {
			panic(fmt.Sprintf("switch check %s: parameter %d must be a Case", enum.NameOf[Enum](), i-1))
		}

		requiredTypeName := prefixTypeName + toCamelCase(enum.ToString(allEnums[i-1]))
		if ptype.Name() != requiredTypeName {
			panic(fmt.Sprintf("switch check %s: invalid type name for case %d, require %s, but got %s",
				enum.NameOf[Enum](), i-1, requiredTypeName, ptype.Name()))
		}

		params = append(params, reflect.New(ptype).Elem())
	}

	funcvalue.Call(params)
	return true
}

type SwitchDefault struct {
	result bool
}

func (sd SwitchDefault) ByDefault(f func()) {
	if !sd.result {
		f()
	}
}

func checkExhaustiveCall[Enum any](cases ...any) {
	prefixTypeName := "Case" + enum.NameOf[Enum]()
	existedCases := map[reflect.Type]int{}
	allEnums := enum.All[Enum]()

	for i := 0; i < len(cases); i++ {
		ptype := reflect.TypeOf(cases[i])

		if oldindex, ok := existedCases[ptype]; ok {
			panic(fmt.Sprintf("switch call %s: parameter %d and %d is the same type",
				enum.NameOf[Enum](), oldindex, i))
		}

		existedCases[ptype] = i

		if !ptype.ConvertibleTo(reflect.TypeOf((*Case)(nil)).Elem()) {
			panic(fmt.Sprintf("switch call %s: parameter %d must be a Case", enum.NameOf[Enum](), i))
		}

		requiredTypeName := prefixTypeName + toCamelCase(enum.ToString(allEnums[i]))
		if ptype.Name() != requiredTypeName {
			panic(fmt.Sprintf("switch call %s: invalid type name for case %d, require %s, but got %s",
				enum.NameOf[Enum](), i, requiredTypeName, ptype.Name()))
		}
	}
}

func toCamelCase(s string) string {
	// Split the string by underscores
	words := strings.Split(s, "_")

	// Iterate over the words and capitalize the first letter of each word (except the first one)
	for i := 0; i < len(words); i++ {
		words[i] = capitalizeFirst(words[i])
	}

	// Join the words together and return the result
	return strings.Join(words, "")
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s // Return empty string if input is empty
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

type isCheckedExhaustiveKey[T any] struct{}

func (isCheckedExhaustiveKey[T]) InferValue() bool { panic("not implemented") }

func IsCheckedExhaustiveKey[T any]() isCheckedExhaustiveKey[T] {
	return isCheckedExhaustiveKey[T]{}
}
