package dto

import (
	"github.com/KowalskiPiotr98/ludivault/notes"
	"github.com/google/uuid"
	"time"
)

type NoteDto struct {
	Id      uuid.UUID `json:"id"`
	GameId  uuid.UUID `json:"gameId"`
	Title   string    `json:"title"`
	Value   string    `json:"value"`
	Kind    int8      `json:"kind"`
	AddedOn time.Time `json:"addedOn"`
	Pinned  bool      `json:"pinned"`
}

func MapNoteToDto(note *notes.GameNote) *NoteDto {
	return &NoteDto{
		Id:      note.Id,
		GameId:  note.GameId,
		Title:   note.Title,
		Value:   note.Value,
		Kind:    int8(note.Kind),
		AddedOn: note.AddedOn,
		Pinned:  note.Pinned,
	}
}

type ModifyNoteDto struct {
	GameId uuid.UUID `json:"gameId" binding:"required"`
	Title  string    `json:"title" binding:"required,max=100"`
	Value  string    `json:"value" binding:"required,max=10000"`
	Kind   int8      `json:"kind" binding:"min=0,max=1"`
	Pinned bool      `json:"pinned"`
}

func MapCreateDtoToNote(note *ModifyNoteDto) *notes.GameNote {
	return &notes.GameNote{
		GameId: note.GameId,
		Title:  note.Title,
		Value:  note.Value,
		Kind:   notes.NoteType(note.Kind),
		Pinned: note.Pinned,
	}
}

func MapModifyDtoToNote(noteId uuid.UUID, note *ModifyNoteDto) *notes.GameNote {
	return &notes.GameNote{
		Id:     noteId,
		GameId: note.GameId,
		Title:  note.Title,
		Value:  note.Value,
		Kind:   notes.NoteType(note.Kind),
		Pinned: note.Pinned,
	}
}
