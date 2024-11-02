package platforms

import (
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/google/uuid"
)

type Platform struct {
	Id        uuid.UUID
	Name      string
	ShortName string
}

func (p *Platform) SetId(id uuid.UUID) {
	p.Id = id
}

func scanPlatform(row gotabase.Row) (*Platform, error) {
	var platform Platform
	err := row.Scan(&platform.Id, &platform.Name, &platform.ShortName)
	return &platform, err
}
