package games

import (
	"database/sql"
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/KowalskiPiotr98/ludivault/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func makePlatform() uuid.UUID {
	id := tests.GetRandomUuid()
	query := `insert into platforms (id, name, short_name) values ($1, 'aa', 'aa')`
	_, err := getDatabase().Exec(query, id)
	tests.PanicOnErr(err)
	return id
}

func makeDefaultTestGame(platformId uuid.UUID) Game {
	return Game{
		PlatformId:  platformId,
		Title:       "test game",
		Owned:       true,
		ReleaseDate: sql.NullTime{Valid: true, Time: tests.GetRandomTestTime()},
		Released:    true,
	}
}

func TestCreateGame(t *testing.T) {
	t.Run("New game created", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		game := makeDefaultTestGame(makePlatform())

		err := CreateGame(&game)

		assert.NoError(t, err)
		dbRow, err := getDatabase().QueryRow("select id, platform_id, title, owned, release_date, released from games where id = $1", game.Id)
		tests.PanicOnErr(err)
		dbGame, err := scanGame(dbRow)
		dbGame.ReleaseDate.Time = dbGame.ReleaseDate.Time.UTC()
		tests.PanicOnErr(err)
		assert.Equal(t, game, *dbGame)
	})

	t.Run("Missing platform", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		game := makeDefaultTestGame(tests.GetRandomUuid())

		err := CreateGame(&game)

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestGetGames(t *testing.T) {
	t.Run("Get all games", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		platformId := makePlatform()
		games := []Game{
			{
				PlatformId:  platformId,
				Title:       "test game 1",
				Owned:       true,
				ReleaseDate: sql.NullTime{Valid: true, Time: tests.GetRandomTestTime()},
				Released:    true,
			},
			{
				PlatformId:  platformId,
				Title:       "test game 2",
				Owned:       true,
				ReleaseDate: sql.NullTime{Valid: true, Time: tests.GetRandomTestTime()},
				Released:    true,
			},
			{
				PlatformId:  platformId,
				Title:       "test game 3",
				Owned:       true,
				ReleaseDate: sql.NullTime{Valid: true, Time: tests.GetRandomTestTime()},
				Released:    true,
			},
		}
		for i := range games {
			tests.PanicOnErr(CreateGame(&games[i]))
		}

		list, err := GetGames(0, 100)

		assert.NoError(t, err)
		assert.Len(t, list, len(games))
		for _, game := range list {
			game.ReleaseDate.Time = game.ReleaseDate.Time.UTC()
			assert.Contains(t, games, *game)
		}
	})
}

func TestGetGame(t *testing.T) {
	t.Run("Game exists - returned", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		game := makeDefaultTestGame(makePlatform())
		tests.PanicOnErr(CreateGame(&game))

		dbGame, err := GetGame(game.Id)

		assert.NoError(t, err)
		dbGame.ReleaseDate.Time = game.ReleaseDate.Time.UTC()
		assert.Equal(t, game, *dbGame)
	})
	t.Run("Game does not exist - error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)

		_, err := GetGame(tests.GetRandomUuid())

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestUpdateGame(t *testing.T) {
	t.Run("Game exists - updated", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		game := makeDefaultTestGame(makePlatform())
		tests.PanicOnErr(CreateGame(&game))
		game.Title = "updated game"
		game.Owned = false
		game.ReleaseDate = sql.NullTime{}
		game.Released = false

		err := UpdateGame(&game)

		assert.NoError(t, err)
		dbGame, err := GetGame(game.Id)
		tests.PanicOnErr(err)
		assert.Equal(t, game, *dbGame)
	})
	t.Run("Game not found - error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		game := Game{Id: tests.GetRandomUuid()}

		err := UpdateGame(&game)

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestDeleteGame(t *testing.T) {
	t.Run("Game exists - deleted", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		game := makeDefaultTestGame(makePlatform())
		tests.PanicOnErr(CreateGame(&game))

		err := DeleteGame(game.Id)

		assert.NoError(t, err)
		_, err = GetGame(game.Id)
		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
	t.Run("Game not found - error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)

		err := DeleteGame(tests.GetRandomUuid())

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}
