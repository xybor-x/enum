package enum

import (
	"database/sql/driver"
	"encoding/xml"
	"fmt"

	"github.com/xybor-x/enum/internal/core"
	"gopkg.in/yaml.v3"
)

var _ newableEnum = SafeEnum[int]{}
var _ hookAfterEnum = SafeEnum[int]{}

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

func (e SafeEnum[underlyingEnum]) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return MarshalXML(encoder, start, e)
}

func (e *SafeEnum[underlyingEnum]) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return UnmarshalXML(decoder, start, e)
}

func (e SafeEnum[underlyingEnum]) MarshalYAML() (any, error) {
	return MarshalYAML(e)
}

func (e *SafeEnum[underlyingEnum]) UnmarshalYAML(node *yaml.Node) error {
	return UnmarshalYAML(node, e)
}

func (e SafeEnum[underlyingEnum]) Value() (driver.Value, error) {
	return ValueSQL(e)
}

func (e *SafeEnum[underlyingEnum]) Scan(a any) error {
	return ScanSQL(a, e)
}

func (e SafeEnum[underlyingEnum]) Int() int {
	return MustTo[int](e)
}

// To returns the underlying representation of this enum.
func (e SafeEnum[underlyingEnum]) To() underlyingEnum {
	return MustTo[underlyingEnum](e)
}

func (e SafeEnum[underlyingEnum]) String() string {
	return ToString(e)
}

func (e SafeEnum[underlyingEnum]) GoString() string {
	if !IsValid(e) {
		return "<nil>"
	}

	return fmt.Sprintf("%d (%s)", e.Int(), e.inner)
}

// WARNING: Only use this function if you fully understand its behavior.
// It might cause unexpected results if used improperly.
func (e SafeEnum[underlyingEnum]) newEnum(reprs []any) any {
	str, ok := core.GetStringRepresentation(reprs)
	if !ok {
		panic("SafeEnum requires at least a string representation")
	}

	return core.MapAny(SafeEnum[underlyingEnum]{inner: str}, reprs)
}

// WARNING: Only use this function if you fully understand its behavior.
// It might cause unexpected results if used improperly.
func (e SafeEnum[underlyingEnum]) hookAfter() {
	mustHaveUnderlyingRepr[underlyingEnum](e)
}
