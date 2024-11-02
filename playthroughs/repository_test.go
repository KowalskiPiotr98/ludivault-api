package playthroughs

import (
	"database/sql"
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/KowalskiPiotr98/ludivault/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func makePlatform() uuid.UUID {
	id := tests.GetRandomUuid()
	query := `insert into platforms (id, name, short_name) values ($1, 'aa', 'aa')`
	_, err := getDatabase().Exec(query, id)
	tests.PanicOnErr(err)
	return id
}

func makeGame(title string, platformId uuid.UUID) uuid.UUID {
	id := tests.GetRandomUuid()
	query := `insert into games (id, title, platform_id, release_date, released) values ($1, $3, $2, null, true)`
	_, err := getDatabase().Exec(query, id, platformId, title)
	tests.PanicOnErr(err)
	return id
}

func TestCreatePlaythrough(t *testing.T) {
	t.Run("New playthrough created", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		playthrough := Playthrough{
			GameId:    makeGame("test", makePlatform()),
			StartDate: time.Now().UTC(),
			EndDate:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
			Status:    PlaythroughCompleted,
			Runtime:   123123,
		}

		err := CreatePlaythrough(&playthrough)

		assert.NoError(t, err)
		dbRow, err := getDatabase().QueryRow("select id, game_id, start_date, end_date, status, runtime_minutes from playthroughs where id = $1", playthrough.Id)
		tests.PanicOnErr(err)
		dbPlaythrough, err := scanPlaythrough(dbRow)
		tests.PanicOnErr(err)
		dbPlaythrough.StartDate = dbPlaythrough.StartDate.UTC()
		dbPlaythrough.EndDate.Time = dbPlaythrough.EndDate.Time.UTC()
		assert.Equal(t, playthrough, *dbPlaythrough)
	})

	t.Run("Game missing returns error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		playthrough := Playthrough{
			GameId:    tests.GetRandomUuid(),
			StartDate: time.Now().UTC(),
			EndDate:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
			Status:    PlaythroughCompleted,
			Runtime:   123123,
		}

		err := CreatePlaythrough(&playthrough)

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestGetPlaythroughs(t *testing.T) {
	makePlaythroughs := func() []Playthrough {
		platform := makePlatform()
		game1 := makeGame("game 1", platform)
		game2 := makeGame("game 2", platform)
		playthroughs := []Playthrough{
			{
				GameId:    game1,
				StartDate: time.Now().UTC(),
				EndDate:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
				Status:    PlaythroughCompleted,
				Runtime:   123123,
			},
			{
				GameId:    game1,
				StartDate: time.Now().UTC(),
				EndDate:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
				Status:    PlaythroughCompleted,
				Runtime:   123123,
			},
			{
				GameId:    game1,
				StartDate: time.Now().UTC(),
				EndDate:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
				Status:    PlaythroughCompleted,
				Runtime:   123123,
			},
			{
				GameId:    game2,
				StartDate: time.Now().UTC(),
				EndDate:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
				Status:    PlaythroughCompleted,
				Runtime:   123123,
			},
			{
				GameId:    game2,
				StartDate: time.Now().UTC(),
				EndDate:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
				Status:    PlaythroughCompleted,
				Runtime:   123123,
			},
			{
				GameId:    game2,
				StartDate: time.Now().UTC(),
				EndDate:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
				Status:    PlaythroughCompleted,
				Runtime:   123123,
			},
		}
		for i := range playthroughs {
			tests.PanicOnErr(CreatePlaythrough(&playthroughs[i]))
		}
		return playthroughs
	}

	t.Run("Returns all playthroughs", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		playthroughs := makePlaythroughs()

		list, err := GetPlaythroughs(uuid.Nil)

		assert.NoError(t, err)
		assert.Len(t, list, len(playthroughs))
		for _, p := range list {
			p.StartDate = p.StartDate.UTC()
			p.EndDate.Time = p.EndDate.Time.UTC()
			assert.Contains(t, playthroughs, *p)
		}
	})

	t.Run("Returns playthroughs limited to game", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		playthroughs := makePlaythroughs()

		list, err := GetPlaythroughs(playthroughs[0].GameId)

		filtered := playthroughs[:3]
		assert.NoError(t, err)
		assert.Len(t, list, len(filtered))
		for _, p := range list {
			p.StartDate = p.StartDate.UTC()
			p.EndDate.Time = p.EndDate.Time.UTC()
			assert.Contains(t, filtered, *p)
		}
	})
}

func TestGetPlaythrough(t *testing.T) {
	t.Run("Returns playthrough", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		playthrough := Playthrough{
			GameId:    makeGame("test", makePlatform()),
			StartDate: time.Now().UTC(),
			EndDate:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
			Status:    PlaythroughCompleted,
			Runtime:   123123,
		}
		tests.PanicOnErr(CreatePlaythrough(&playthrough))

		db, err := GetPlaythrough(playthrough.Id)

		assert.NoError(t, err)
		db.StartDate = playthrough.StartDate.UTC()
		db.EndDate.Time = playthrough.EndDate.Time.UTC()
		assert.Equal(t, playthrough, *db)
	})

	t.Run("Playthrough not found", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)

		_, err := GetPlaythrough(tests.GetRandomUuid())

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestUpdatePlaythrough(t *testing.T) {
	t.Run("Updates existing playthrough", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		playthrough := Playthrough{
			GameId:    makeGame("test", makePlatform()),
			StartDate: time.Now().UTC(),
			EndDate:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
			Status:    PlaythroughCompleted,
			Runtime:   123123,
		}
		tests.PanicOnErr(CreatePlaythrough(&playthrough))
		playthrough.StartDate = time.Now().AddDate(0, 1, 1).UTC()
		playthrough.EndDate = sql.NullTime{}
		playthrough.Status = PlaythroughSuspended
		playthrough.Runtime = 9

		err := UpdatePlaythrough(&playthrough)

		assert.NoError(t, err)
		dbPlaythrough, err := GetPlaythrough(playthrough.Id)
		tests.PanicOnErr(err)
		dbPlaythrough.StartDate = playthrough.StartDate.UTC()
		assert.Equal(t, playthrough, *dbPlaythrough)
	})

	t.Run("Playthrough not found", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)

		err := UpdatePlaythrough(&Playthrough{Id: tests.GetRandomUuid()})

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestDeletePlaythrough(t *testing.T) {
	t.Run("Deletes existing playthrough", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		playthrough := Playthrough{
			GameId:    makeGame("test", makePlatform()),
			StartDate: time.Now().UTC(),
			EndDate:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
			Status:    PlaythroughCompleted,
			Runtime:   123123,
		}
		tests.PanicOnErr(CreatePlaythrough(&playthrough))

		err := DeletePlaythrough(playthrough.Id)

		assert.NoError(t, err)
		_, err = GetPlaythrough(playthrough.Id)
		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})

	t.Run("Playthrough not found", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)

		err := DeletePlaythrough(tests.GetRandomUuid())

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}
