package playthroughs

import (
	"fmt"
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/KowalskiPiotr98/ludivault/games"
	"github.com/google/uuid"
)

// GetPlaythroughs returns a list of playthroughs.
func GetPlaythroughs(gameId uuid.UUID, userId uuid.UUID) ([]*Playthrough, error) {
	query := `select id, game_id, start_date, end_date, status, runtime_minutes from playthroughs where check_user_playthrough($1, id) %s order by start_date desc`
	args := make([]interface{}, 1)
	args[0] = userId

	if gameId != uuid.Nil {
		query = fmt.Sprintf(query, "and game_id = $2 %s")
		args = append(args, gameId)
	}
	query = fmt.Sprintf(query, "")

	return operations.QueryRows(getDatabase(), scanPlaythrough, query, args...)
}

// GetPlaythrough returns a single playthrough selected by id.
func GetPlaythrough(id uuid.UUID, userId uuid.UUID) (*Playthrough, error) {
	query := `select id, game_id, start_date, end_date, status, runtime_minutes from playthroughs where id = $1 and check_user_playthrough($2, id)`
	return operations.QueryRow(getDatabase(), scanPlaythrough, query, id, userId)
}

// CreatePlaythrough creates a new playthrough.
func CreatePlaythrough(playthrough *Playthrough, userId uuid.UUID) error {
	if !games.IsUserAuthorised(getDatabase(), playthrough.GameId, userId) {
		return operations.Errors.DataNotFoundErr
	}

	query := `insert into playthroughs (game_id, start_date, end_date, status, runtime_minutes) values ($1, $2, $3, $4, $5) returning id`
	return operations.CreateRowWithId(getDatabase(), playthrough, query, playthrough.GameId, playthrough.StartDate, playthrough.EndDate, playthrough.Status, playthrough.Runtime)
}

// UpdatePlaythrough updates details about a single playthrough in the database.
func UpdatePlaythrough(playthrough *Playthrough, userId uuid.UUID) error {
	query := `update playthroughs set start_date = $2, end_date = $3, status = $4, runtime_minutes = $5 where id = $1 and check_user_playthrough($6, id)`
	return operations.UpdateRow(getDatabase(), query, playthrough.Id, playthrough.StartDate, playthrough.EndDate, playthrough.Status, playthrough.Runtime, userId)
}

// DeletePlaythrough deletes a single playthrough from the database
func DeletePlaythrough(id uuid.UUID, userId uuid.UUID) error {
	query := `delete from playthroughs where id = $1 and check_user_playthrough($2, id)`
	return operations.DeleteRow(getDatabase(), query, id, userId)
}
