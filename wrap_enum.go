package enum

import (
	"database/sql/driver"
	"fmt"
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
		return "<nil>"
	}

	return fmt.Sprintf("%d (%s)", e, e)
}
