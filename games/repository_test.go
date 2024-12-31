package games

import (
	"database/sql"
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/KowalskiPiotr98/ludivault/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func makePlatform(userId uuid.UUID) uuid.UUID {
	id := tests.GetRandomUuid()
	query := `insert into platforms (id, name, short_name, user_id) values ($1, 'aa', 'aa', $2)`
	_, err := getDatabase().Exec(query, id, userId)
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
		userId := tests.MakeTestUserId(getDatabase())
		game := makeDefaultTestGame(makePlatform(userId))

		err := CreateGame(&game, userId)

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
		userId := tests.MakeTestUserId(getDatabase())
		game := makeDefaultTestGame(tests.GetRandomUuid())

		err := CreateGame(&game, userId)

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestGetGames(t *testing.T) {
	t.Run("Get all games", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		platformId := makePlatform(userId)
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
			tests.PanicOnErr(CreateGame(&games[i], userId))
		}

		list, err := GetGames(0, 100, userId)

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
		userId := tests.MakeTestUserId(getDatabase())
		game := makeDefaultTestGame(makePlatform(userId))
		tests.PanicOnErr(CreateGame(&game, userId))

		dbGame, err := GetGame(game.Id, userId)

		assert.NoError(t, err)
		dbGame.ReleaseDate.Time = game.ReleaseDate.Time.UTC()
		assert.Equal(t, game, *dbGame)
	})
	t.Run("Game does not exist - error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())

		_, err := GetGame(tests.GetRandomUuid(), userId)

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestUpdateGame(t *testing.T) {
	t.Run("Game exists - updated", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		game := makeDefaultTestGame(makePlatform(userId))
		tests.PanicOnErr(CreateGame(&game, userId))
		game.Title = "updated game"
		game.Owned = false
		game.ReleaseDate = sql.NullTime{}
		game.Released = false

		err := UpdateGame(&game, userId)

		assert.NoError(t, err)
		dbGame, err := GetGame(game.Id, userId)
		tests.PanicOnErr(err)
		assert.Equal(t, game, *dbGame)
	})
	t.Run("Game not found - error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		game := Game{Id: tests.GetRandomUuid()}
		userId := tests.MakeTestUserId(getDatabase())

		err := UpdateGame(&game, userId)

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestDeleteGame(t *testing.T) {
	t.Run("Game exists - deleted", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		game := makeDefaultTestGame(makePlatform(userId))
		tests.PanicOnErr(CreateGame(&game, userId))

		err := DeleteGame(game.Id, userId)

		assert.NoError(t, err)
		_, err = GetGame(game.Id, userId)
		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
	t.Run("Game not found - error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())

		err := DeleteGame(tests.GetRandomUuid(), userId)

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}
