package repositories

import "github.com/KowalskiPiotr98/gotabase"

// ScanRows uses the provided scanner function to read multiple rows of type T from database to a slice.
func ScanRows[T any](rows gotabase.Rows, scanner func(row gotabase.Row) (*T, error)) ([]*T, error) {
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
