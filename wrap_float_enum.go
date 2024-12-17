package enum

import (
	"database/sql/driver"
	"fmt"
	"strconv"

	"github.com/xybor-x/enum/internal/core"
)

// WrapFloatEnum provides a set of built-in methods to simplify working with
// float64 enums.
type WrapFloatEnum[underlyingEnum any] float64

func (e WrapFloatEnum[underlyingEnum]) IsValid() bool {
	return IsValid(e)
}

func (e WrapFloatEnum[underlyingEnum]) MarshalJSON() ([]byte, error) {
	return MarshalJSON(e)
}

func (e *WrapFloatEnum[underlyingEnum]) UnmarshalJSON(data []byte) error {
	return UnmarshalJSON(data, e)
}

func (e WrapFloatEnum[underlyingEnum]) Value() (driver.Value, error) {
	return ValueSQL(e)
}

func (e *WrapFloatEnum[underlyingEnum]) Scan(a any) error {
	return ScanSQL(a, e)
}

func (e WrapFloatEnum[underlyingEnum]) String() string {
	return ToString(e)
}

func (e WrapFloatEnum[underlyingEnum]) GoString() string {
	if !e.IsValid() {
		return strconv.Itoa(int(e))
	}

	return fmt.Sprintf("%f (%s)", e, e)
}

// WARNING: Only use this function if you fully understand its behavior.
// It might cause unexpected results if used improperly.
func (e WrapFloatEnum[underlyingEnum]) newEnum(id int64, s string) any {
	return core.MapAny(id, WrapFloatEnum[underlyingEnum](id), s)
}
