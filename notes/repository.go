package notes

import (
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/KowalskiPiotr98/ludivault/games"
	"github.com/google/uuid"
)

func GetNoteTitles(gameId uuid.UUID, userId uuid.UUID) ([]*GameNote, error) {
	query := `select id, game_id, title, '', kind, added_on, pinned from game_notes where game_id = $1 and check_user_note($2, id)`
	return operations.QueryRows(getDatabase(), scanGameNote, query, gameId, userId)
}

func GetNote(noteId uuid.UUID, userId uuid.UUID) (*GameNote, error) {
	query := `select id, game_id, title, value, kind, added_on, pinned from game_notes where id = $1 and check_user_note($2, id)`
	return operations.QueryRow(getDatabase(), scanGameNote, query, noteId, userId)
}

func CreateNote(note *GameNote, userId uuid.UUID) error {
	if !games.IsUserAuthorised(getDatabase(), note.GameId, userId) {
		return operations.Errors.DataNotFoundErr
	}

	query := `insert into game_notes (game_id, title, value, kind, pinned) values ($1, $2, $3, $4, false) returning id`
	return operations.CreateRowWithId(getDatabase(), note, query, note.GameId, note.Title, note.Value, note.Kind)
}

func UpdateNote(note *GameNote, userId uuid.UUID) error {
	query := `update game_notes set title = $3, value = $4, kind = $5, pinned = $6 where id = $1 and check_user_note($2, id)`
	return operations.UpdateRow(getDatabase(), query, note.Id, userId, note.Title, note.Value, note.Kind, note.Pinned)
}

func DeleteNote(noteId uuid.UUID, userId uuid.UUID) error {
	query := `delete from game_notes where id = $1 and check_user_note($2, id)`
	return operations.DeleteRow(getDatabase(), query, noteId, userId)
}
