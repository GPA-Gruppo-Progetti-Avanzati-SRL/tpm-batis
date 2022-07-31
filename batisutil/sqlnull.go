package batisutil

import (
	"database/sql"
	"time"
)

// ToSqlNullTime transforms a time.Time into a sql.NullTime. If the tm is empty it returns an empty struct.
func ToSqlNullTime(tm time.Time) sql.NullTime {
	if tm.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{Time: tm, Valid: true}
}

// ToSqlNullString transforms a string into a sql.NullString. If the string is empty it returns an empty struct.
func ToSqlNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}

	return sql.NullString{String: s, Valid: true}
}
