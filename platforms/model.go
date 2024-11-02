package platforms

import (
	"github.com/google/uuid"
)

type Platform struct {
	Id        uuid.UUID
	Name      string
	ShortName string
}
