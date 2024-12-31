package users

import (
	"github.com/KowalskiPiotr98/ludivault/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("new user created", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		user := &User{
			Email:        "user@localhost",
			ProviderId:   "localhost",
			ProviderName: "test",
		}

		err := GetOrCreate(user)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, user.Id)
	})

	t.Run("existing user updated", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		user := &User{
			Email:        "user@localhost",
			ProviderId:   "localhost",
			ProviderName: "test",
		}
		tests.PanicOnErr(GetOrCreate(user))

		user.Email = "user2@localhost"
		user.Id = uuid.Nil
		err := GetOrCreate(user)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, user.Id)
	})
}
