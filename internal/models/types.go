package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type DateOnly struct {
	time.Time
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		return nil
	}
	// Remove quotes
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	return d.parseString(s)
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Time.Format("2006-01-02") + `"`), nil
}

// Value implements the driver Valuer interface.
func (d DateOnly) Value() (driver.Value, error) {
	return d.Time, nil
}

// Scan implements the Scanner interface.
func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case []byte:
		return d.parseString(string(v))
	case string:
		return d.parseString(v)
	default:
		return fmt.Errorf("cannot scan type %T into DateOnly", value)
	}
}

func (d *DateOnly) parseString(s string) error {
	if t, err := time.Parse("2006-01-02", s); err == nil {
		d.Time = t
		return nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		d.Time = t
		return nil
	}
	return fmt.Errorf("failed to parse date: %s", s)
}
