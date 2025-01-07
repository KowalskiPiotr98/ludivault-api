package games

import (
	"database/sql"
	"fmt"
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

func makePlaythrough(gameId uuid.UUID, status int) uuid.UUID {
	id := tests.GetRandomUuid()
	query := `insert into playthroughs (game_id, status, start_date) values ($1, $2, now())`
	_, err := getDatabase().Exec(query, gameId, status)
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
	tests.GetDatabaseWithCleanup(t)
	userId := tests.MakeTestUserId(getDatabase())
	platformId := makePlatform(userId)
	games := []Game{
		{
			PlatformId:  platformId,
			Title:       "test game 1",
			Owned:       false,
			ReleaseDate: sql.NullTime{Valid: true, Time: tests.GetRandomTestTime()},
			Released:    true,
		},
		{
			PlatformId:  platformId,
			Title:       "test game 2",
			Owned:       true,
			ReleaseDate: sql.NullTime{Valid: true, Time: tests.GetRandomTestTime()},
			Released:    false,
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
	makePlaythrough(games[0].Id, 0)
	makePlaythrough(games[1].Id, 1)
	clone := games[0]
	tests.PanicOnErr(CreateGame(&clone, tests.MakeTestUserId(getDatabase())))

	t.Run("Get all games", func(t *testing.T) {
		list, err := GetGames(0, 100, userId, "game", nil, nil, nil)

		assert.NoError(t, err)
		assert.Len(t, list, len(games))
		for _, game := range list {
			game.ReleaseDate.Time = game.ReleaseDate.Time.UTC()
			assert.Contains(t, games, *game)
		}
	})

	bools := []bool{true, false}

	t.Run("Get only owned", func(t *testing.T) {
		for _, arg := range bools {
			t.Run(fmt.Sprint(arg), func(t *testing.T) {
				list, err := GetGames(0, 100, userId, "", tests.GetPointerFromValue(arg), nil, nil)

				assert.NoError(t, err)
				assert.NotEmpty(t, list)
				for _, game := range list {
					assert.Equal(t, arg, game.Owned)
				}
			})
		}
	})

	t.Run("Get only released", func(t *testing.T) {
		for _, arg := range bools {
			t.Run(fmt.Sprint(arg), func(t *testing.T) {
				list, err := GetGames(0, 100, userId, "", nil, tests.GetPointerFromValue(arg), nil)

				assert.NoError(t, err)
				assert.NotEmpty(t, list)
				for _, game := range list {
					assert.Equal(t, arg, game.Released)
				}
			})
		}
	})

	t.Run("Only in progress", func(t *testing.T) {
		list, err := GetGames(0, 100, userId, "", nil, nil, tests.GetPointerFromValue(true))

		assert.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, games[0].Id, list[0].Id)
	})

	t.Run("Only not in progress", func(t *testing.T) {
		list, err := GetGames(0, 100, userId, "", nil, nil, tests.GetPointerFromValue(false))

		assert.NoError(t, err)
		assert.Len(t, list, 2)
		for _, game := range list {
			assert.NotEqual(t, games[0].Id, game.Id)
		}
	})

	t.Run("Filter by title", func(t *testing.T) {
		list, err := GetGames(0, 100, userId, "1", nil, nil, nil)

		assert.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, games[0].Id, list[0].Id)
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

	t.Run("User not authorised", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		game := makeDefaultTestGame(makePlatform(userId))
		tests.PanicOnErr(CreateGame(&game, userId))

		_, err := GetGame(game.Id, tests.MakeTestUserId(getDatabase()))

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

	t.Run("User not authorised", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		game := makeDefaultTestGame(makePlatform(userId))
		tests.PanicOnErr(CreateGame(&game, userId))
		game.Title = "updated game"
		game.Owned = false
		game.ReleaseDate = sql.NullTime{}
		game.Released = false

		err := UpdateGame(&game, tests.MakeTestUserId(getDatabase()))

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

	t.Run("User not authorised", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		game := makeDefaultTestGame(makePlatform(userId))
		tests.PanicOnErr(CreateGame(&game, userId))

		err := DeleteGame(tests.GetRandomUuid(), tests.MakeTestUserId(getDatabase()))

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}
