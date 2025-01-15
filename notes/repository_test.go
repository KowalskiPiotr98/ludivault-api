package notes

import (
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/KowalskiPiotr98/ludivault/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func makePlatform(userId uuid.UUID) uuid.UUID {
	id := tests.GetRandomUuid()
	query := `insert into platforms (id, name, short_name, user_id) values ($1, 'aa', 'aa', $2)`
	_, err := getDatabase().Exec(query, id, userId)
	tests.PanicOnErr(err)
	return id
}

func makeGame(title string, platformId uuid.UUID, userId uuid.UUID) uuid.UUID {
	id := tests.GetRandomUuid()
	query := `insert into games (id, title, platform_id, release_date, released, user_id) values ($1, $3, $2, null, true, $4)`
	_, err := getDatabase().Exec(query, id, platformId, title, userId)
	tests.PanicOnErr(err)
	return id
}

func TestCreateNote(t *testing.T) {
	t.Run("New note created", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		note := GameNote{
			GameId: makeGame("test", makePlatform(userId), userId),
			Title:  "test title",
			Value:  "test value",
			Kind:   TextNote,
			Pinned: false,
		}

		err := CreateNote(&note, userId)

		assert.NoError(t, err)
		dbRow, err := getDatabase().QueryRow("select id, game_id, title, value, kind, added_on, pinned from game_notes where id = $1", note.Id)
		tests.PanicOnErr(err)
		dbNote, err := scanGameNote(dbRow)
		tests.PanicOnErr(err)
		dbNote.AddedOn = note.AddedOn
		assert.Equal(t, note, *dbNote)
	})

	t.Run("Game missing returns error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		note := GameNote{
			GameId: tests.GetRandomUuid(),
			Title:  "test title",
			Value:  "test value",
			Kind:   TextNote,
			Pinned: false,
		}

		err := CreateNote(&note, userId)

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})

	t.Run("User not authorised for game", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		note := GameNote{
			GameId: makeGame("test", makePlatform(userId), userId),
			Title:  "test title",
			Value:  "test value",
			Kind:   TextNote,
			Pinned: false,
		}

		err := CreateNote(&note, tests.MakeTestUserId(getDatabase()))

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestGetNoteTitles(t *testing.T) {
	makeNotes := func(userId uuid.UUID) []GameNote {
		platform := makePlatform(userId)
		game1 := makeGame("game 1", platform, userId)
		notes := []GameNote{
			{
				GameId: game1,
				Title:  "test title 1",
				Value:  "aaa",
				Kind:   TextNote,
				Pinned: false,
			},
			{
				GameId: game1,
				Title:  "test title 2",
				Value:  "aaa",
				Kind:   TextNote,
				Pinned: false,
			},
			{
				GameId: game1,
				Title:  "test title 3",
				Value:  "aaa",
				Kind:   TextNote,
				Pinned: false,
			},
		}
		for i := range notes {
			tests.PanicOnErr(CreateNote(&notes[i], userId))
			notes[i].AddedOn = time.Time{}
		}

		otherUser := tests.MakeTestUserId(getDatabase())
		gameUnauthorised := makeGame("unauthorised", platform, otherUser)
		tests.PanicOnErr(CreateNote(&GameNote{
			GameId: gameUnauthorised,
			Title:  "aaa",
			Value:  "aaa",
			Kind:   TextNote,
			Pinned: false,
		}, otherUser))

		return notes
	}

	t.Run("Returns all notes for game", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		notes := makeNotes(userId)

		list, err := GetNoteTitles(notes[0].GameId, userId)

		assert.NoError(t, err)
		assert.Len(t, list, len(notes))
		for _, p := range list {
			p.Value = "aaa"
			p.AddedOn = time.Time{}
			assert.Contains(t, notes, *p)
		}
	})
}

func TestGetNote(t *testing.T) {
	t.Run("Returns note", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		note := GameNote{
			GameId: makeGame("test", makePlatform(userId), userId),
			Title:  "test title 1",
			Value:  "aaa",
			Kind:   TextNote,
			Pinned: false,
		}
		tests.PanicOnErr(CreateNote(&note, userId))

		db, err := GetNote(note.Id, userId)

		assert.NoError(t, err)
		db.AddedOn = note.AddedOn
		assert.Equal(t, note, *db)
	})

	t.Run("Note not found", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)

		_, err := GetNote(tests.GetRandomUuid(), tests.MakeTestUserId(getDatabase()))

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})

	t.Run("User not authorised for game", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		note := GameNote{
			GameId: makeGame("test", makePlatform(userId), userId),
			Title:  "test title 1",
			Value:  "aaa",
			Kind:   TextNote,
			Pinned: false,
		}
		tests.PanicOnErr(CreateNote(&note, userId))

		_, err := GetNote(note.Id, tests.MakeTestUserId(getDatabase()))

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestUpdateNote(t *testing.T) {
	t.Run("Updates existing note", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		note := GameNote{
			GameId: makeGame("test", makePlatform(userId), userId),
			Title:  "test title 1",
			Value:  "aaa",
			Kind:   TextNote,
			Pinned: false,
		}
		tests.PanicOnErr(CreateNote(&note, userId))
		note.Title = "modified"
		note.Value = "modified value"
		note.Kind = LinkNote
		note.Pinned = false

		err := UpdateNote(&note, userId)

		assert.NoError(t, err)
		dbNote, err := GetNote(note.Id, userId)
		tests.PanicOnErr(err)
		dbNote.AddedOn = note.AddedOn
		assert.Equal(t, note, *dbNote)
	})

	t.Run("Note not found", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)

		err := UpdateNote(&GameNote{Id: tests.GetRandomUuid()}, tests.MakeTestUserId(getDatabase()))

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})

	t.Run("User not authorised for game", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		note := GameNote{
			GameId: makeGame("test", makePlatform(userId), userId),
			Title:  "test title 1",
			Value:  "aaa",
			Kind:   TextNote,
			Pinned: false,
		}
		tests.PanicOnErr(CreateNote(&note, userId))

		err := UpdateNote(&note, tests.MakeTestUserId(getDatabase()))

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestDeleteNote(t *testing.T) {
	t.Run("Deletes existing note", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		note := GameNote{
			GameId: makeGame("test", makePlatform(userId), userId),
			Title:  "test title 1",
			Value:  "aaa",
			Kind:   TextNote,
			Pinned: false,
		}
		tests.PanicOnErr(CreateNote(&note, userId))

		err := DeleteNote(note.Id, userId)

		assert.NoError(t, err)
		_, err = GetNote(note.Id, userId)
		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})

	t.Run("Note not found", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)

		err := DeleteNote(tests.GetRandomUuid(), tests.MakeTestUserId(getDatabase()))

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})

	t.Run("User not authorised for game", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		note := GameNote{
			GameId: makeGame("test", makePlatform(userId), userId),
			Title:  "test title 1",
			Value:  "aaa",
			Kind:   TextNote,
			Pinned: false,
		}
		tests.PanicOnErr(CreateNote(&note, userId))

		err := DeleteNote(note.Id, tests.MakeTestUserId(getDatabase()))

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}
