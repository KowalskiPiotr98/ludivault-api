package dto

import (
	"database/sql"
	"time"
)

func MapMany[TSource any, TDest any](source []*TSource, mapper func(*TSource) *TDest) []*TDest {
	result := make([]*TDest, len(source))
	for i, tSource := range source {
		result[i] = mapper(tSource)
	}
	return result
}

func makeNullTimeFromPointer(time *time.Time) sql.NullTime {
	if time == nil {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Valid: true,
		Time:  *time,
	}
}

func makePointerFromNullTime(time sql.NullTime) *time.Time {
	if time.Valid {
		return &time.Time
	}
	return nil
}

func makePointerFromNullInt(value sql.NullInt32) *int {
	if value.Valid {
		typedInt := int(value.Int32)
		return &typedInt
	}
	return nil
}

func makeNullIntFromPointer(value *int) sql.NullInt32 {
	if value == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Valid: true,
		Int32: int32(*value),
	}
}
