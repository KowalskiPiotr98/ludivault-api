package auth

import "errors"

var (
	UserIdNotStored = errors.New("user id was not stored in user session")
)
