package games

import (
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/google/uuid"
)

// GetGames returns a list of games with offset and limit used for pagination.
func GetGames(offset int, limit int) ([]*Game, error) {
	query := `select id, platform_id, title, owned, release_date, released from games order by title offset $1 limit $2`
	return operations.QueryRows(getDatabase(), scanGame, query, offset, limit)
}

// GetGame returns a single game selected by id.
func GetGame(id uuid.UUID) (*Game, error) {
	query := `select id, platform_id, title, owned, release_date, released from games where id = $1`
	return operations.QueryRow(getDatabase(), scanGame, query, id)
}

// CreateGame creates a new game.
func CreateGame(game *Game) error {
	query := `insert into games (title, platform_id, owned, release_date, released) values ($1, $2, $3, $4, $5) returning id`
	return operations.CreateRowWithId(getDatabase(), game, query, game.Title, game.PlatformId, game.Owned, game.ReleaseDate, game.Released)
}

// UpdateGame updates details about a single game in the database.
func UpdateGame(game *Game) error {
	query := `update games set title = $2, platform_id = $3, owned = $4, release_date = $5, released = $6 where id = $1`
	return operations.UpdateRow(getDatabase(), query, game.Id, game.Title, game.PlatformId, game.Owned, game.ReleaseDate, game.Released)
}

// DeleteGame deletes a single game from the database
func DeleteGame(id uuid.UUID) error {
	query := `delete from games where id = $1`
	return operations.DeleteRow(getDatabase(), query, id)
}
