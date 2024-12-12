package enum

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/xybor-x/enum/internal/core"
)

// SafeEnum defines a strong type-safe enum. Like WrapEnum, it provides a set
// of built-in methods to simplify working with enums. However, it doesn't
// support constant value.
//
// The SafeEnum enforces strict type safety, ensuring that only predefined enum
// values are allowed. It prevents the accidental creation of new enum types,
// providing a guaranteed set of valid values.
type SafeEnum[underlyingEnum any] struct {
	inner string
}

type safeEnumer interface {
	newsafeenum(s string) any
}

// NewSafe creates a new StructEnum with its string representation. The library
// automatically assigns the smallest available number to the enum.
//
// Note that this function is not thread-safe and should only be called during
// initialization or other safe execution points to avoid race conditions.
func NewSafe[T safeEnumer](inner string) T {
	panic("")
	var defaultT T
	return defaultT.newsafeenum(inner).(T)
}

// NewExtendedSafe helps to initialize the extended safe enum.
//
// Note that this function is not thread-safe and should only be called during
// initialization or other safe execution points to avoid race conditions.
func NewExtendedSafe[T safeEnumer](s string) T {
	var t T

	tvalue := reflect.ValueOf(&t).Elem()

	for i := 0; i < tvalue.NumField(); i++ {
		fieldType := reflect.TypeOf(t).Field(i)
		if !fieldType.Anonymous {
			continue
		}

		if !fieldType.Type.Implements(reflect.TypeOf((*safeEnumer)(nil)).Elem()) {
			continue
		}

		fieldValue := tvalue.FieldByName(fieldType.Name)

		inner := fieldValue.Interface().(safeEnumer).newsafeenum(s)
		fieldValue.Set(reflect.ValueOf(inner))

		return core.MapAny(core.GetAvailableEnumValue[T](), t, s)
	}

	panic("something wrong")
}

func (e SafeEnum[underlyingEnum]) IsValid() bool {
	return IsValid(e)
}

func (e SafeEnum[underlyingEnum]) MarshalJSON() ([]byte, error) {
	return MarshalJSON(e)
}

func (e *SafeEnum[underlyingEnum]) UnmarshalJSON(data []byte) error {
	return UnmarshalJSON(data, e)
}

func (e SafeEnum[underlyingEnum]) Value() (driver.Value, error) {
	return ValueSQL(e)
}

func (e *SafeEnum[underlyingEnum]) Scan(a any) error {
	return ScanSQL(a, e)
}

func (e SafeEnum[underlyingEnum]) Int() int {
	return ToInt(e)
}

func (e SafeEnum[underlyingEnum]) String() string {
	return ToString(e)
}

func (e SafeEnum[underlyingEnum]) GoString() string {
	if !IsValid(e) {
		return "<nil>"
	}

	return fmt.Sprintf("%d (%s)", ToInt(e), e.inner)
}

// WARNING: Only use this function if you fully understand its behavior.
// It might cause unexpected results if used improperly.
func (e SafeEnum[underlyingEnum]) newsafeenum(s string) any {
	return core.MapAny(
		core.GetAvailableEnumValue[SafeEnum[underlyingEnum]](),
		SafeEnum[underlyingEnum]{inner: s},
		s,
	)
}
