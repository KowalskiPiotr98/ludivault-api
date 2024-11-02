package tests

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func GetRandomUuid() uuid.UUID {
	id, _ := uuid.NewUUID()
	return id
}

func GetRandomTestTime() time.Time {
	return time.Date(2000, time.Month((rand.Int()%12)+1), (rand.Int()%20)+1, 0, 0, 0, 0, time.UTC)
}
