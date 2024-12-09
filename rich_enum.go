package enum

import (
	"database/sql/driver"
	"fmt"
)

// RichEnum provides a set of built-in methods to simplify working with enums.
type RichEnum[dummyEnum any] int

func (e RichEnum[dummyEnum]) IsValid() bool {
	return IsValid(e)
}

func (e RichEnum[dummyEnum]) MarshalJSON() ([]byte, error) {
	return MarshalJSON(e)
}

func (e *RichEnum[dummyEnum]) UnmarshalJSON(data []byte) error {
	return UnmarshalJSON(data, e)
}

func (e RichEnum[dummyEnum]) Value() (driver.Value, error) {
	return ValueSQL(e)
}

func (e *RichEnum[dummyEnum]) Scan(a any) error {
	return ScanSQL(a, e)
}

func (e RichEnum[dummyEnum]) Int() int {
	return ToInt(e)
}

func (e RichEnum[dummyEnum]) String() string {
	return ToString(e)
}

func (e RichEnum[dummyEnum]) GoString() string {
	if !e.IsValid() {
		return fmt.Sprintf("%d (<<undefined>>)", e)
	}

	return fmt.Sprintf("%d (%s)", e, e)
}
