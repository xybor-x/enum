package enum

import (
	"database/sql/driver"
	"fmt"

	"github.com/xybor-x/enum/internal/core"
)

// WrapUintEnum provides a set of built-in methods to simplify working with uint
// enums.
type WrapUintEnum[underlyingEnum any] uint

func (e WrapUintEnum[underlyingEnum]) IsValid() bool {
	return IsValid(e)
}

func (e WrapUintEnum[underlyingEnum]) MarshalJSON() ([]byte, error) {
	return MarshalJSON(e)
}

func (e *WrapUintEnum[underlyingEnum]) UnmarshalJSON(data []byte) error {
	return UnmarshalJSON(data, e)
}

func (e WrapUintEnum[underlyingEnum]) Value() (driver.Value, error) {
	return ValueSQL(e)
}

func (e *WrapUintEnum[underlyingEnum]) Scan(a any) error {
	return ScanSQL(a, e)
}

func (e WrapUintEnum[underlyingEnum]) String() string {
	return ToString(e)
}

func (e WrapUintEnum[underlyingEnum]) GoString() string {
	if !e.IsValid() {
		return fmt.Sprintf("%d", e)
	}

	return fmt.Sprintf("%d (%s)", e, e)
}

// WARNING: Only use this function if you fully understand its behavior.
// It might cause unexpected results if used improperly.
func (e WrapUintEnum[underlyingEnum]) newEnum(id int64, s string) any {
	return core.MapAny(id, WrapUintEnum[underlyingEnum](id), s)
}
