package auth

import "errors"

var (
	UserIdNotStored             = errors.New("user id was not stored in user session")
	NoProvidersSet              = errors.New("no providers were set during initialisation")
	ProvidersAlreadyInitialised = errors.New("providers were already initialised")
)
