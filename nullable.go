package enum

import (
	"database/sql/driver"
)

// Nullable allows handling nullable enums in both JSON and SQL.
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
