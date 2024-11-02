package playthroughs

import (
	"database/sql"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/google/uuid"
	"time"
)

type PlaythroughStatus int8

const (
	PlaythroughInProgress PlaythroughStatus = iota
	PlaythroughCompleted
	PlaythroughDropped
	PlaythroughRetired
	PlaythroughSuspended
)

type Playthrough struct {
	Id        uuid.UUID
	GameId    uuid.UUID
	StartDate time.Time
	EndDate   sql.NullTime
	Status    PlaythroughStatus
	Runtime   int
}

func (p *Playthrough) SetId(id uuid.UUID) {
	p.Id = id
}

func scanPlaythrough(row gotabase.Row) (*Playthrough, error) {
	var p Playthrough
	if err := row.Scan(&p.Id, p.GameId, &p.StartDate, &p.EndDate, &p.Status, &p.Runtime); err != nil {
		return nil, err
	}
	return &p, nil
}
