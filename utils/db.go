package utils

import (
	"database/sql"
	"time"
)

func MakeNullTime(value *time.Time) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Valid: true,
		Time:  *value,
	}
}
