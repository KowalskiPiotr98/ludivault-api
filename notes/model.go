package notes

import (
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/google/uuid"
	"time"
)

type NoteType int8

const (
	TextNote NoteType = iota
	LinkNote
)

type GameNote struct {
	Id      uuid.UUID
	GameId  uuid.UUID
	Title   string
	Value   string
	Kind    NoteType
	AddedOn time.Time
	Pinned  bool
}

func (g *GameNote) SetId(id uuid.UUID) {
	g.Id = id
}

func scanGameNote(row gotabase.Row) (*GameNote, error) {
	var n GameNote
	if err := row.Scan(&n.Id, &n.GameId, &n.Title, &n.Value, &n.Kind, &n.AddedOn, &n.Pinned); err != nil {
		return nil, err
	}
	return &n, nil
}
