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
