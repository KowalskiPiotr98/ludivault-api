package repositories

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

var (
	DataNotFoundErr  = errors.New("requested data was not found in the database")
	DuplicateDataErr = errors.New("this data already exists")
)

func IsNotFoundErr(err error) bool {
	//todo: placeholder
	return false
}

func IsDuplicateDataErr(err error) bool {
	//todo: placeholder
	return false
}

func HandleKnownError(err error) error {
	if IsNotFoundErr(err) {
		return DataNotFoundErr
	}
	if IsDuplicateDataErr(err) {
		return DuplicateDataErr
	}
	log.Warnf("Unexpected database error: %v", err)
	return err
}
