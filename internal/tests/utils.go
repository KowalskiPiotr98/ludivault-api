package tests

import (
	"github.com/google/uuid"
)

func GetRandomUuid() uuid.UUID {
	id, _ := uuid.NewUUID()
	return id
}
