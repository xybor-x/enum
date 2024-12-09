package safeenum

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xybor-x/enum"
	"github.com/xybor-x/enum/internal/core"
)

// SafeEnum defines a strong type-safe enum. Like IntEnum, it provides a set of
// built-in methods to simplify working with enums.
//
// The SafeEnum enforces strict type safety, ensuring that only predefined enum
// values are allowed. It prevents the accidental creation of new enum types,
// providing a guaranteed set of valid values.
//
// Note: This interface does not include deserialization capabilities. If you
// require serialization and deserialization functionality, consider using
// enum.Serde, which provides additional support for those operations.
type SafeEnum[unsafeEnum any] interface {
	fmt.Stringer
	fmt.GoStringer

	json.Marshaler
	driver.Valuer

	Int() int

	safeenum() unsafeEnum
}

func New[unsafeEnum any, P positioner](s string) SafeEnum[unsafeEnum] {
	var p P
	return core.MapAny[SafeEnum[unsafeEnum]](int64(p.position()), safeEnum[unsafeEnum, P]{}, s)
}

type safeEnum[unsafeEnum any, P positioner] struct{}

func (e safeEnum[unsafeEnum, P]) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSON[SafeEnum[unsafeEnum]](e)
}

// DO NOT USE
func (e *safeEnum[unsafeEnum, P]) UnmarshalJSON([]byte) error {
	return errors.New("not implemented")
}

func (e safeEnum[unsafeEnum, P]) Value() (driver.Value, error) {
	return enum.ValueSQL[SafeEnum[unsafeEnum]](e)
}

// DO NOT USE
func (e *safeEnum[unsafeEnum, P]) Scan(data any) error {
	return errors.New("not implemented")
}

func (e safeEnum[unsafeEnum, P]) Int() int {
	return enum.ToInt[SafeEnum[unsafeEnum]](e)
}

func (e safeEnum[unsafeEnum, P]) String() string {
	return enum.ToString[SafeEnum[unsafeEnum]](e)
}

func (e safeEnum[unsafeEnum, P]) GoString() string {
	return fmt.Sprintf("%d (%s)", e.Int(), e)
}

func (e safeEnum[unsafeEnum, P]) safeenum() unsafeEnum { panic("not implemented") }
