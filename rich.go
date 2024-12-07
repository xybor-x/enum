package enum

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

var (
	_ json.Marshaler   = RichEnum[int](0)
	_ json.Unmarshaler = (*RichEnum[int])(nil)
	_ driver.Valuer    = RichEnum[int](0)
	_ sql.Scanner      = (*RichEnum[int])(nil)
	_ fmt.Stringer     = RichEnum[int](0)
)

// RichEnum provides a set of utility methods to simplify working with enums.
//
// It includes various helper functions for operations like serialization,
// deserialization, string conversion, and validation, making it easier to
// manage and manipulate enum values across your codebase.
type RichEnum[T any] int

func (e RichEnum[T]) IsValid() bool {
	return IsValid(e)
}

func (e RichEnum[T]) MarshalJSON() ([]byte, error) {
	return MarshalJSON(e)
}

func (e *RichEnum[T]) UnmarshalJSON(data []byte) error {
	return UnmarshalJSON(data, e)
}

func (e RichEnum[T]) Value() (driver.Value, error) {
	return ValueSQL(e)
}

func (e *RichEnum[T]) Scan(a any) error {
	return ScanSQL(a, e)
}

func (e RichEnum[T]) Int() int {
	return int(e)
}

func (e RichEnum[T]) Repr() string {
	return fmt.Sprintf("%d (%s)", e.Int(), e.String())
}

func (e RichEnum[T]) String() string {
	return StringOf(e)
}
