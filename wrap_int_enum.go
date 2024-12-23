package enum

import (
	"database/sql/driver"
	"fmt"

	"github.com/xybor-x/enum/internal/core"
	"github.com/xybor-x/enum/internal/mtkey"
	"github.com/xybor-x/enum/internal/mtmap"
	"github.com/xybor-x/enum/internal/xreflect"
)

var _ newEnumable = (WrapEnum[int](0))

// WrapEnum provides a set of built-in methods to simplify working with int
// enums.
type WrapEnum[underlyingEnum any] int

func (e WrapEnum[underlyingEnum]) IsValid() bool {
	return IsValid(e)
}

func (e WrapEnum[underlyingEnum]) MarshalJSON() ([]byte, error) {
	return MarshalJSON(e)
}

func (e *WrapEnum[underlyingEnum]) UnmarshalJSON(data []byte) error {
	return UnmarshalJSON(data, e)
}

func (e WrapEnum[underlyingEnum]) Value() (driver.Value, error) {
	return ValueSQL(e)
}

func (e *WrapEnum[underlyingEnum]) Scan(a any) error {
	return ScanSQL(a, e)
}

// Int returns the int representation of the enum. This method returns the value
// of math.MinInt32 if the enum is invalid.
//
// DEPRECATED: directly cast the enum to int instead.
func (e WrapEnum[underlyingEnum]) Int() int {
	return mtmap.Get(mtkey.Enum2Number[WrapEnum[underlyingEnum], int](e))
}

// To returns the underlying representation of this enum.
func (e WrapEnum[underlyingEnum]) To() underlyingEnum {
	return To[underlyingEnum](e)
}

func (e WrapEnum[underlyingEnum]) String() string {
	return ToString(e)
}

func (e WrapEnum[underlyingEnum]) GoString() string {
	if !e.IsValid() {
		return fmt.Sprintf("%d", e)
	}

	return fmt.Sprintf("%d (%s)", e, e)
}

// WARNING: Only use this function if you fully understand its behavior.
// It might cause unexpected results if used improperly.
func (e WrapEnum[underlyingEnum]) newEnum(repr []any) any {
	numeric := core.GetNumericRepresentation(repr)
	if numeric == nil {
		numeric = core.GetAvailableEnumValue[WrapEnum[underlyingEnum]]()
	}

	return core.MapAny(xreflect.Convert[WrapEnum[underlyingEnum]](numeric), repr)
}
