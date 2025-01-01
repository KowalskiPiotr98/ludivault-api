package controllers

import (
	"github.com/KowalskiPiotr98/ludivault/auth"
	"github.com/KowalskiPiotr98/ludivault/controllers/dto"
	"github.com/KowalskiPiotr98/ludivault/playthroughs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"net/http"
)

func getPlaythroughs(c *gin.Context) {
	var query struct {
		GameId uuid.UUID `form:"gameId"`
	}
	if c.MustBindWith(&query, binding.Query) != nil {
		return
	}

	list, err := playthroughs.GetPlaythroughs(query.GameId, auth.GetUserId(c))

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapMany(list, dto.MapPlaythroughToDto))
}

func getPlaythrough(c *gin.Context) {
	id, err := parseUuidFromPath(c)
	if err != nil {
		return
	}

	item, err := playthroughs.GetPlaythrough(id, auth.GetUserId(c))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapPlaythroughToDto(item))
}

func createPlaythrough(c *gin.Context) {
	var model dto.PlaythroughEditDto
	if c.MustBindWith(&model, binding.JSON) != nil {
		return
	}

	mapped := dto.MapPlaythroughEditDtoToObject(uuid.Nil, &model)
	if err := playthroughs.CreatePlaythrough(mapped, auth.GetUserId(c)); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.MapPlaythroughToDto(mapped))
}

func updatePlaythrough(c *gin.Context) {
	id, err := parseUuidFromPath(c)
	if err != nil {
		return
	}
	var model dto.PlaythroughEditDto
	if c.MustBindWith(&model, binding.JSON) != nil {
		return
	}

	mapped := dto.MapPlaythroughEditDtoToObject(id, &model)
	if err = playthroughs.UpdatePlaythrough(mapped, auth.GetUserId(c)); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapPlaythroughToDto(mapped))
}

func deletePlaythrough(c *gin.Context) {
	id, err := parseUuidFromPath(c)
	if err != nil {
		return
	}

	err = playthroughs.DeletePlaythrough(id, auth.GetUserId(c))
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
