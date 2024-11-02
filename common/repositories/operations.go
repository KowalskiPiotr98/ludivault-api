package repositories

import (
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/google/uuid"
)

type HasIdSetter interface {
	SetId(id uuid.UUID)
}

// QueryRows is a helper function to run a multiple rows query based on a query string.
func QueryRows[T any](connector gotabase.Connector, scanner func(row gotabase.Row) (*T, error), query string, args ...any) ([]*T, error) {
	rows, err := connector.QueryRows(query, args...)
	if err != nil {
		return nil, HandleKnownError(err)
	}
	return scanRows(rows, scanner)
}

// QueryRow is a helper function to run a single row query.
func QueryRow[T any](connector gotabase.Connector, scanner func(row gotabase.Row) (*T, error), query string, args ...any) (*T, error) {
	row, err := connector.QueryRow(query, args...)
	if err != nil {
		return nil, HandleKnownError(err)
	}

	result, err := scanner(row)
	if err != nil {
		return nil, HandleKnownError(err)
	}
	return result, nil
}

// CreateDataWithId creates a new data in the database.
// The query is expected to return a new id, that will be set in the object using the HasIdSetter interface method.
func CreateDataWithId[T HasIdSetter](connector gotabase.Connector, object T, query string, args ...any) error {
	row, err := connector.QueryRow(query, args...)
	if err != nil {
		return HandleKnownError(err)
	}

	var newId uuid.UUID
	if err = row.Scan(&newId); err != nil {
		return HandleKnownError(err)
	}

	object.SetId(newId)
	return nil
}

// UpdateRow runs the query to update a single row in the database.
// If none or more than one row are updated, then error will be returned.
func UpdateRow(connector gotabase.Connector, query string, args ...any) error {
	return runSingleRowAffectedQuery(connector, query, args...)
}

// DeleteRow runs the query to remove a single row from the database.
// If none or more than one row are affected, then error will be returned.
func DeleteRow(connector gotabase.Connector, query string, args ...any) error {
	return runSingleRowAffectedQuery(connector, query, args...)
}

func runSingleRowAffectedQuery(connector gotabase.Connector, query string, args ...any) error {
	result, err := connector.Exec(query, args...)
	if err != nil {
		return HandleKnownError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return HandleKnownError(err)
	}

	if rowsAffected > 1 {
		return TooManyRowsAffectedErr
	}
	if rowsAffected == 0 {
		return DataNotFoundErr
	}
	return nil
}

func scanRows[T any](rows gotabase.Rows, scanner func(row gotabase.Row) (*T, error)) ([]*T, error) {
	defer rows.Close()
	result := make([]*T, 0)

	for rows.Next() {
		item, err := scanner(rows)
		if err != nil {
			return nil, HandleKnownError(err)
		}

		result = append(result, item)
	}

	return result, nil
}
