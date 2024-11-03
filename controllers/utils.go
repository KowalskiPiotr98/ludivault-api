package controllers

import (
	"errors"
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
