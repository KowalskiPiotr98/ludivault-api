package controllers

import (
	"errors"
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/KowalskiPiotr98/ludivault/auth"
	"github.com/KowalskiPiotr98/ludivault/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/markbates/goth"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func handleError(c *gin.Context, err error) {
	if errors.Is(err, operations.Errors.DataNotFoundErr) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if errors.Is(err, operations.Errors.DataAlreadyExistErr) {
		c.AbortWithStatus(http.StatusConflict)
		return
	}
	if errors.Is(err, operations.Errors.DataUsedErr) {
		c.AbortWithStatus(http.StatusConflict)
		return
	}
	if errors.Is(err, operations.Errors.RowNumberUnexpectedErr) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusInternalServerError)
}

func parseUuidFromPath(c *gin.Context) (uuid.UUID, error) {
	value := c.Param("id")
	id, err := uuid.Parse(value)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return uuid.Nil, err
	}
	return id, nil
}

func initUserSession(user *goth.User, c *gin.Context) error {
	localUser := users.NewFromProvider(user)
	if err := users.GetOrCreate(localUser); err != nil {
		log.Warnf("Failed to initialise user session: %v", err)
		return err
	}
	if err := auth.StoreUserInSession(c, localUser.Id); err != nil {
		log.Warnf("Failed to initialise user session: %v", err)
		return err
	}
	return nil
}
