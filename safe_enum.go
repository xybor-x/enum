package enum

import (
	"database/sql/driver"
	"fmt"

	"github.com/xybor-x/enum/internal/core"
	"github.com/xybor-x/enum/internal/mtkey"
	"github.com/xybor-x/enum/internal/mtmap"
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
	return mtmap.MustGet(mtkey.Enum2Number[SafeEnum[underlyingEnum], int](e))
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
func (e SafeEnum[underlyingEnum]) newEnum(id int64, s string) any {
	return core.MapAny(id, SafeEnum[underlyingEnum]{inner: s}, s)
}
