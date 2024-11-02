package platforms

import (
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/KowalskiPiotr98/ludivault/internal/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePlatform(t *testing.T) {
	t.Run("New platform created", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		newPlatform := Platform{
			Name:      "test platform",
			ShortName: "tp",
		}

		err := CreatePlatform(&newPlatform)

		assert.NoError(t, err)
		dbRow, err := getDatabase().QueryRow("select id, name, short_name from platforms where id = $1", newPlatform.Id)
		tests.PanicOnErr(err)
		var dbPlatform Platform
		tests.PanicOnErr(dbRow.Scan(&dbPlatform.Id, &dbPlatform.Name, &dbPlatform.ShortName))
		assert.Equal(t, newPlatform, dbPlatform)
	})

	t.Run("Platform already exists", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		newPlatform := Platform{
			Name:      "test platform",
			ShortName: "tp",
		}
		tests.PanicOnErr(CreatePlatform(&newPlatform))

		err := CreatePlatform(&newPlatform)

		assert.Equal(t, operations.Errors.DataAlreadyExistErr, err)
	})
}

func TestGetPlatforms(t *testing.T) {
	t.Run("Get all platforms", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		platforms := []Platform{
			{
				Name:      "platform 1",
				ShortName: "pm1",
			},
			{
				Name:      "platform 2",
				ShortName: "pm2",
			},
			{
				Name:      "platform 3",
				ShortName: "pm3",
			},
		}
		for i := range platforms {
			tests.PanicOnErr(CreatePlatform(&platforms[i]))
		}

		list, err := GetPlatforms()

		assert.NoError(t, err)
		assert.Len(t, list, len(platforms))
		for _, platform := range list {
			assert.Contains(t, platforms, *platform)
		}
	})
}

func TestGetPlatform(t *testing.T) {
	t.Run("Platform exists - returned", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		newPlatform := Platform{
			Name:      "test platform",
			ShortName: "tp",
		}
		tests.PanicOnErr(CreatePlatform(&newPlatform))

		dbPlatform, err := GetPlatform(newPlatform.Id)

		assert.NoError(t, err)
		assert.Equal(t, newPlatform, *dbPlatform)
	})
	t.Run("Platform does not exist - error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)

		_, err := GetPlatform(tests.GetRandomUuid())

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestUpdatePlatform(t *testing.T) {
	t.Run("Platform exists - updated", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		newPlatform := Platform{
			Name:      "test platform",
			ShortName: "tp",
		}
		tests.PanicOnErr(CreatePlatform(&newPlatform))
		newPlatform.Name = "updated platform"
		newPlatform.ShortName = "up"

		err := UpdatePlatform(&newPlatform)

		assert.NoError(t, err)
		dbPlatform, err := GetPlatform(newPlatform.Id)
		tests.PanicOnErr(err)
		assert.Equal(t, newPlatform, *dbPlatform)
	})
	t.Run("Platform not found - error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		platform := Platform{
			Id: tests.GetRandomUuid(),
		}

		err := UpdatePlatform(&platform)

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
	t.Run("Duplicate platform - returns error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		platform1 := Platform{
			Name:      "test platform1",
			ShortName: "tp",
		}
		tests.PanicOnErr(CreatePlatform(&platform1))
		platform2 := Platform{
			Name:      "test platform2",
			ShortName: "tp2",
		}
		tests.PanicOnErr(CreatePlatform(&platform2))
		platform2.Name = platform1.Name

		err := UpdatePlatform(&platform2)

		assert.Equal(t, operations.Errors.DataAlreadyExistErr, err)
	})
}

func TestDeletePlatform(t *testing.T) {
	t.Run("Platform exists - deleted", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		platform := Platform{
			Name:      "test platform",
			ShortName: "tp",
		}
		tests.PanicOnErr(CreatePlatform(&platform))

		err := DeletePlatform(platform.Id)

		assert.NoError(t, err)
		_, err = GetPlatform(platform.Id)
		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
	t.Run("Platform not found - error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)

		err := DeletePlatform(tests.GetRandomUuid())

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}
