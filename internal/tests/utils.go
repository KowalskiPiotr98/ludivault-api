package tests

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func GetRandomUuid() uuid.UUID {
	id, _ := uuid.NewUUID()
	return id
}

func PanicOnError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
