package platforms

import (
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/google/uuid"
)

// GetPlatforms returns a complete list of [Platform] items from the database.
func GetPlatforms(userId uuid.UUID) ([]*Platform, error) {
	query := `select id, name, short_name from platforms where user_id = $1 order by name`
	return operations.QueryRows(getDatabase(), scanPlatform, query, userId)
}

// GetPlatform returns a single [Platform] based on the id provided.
func GetPlatform(id uuid.UUID, userId uuid.UUID) (*Platform, error) {
	query := `select id, name, short_name from platforms where id = $1 and user_id = $2`
	return operations.QueryRow(getDatabase(), scanPlatform, query, id, userId)
}

// CreatePlatform creates a new [Platform] in the database and sets the id in the provided struct.
func CreatePlatform(platform *Platform, userId uuid.UUID) error {
	query := `insert into platforms (name, short_name, user_id) values ($1, $2, $3) returning id`
	return operations.CreateRowWithId(getDatabase(), platform, query, platform.Name, platform.ShortName, userId)
}

// UpdatePlatform updates values of the [Platform] with the id as provided.
func UpdatePlatform(platform *Platform, userId uuid.UUID) error {
	query := `update platforms set name = $3, short_name = $4 where id = $1 and user_id = $2`
	return operations.UpdateRow(getDatabase(), query, platform.Id, userId, platform.Name, platform.ShortName)
}

// DeletePlatform removes a single [Platform] with the id provided from the database.
func DeletePlatform(id uuid.UUID, userId uuid.UUID) error {
	query := `delete from platforms where id = $1 and user_id = $2`
	return operations.DeleteRow(getDatabase(), query, id, userId)
}
