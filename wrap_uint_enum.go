package enum

import (
	"database/sql/driver"
	"encoding/xml"
	"fmt"

	"github.com/xybor-x/enum/internal/core"
	"github.com/xybor-x/enum/internal/xreflect"
	"gopkg.in/yaml.v3"
)

var _ newableEnum = WrapUintEnum[int](0)
var _ hookAfterEnum = WrapUintEnum[int](0)

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

func (e WrapUintEnum[underlyingEnum]) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return MarshalXML(encoder, start, e)
}

func (e *WrapUintEnum[underlyingEnum]) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return UnmarshalXML(decoder, start, e)
}

func (e WrapUintEnum[underlyingEnum]) MarshalYAML() (any, error) {
	return MarshalYAML(e)
}

func (e *WrapUintEnum[underlyingEnum]) UnmarshalYAML(node *yaml.Node) error {
	return UnmarshalYAML(node, e)
}

func (e WrapUintEnum[underlyingEnum]) Value() (driver.Value, error) {
	return ValueSQL(e)
}

func (e *WrapUintEnum[underlyingEnum]) Scan(a any) error {
	return ScanSQL(a, e)
}

// To returns the underlying representation of this enum.
func (e WrapUintEnum[underlyingEnum]) To() underlyingEnum {
	return MustTo[underlyingEnum](e)
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
func (e WrapUintEnum[underlyingEnum]) newEnum(repr []any) any {
	numeric := core.GetNumericRepresentation(repr)
	if numeric == nil {
		numeric = core.GetAvailableEnumValue[WrapUintEnum[underlyingEnum]]()
	} else {
		repr = core.RemoveNumericRepresentation(repr)
	}

	return core.MapAny(xreflect.Convert[WrapUintEnum[underlyingEnum]](numeric), repr)
}

// WARNING: Only use this function if you fully understand its behavior.
// It might cause unexpected results if used improperly.
func (e WrapUintEnum[underlyingEnum]) hookAfter() {
	mustHaveUnderlyingRepr[underlyingEnum](e)
}
