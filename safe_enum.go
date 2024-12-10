package enum

import (
	"database/sql/driver"
	"fmt"

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

// NewSafe creates a new StructEnum with its string representation. The library
// automatically assigns the smallest available number to the enum.
//
// Note that this function is not thread-safe and should only be called during
// initialization or other safe execution points to avoid race conditions.
func NewSafe[underlyingEnum any](inner string) SafeEnum[underlyingEnum] {
	return core.MapAny(
		core.GetAvailableEnumValue[SafeEnum[underlyingEnum]](),
		SafeEnum[underlyingEnum]{inner: inner},
		inner,
	)
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
