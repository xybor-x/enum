package enum

import (
	"database/sql/driver"
	"fmt"
	"strconv"

	"github.com/xybor-x/enum/internal/core"
)

// WrapEnum provides a set of built-in methods to simplify working with enums.
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

func (e WrapEnum[underlyingEnum]) Int() int {
	return ToInt(e)
}

func (e WrapEnum[underlyingEnum]) String() string {
	return ToString(e)
}

func (e WrapEnum[underlyingEnum]) GoString() string {
	if !e.IsValid() {
		return strconv.Itoa(int(e))
	}

	return fmt.Sprintf("%d (%s)", e, e)
}

// WARNING: Only use this function if you fully understand its behavior.
// It might cause unexpected results if used improperly.
func (e WrapEnum[underlyingEnum]) newInnerEnum(s string) any {
	id := core.GetAvailableEnumValue[WrapEnum[underlyingEnum]]()
	return core.MapAny(id, WrapEnum[underlyingEnum](id), s)
}
