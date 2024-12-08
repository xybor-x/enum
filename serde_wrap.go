package enum

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/xybor-x/enum/internal/common"
)

// Serde provides functionality for serializing and deserializing enums
// that cannot be directly serialized or deserialized.
//
// Note: This struct is intentionally non-comparable. If you need a comparable
// version, use ComparableSerde instead.
//
// Why is it non-comparable?
// The primary purpose of this struct is to handle serialization and
// deserialization. It is not intended to replace the role or purpose of the
// enum itself.
type Serde[Enum any] struct {
	ComparableSerde[Enum]
	_ []byte // prevent comparison
}

func SerdeWrap[Enum any](enum Enum) Serde[Enum] {
	return Serde[Enum]{ComparableSerde: ComparableSerdeWrap(enum)}
}

// ComparableSerde facilitates the serialization and deserialization of enums
// that cannot be directly serialized or deserialized.
//
// It is similar to Store but comparable.
type ComparableSerde[Enum any] struct {
	enum Enum
}

func ComparableSerdeWrap[Enum any](enum Enum) ComparableSerde[Enum] {
	return ComparableSerde[Enum]{enum: enum}
}

// Enum returns the inner enum.
func (e ComparableSerde[Enum]) Enum() Enum {
	return e.enum
}

func (e ComparableSerde[Enum]) MarshalJSON() ([]byte, error) {
	return MarshalJSON(e.enum)
}

func (e *ComparableSerde[Enum]) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	enum, ok := FromString[Enum](str)
	if !ok {
		return fmt.Errorf("enum %s: invalid string %s", common.NameOf[Enum](), str)
	}

	e.enum = enum
	return nil
}

func (e ComparableSerde[Enum]) Value() (driver.Value, error) {
	return ValueSQL(e.enum)
}

func (e *ComparableSerde[Enum]) Scan(a any) error {
	var data string
	switch t := a.(type) {
	case string:
		data = t
	case []byte:
		data = string(t)
	default:
		return fmt.Errorf("not support type %T", a)
	}

	enum, ok := FromString[Enum](data)
	if !ok {
		return fmt.Errorf("enum %s: invalid string %s", common.NameOf[Enum](), a)
	}

	e.enum = enum
	return nil
}
