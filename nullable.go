package enum

import (
	"database/sql/driver"

	"gopkg.in/yaml.v3"
)

// Nullable allows handling nullable enums in JSON, YAML, and SQL.
type Nullable[Enum any] struct {
	Enum  Enum
	Valid bool
}

func (e Nullable[Enum]) MarshalJSON() ([]byte, error) {
	if !e.Valid {
		return []byte("null"), nil
	}

	return MarshalJSON(e.Enum)
}

func (e *Nullable[Enum]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		var defaultEnum Enum
		e.Enum, e.Valid = defaultEnum, false
		return nil
	}

	return UnmarshalJSON(data, &e.Enum)
}

func (e Nullable[Enum]) MarshalYAML() (any, error) {
	if !e.Valid {
		return yaml.Node{
			Kind: yaml.ScalarNode,
			Tag:  "!!null", // Use the YAML null tag
		}, nil
	}

	return MarshalYAML(e.Enum)
}

func (e *Nullable[Enum]) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode && node.Tag == "!!null" {
		var defaultEnum Enum
		e.Enum, e.Valid = defaultEnum, false
		return nil
	}

	return UnmarshalYAML(node, &e.Enum)
}

func (e Nullable[Enum]) Value() (driver.Value, error) {
	if !e.Valid {
		return nil, nil
	}

	return ValueSQL(e.Enum)
}

func (e *Nullable[Enum]) Scan(a any) error {
	if a == nil {
		var defaultEnum Enum
		e.Enum, e.Valid = defaultEnum, false
		return nil
	}

	e.Valid = true
	return ScanSQL(a, &e.Enum)
}
