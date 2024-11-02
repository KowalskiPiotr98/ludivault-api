package platforms

import (
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/google/uuid"
)

// GetPlatforms returns a complete list of [Platform] items from the database.
func GetPlatforms() ([]*Platform, error) {
	query := `select id, name, short_name from platforms order by name`
	return operations.QueryRows(getDatabase(), scanPlatform, query)
}

// GetPlatform returns a single [Platform] based on the id provided.
func GetPlatform(id uuid.UUID) (*Platform, error) {
	query := `select id, name, short_name from platforms where id = $1`
	return operations.QueryRow(getDatabase(), scanPlatform, query, id)
}

// CreatePlatform creates a new [Platform] in the database and sets the id in the provided struct.
func CreatePlatform(platform *Platform) error {
	query := `insert into platforms (name, short_name) values ($1, $2) returning id`
	return operations.CreateRowWithId(getDatabase(), platform, query, platform.Name, platform.ShortName)
}

// UpdatePlatform updates values of the [Platform] with the id as provided.
func UpdatePlatform(platform *Platform) error {
	query := `update platforms set name = $2, short_name = $3 where id = $1`
	return operations.UpdateRow(getDatabase(), query, platform.Id, platform.Name, platform.ShortName)
}

// DeletePlatform removes a single [Platform] with the id provided from the database.
func DeletePlatform(id uuid.UUID) error {
	query := `delete from platforms where id = $1`
	return operations.DeleteRow(getDatabase(), query, id)
}
