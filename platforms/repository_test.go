package platforms

import (
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/KowalskiPiotr98/ludivault/internal/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeTestDefaultPlatform() Platform {
	return Platform{
		Name:      "test platform",
		ShortName: "tp",
	}
}

func TestCreatePlatform(t *testing.T) {
	t.Run("New platform created", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		newPlatform := makeTestDefaultPlatform()

		err := CreatePlatform(&newPlatform, tests.MakeTestUserId(getDatabase()))

		assert.NoError(t, err)
		dbRow, err := getDatabase().QueryRow("select id, name, short_name from platforms where id = $1", newPlatform.Id)
		tests.PanicOnErr(err)
		var dbPlatform Platform
		tests.PanicOnErr(dbRow.Scan(&dbPlatform.Id, &dbPlatform.Name, &dbPlatform.ShortName))
		assert.Equal(t, newPlatform, dbPlatform)
	})

	t.Run("Platform already exists", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		newPlatform := makeTestDefaultPlatform()
		userId := tests.MakeTestUserId(getDatabase())
		tests.PanicOnErr(CreatePlatform(&newPlatform, userId))

		err := CreatePlatform(&newPlatform, userId)

		assert.Equal(t, operations.Errors.DataAlreadyExistErr, err)
	})
}

func TestGetPlatforms(t *testing.T) {
	t.Run("Get all platforms", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
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
			tests.PanicOnErr(CreatePlatform(&platforms[i], userId))
		}
		tests.PanicOnErr(CreatePlatform(&Platform{
			Name:      "unauthorised",
			ShortName: "un",
		}, tests.MakeTestUserId(getDatabase())))

		list, err := GetPlatforms(userId)

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
		userId := tests.MakeTestUserId(getDatabase())
		newPlatform := makeTestDefaultPlatform()
		tests.PanicOnErr(CreatePlatform(&newPlatform, userId))

		dbPlatform, err := GetPlatform(newPlatform.Id, userId)

		assert.NoError(t, err)
		assert.Equal(t, newPlatform, *dbPlatform)
	})

	t.Run("Platform does not exist - error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())

		_, err := GetPlatform(tests.GetRandomUuid(), userId)

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})

	t.Run("User not authorised", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		newPlatform := makeTestDefaultPlatform()
		tests.PanicOnErr(CreatePlatform(&newPlatform, userId))

		_, err := GetPlatform(newPlatform.Id, tests.MakeTestUserId(getDatabase()))

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestUpdatePlatform(t *testing.T) {
	t.Run("Platform exists - updated", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		newPlatform := makeTestDefaultPlatform()
		tests.PanicOnErr(CreatePlatform(&newPlatform, userId))
		newPlatform.Name = "updated platform"
		newPlatform.ShortName = "up"

		err := UpdatePlatform(&newPlatform, userId)

		assert.NoError(t, err)
		dbPlatform, err := GetPlatform(newPlatform.Id, userId)
		tests.PanicOnErr(err)
		assert.Equal(t, newPlatform, *dbPlatform)
	})

	t.Run("Platform not found - error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		platform := Platform{
			Id: tests.GetRandomUuid(),
		}

		err := UpdatePlatform(&platform, userId)

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})

	t.Run("Duplicate platform - returns error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		platform1 := Platform{
			Name:      "test platform1",
			ShortName: "tp",
		}
		tests.PanicOnErr(CreatePlatform(&platform1, userId))
		platform2 := Platform{
			Name:      "test platform2",
			ShortName: "tp2",
		}
		tests.PanicOnErr(CreatePlatform(&platform2, userId))
		platform2.Name = platform1.Name

		err := UpdatePlatform(&platform2, userId)

		assert.Equal(t, operations.Errors.DataAlreadyExistErr, err)
	})

	t.Run("User not authorised", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())
		newPlatform := makeTestDefaultPlatform()
		tests.PanicOnErr(CreatePlatform(&newPlatform, userId))
		newPlatform.Name = "updated platform"
		newPlatform.ShortName = "up"

		err := UpdatePlatform(&newPlatform, tests.MakeTestUserId(getDatabase()))

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}

func TestDeletePlatform(t *testing.T) {
	t.Run("Platform exists - deleted", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		platform := makeTestDefaultPlatform()
		userId := tests.MakeTestUserId(getDatabase())
		tests.PanicOnErr(CreatePlatform(&platform, userId))

		err := DeletePlatform(platform.Id, userId)

		assert.NoError(t, err)
		_, err = GetPlatform(platform.Id, userId)
		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})

	t.Run("Platform not found - error", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		userId := tests.MakeTestUserId(getDatabase())

		err := DeletePlatform(tests.GetRandomUuid(), userId)

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})

	t.Run("User not authorised", func(t *testing.T) {
		tests.GetDatabaseWithCleanup(t)
		platform := makeTestDefaultPlatform()
		userId := tests.MakeTestUserId(getDatabase())
		tests.PanicOnErr(CreatePlatform(&platform, userId))

		err := DeletePlatform(platform.Id, tests.MakeTestUserId(getDatabase()))

		assert.Equal(t, operations.Errors.DataNotFoundErr, err)
	})
}
