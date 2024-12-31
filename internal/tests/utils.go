package tests

import (
	"github.com/KowalskiPiotr98/gotabase"
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

func MakeTestUserId(connector gotabase.Connector) uuid.UUID {
	id := GetRandomUuid()
	MakeTestUser(connector, id)
	return id
}

func MakeTestUser(connector gotabase.Connector, id uuid.UUID) {
	query := `insert into users (id, provider_id, provider_name, email) values ($1, $2, $3, $4)`
	_, err := connector.Exec(query, id, id.String(), id.String(), id.String())
	PanicOnErr(err)
}
