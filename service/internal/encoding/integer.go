package encoding

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// NullInt16 represents an int16 that may be null. NullInt16 implements the
// sql.Scanner interface so that it may be used as a scan destination, similar to
// sql.NullString.
type NullInt16 struct {
	Int16 int16

	// Valid is true if Int16 is not NULL
	Valid bool
}

var _ sql.Scanner = (*NullInt16)(nil)
var _ driver.Valuer = (*NullInt16)(nil)

// Scan implements the Scanner interface.
func (ni *NullInt16) Scan(value interface{}) error {
	if value == nil {
		ni.Int16, ni.Valid = 0, false
		return nil
	}
	i64, ok := value.(int64)
	if !ok {
		return fmt.Errorf("cannot convert type %T to int16", value)
	}
	ni.Int16 = int16(i64)
	ni.Valid = true
	return nil
}

// Value implements the driver Valuer interface.
func (ni NullInt16) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int16, nil
}
