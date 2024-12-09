package enum

import (
	"database/sql/driver"
	"fmt"

	"github.com/xybor-x/enum/internal/core"
)

// StructEnum provides a medium type-safe enum, which is better than IntEnum,
// but not as good as SafeEnum. Like IntEnum, it provides a set of built-in
// methods to simplify working with enums.
//
// The StructEnum enforces type safety by unexporting its fields, ensuring that
// only nearly predefined enum values are allowed. It prevents most accidental
// creation of new enum types, exception for zero initialization (note
// that the zero StructEnum is considered an invalid enum).
//
// Unlike SafeEnum, StructEnum supports both serialization and deserialization
// out of the box.
type StructEnum[dummyEnum any] struct {
	inner string
}

func NewStruct[dummyEnum any](s string) StructEnum[dummyEnum] {
	return core.MapAny(core.GetAvailableEnumValue[StructEnum[dummyEnum]](), StructEnum[dummyEnum]{inner: s}, s)
}

func (e StructEnum[dummyEnum]) IsValid() bool {
	return IsValid(e)
}

func (e StructEnum[dummyEnum]) MarshalJSON() ([]byte, error) {
	return MarshalJSON(e)
}

func (e *StructEnum[dummyEnum]) UnmarshalJSON(data []byte) error {
	return UnmarshalJSON(data, e)
}

func (e StructEnum[dummyEnum]) Value() (driver.Value, error) {
	return ValueSQL(e)
}

func (e *StructEnum[dummyEnum]) Scan(a any) error {
	return ScanSQL(a, e)
}

func (e StructEnum[dummyEnum]) Int() int {
	return ToInt(e)
}

func (e StructEnum[dummyEnum]) String() string {
	return ToString(e)
}

func (e StructEnum[dummyEnum]) GoString() string {
	if !IsValid(e) {
		return "<<undefined>>"
	}

	return fmt.Sprintf("%d (%s)", ToInt(e), e.inner)
}
