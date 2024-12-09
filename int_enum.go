package enum

import (
	"database/sql/driver"
	"fmt"
)

// IntEnum provides a set of built-in methods to simplify working with enums.
type IntEnum[dummyEnum any] int

func (e IntEnum[dummyEnum]) IsValid() bool {
	return IsValid(e)
}

func (e IntEnum[dummyEnum]) MarshalJSON() ([]byte, error) {
	return MarshalJSON(e)
}

func (e *IntEnum[dummyEnum]) UnmarshalJSON(data []byte) error {
	return UnmarshalJSON(data, e)
}

func (e IntEnum[dummyEnum]) Value() (driver.Value, error) {
	return ValueSQL(e)
}

func (e *IntEnum[dummyEnum]) Scan(a any) error {
	return ScanSQL(a, e)
}

func (e IntEnum[dummyEnum]) Int() int {
	return ToInt(e)
}

func (e IntEnum[dummyEnum]) String() string {
	return ToString(e)
}

func (e IntEnum[dummyEnum]) GoString() string {
	if !e.IsValid() {
		return fmt.Sprintf("%d (<<undefined>>)", e)
	}

	return fmt.Sprintf("%d (%s)", e, e)
}
