package games

import (
	"database/sql"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/google/uuid"
)

type Game struct {
	Id         uuid.UUID
	PlatformId uuid.UUID

	Title       string
	Owned       bool
	ReleaseDate sql.NullTime
	Released    bool
}

func (g *Game) SetId(id uuid.UUID) {
	g.Id = id
}

func scanGame(row gotabase.Row) (*Game, error) {
	var game Game
	if err := row.Scan(&game.Id, &game.PlatformId, &game.Title, &game.Owned, &game.ReleaseDate, &game.Released); err != nil {
		return nil, err
	}
	return &game, nil
}
