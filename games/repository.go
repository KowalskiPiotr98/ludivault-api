package games

import (
	"fmt"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// GetGames returns a list of games with offset and limit used for pagination.
func GetGames(offset int, limit int, userId uuid.UUID, title string, owned *bool, released *bool, inProgress *bool) ([]*Game, error) {
	query := `select id, platform_id, title, owned, release_date, released from games where user_id = $3 %s order by title offset $1 limit $2`
	args := []interface{}{offset, limit, userId}

	if title != "" {
		args = append(args, fmt.Sprintf("%%%s%%", title))
		query = fmt.Sprintf(query, fmt.Sprintf("and title ilike $%d %%s", len(args)))
	}

	if owned != nil {
		args = append(args, *owned)
		query = fmt.Sprintf(query, fmt.Sprintf("and owned = $%d %%s ", len(args)))
	}

	if released != nil {
		args = append(args, *released)
		query = fmt.Sprintf(query, fmt.Sprintf("and released = $%d %%s ", len(args)))
	}

	if inProgress != nil {
		not := "not"
		if *inProgress {
			not = ""
		}
		query = fmt.Sprintf(query, fmt.Sprintf("and %s exists(select * from playthroughs where game_id = games.id and status = 0 limit 1) %%s", not))
	}

	query = fmt.Sprintf(query, "")
	return operations.QueryRows(getDatabase(), scanGame, query, args...)
}

// GetGame returns a single game selected by id.
func GetGame(id uuid.UUID, userId uuid.UUID) (*Game, error) {
	query := `select id, platform_id, title, owned, release_date, released from games where id = $1 and user_id = $2`
	return operations.QueryRow(getDatabase(), scanGame, query, id, userId)
}

// CreateGame creates a new game.
func CreateGame(game *Game, userId uuid.UUID) error {
	query := `insert into games (title, platform_id, owned, release_date, released, user_id) values ($1, $2, $3, $4, $5, $6) returning id`
	return operations.CreateRowWithId(getDatabase(), game, query, game.Title, game.PlatformId, game.Owned, game.ReleaseDate, game.Released, userId)
}

// UpdateGame updates details about a single game in the database.
func UpdateGame(game *Game, userId uuid.UUID) error {
	query := `update games set title = $2, platform_id = $3, owned = $4, release_date = $5, released = $6 where id = $1 and user_id = $7`
	return operations.UpdateRow(getDatabase(), query, game.Id, game.Title, game.PlatformId, game.Owned, game.ReleaseDate, game.Released, userId)
}

// DeleteGame deletes a single game from the database
func DeleteGame(id uuid.UUID, userId uuid.UUID) error {
	query := `delete from games where id = $1 and user_id = $2`
	return operations.DeleteRow(getDatabase(), query, id, userId)
}

func IsUserAuthorised(connector gotabase.Connector, gameId uuid.UUID, userId uuid.UUID) bool {
	query := `select count(1) from games where id = $1 and user_id = $2`
	row, err := connector.QueryRow(query, gameId, userId)
	if err != nil {
		log.Warnf("Failed to check user authorised: %v", err)
		return false
	}
	var count int
	err = row.Scan(&count)
	if err != nil {
		log.Warnf("Failed to check user authorised: %v", err)
		return false
	}
	return count > 0
}
