package controllers

import (
	"github.com/KowalskiPiotr98/ludivault/controllers/dto"
	"github.com/KowalskiPiotr98/ludivault/games"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"net/http"
)

func getGames(c *gin.Context) {
	model := struct {
		Limit  int `form:"limit" binding:"min=1,max=100"`
		Offset int `form:"offset" binding:"min=0"`
	}{
		Limit:  20,
		Offset: 0,
	}
	if err := c.MustBindWith(&model, binding.Query); err != nil {
		return
	}

	list, err := games.GetGames(model.Offset, model.Limit)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapMany(list, dto.MapGameToDto))
}

func getGame(c *gin.Context) {
	id, err := parseUuidFromPath(c)
	if err != nil {
		return
	}

	item, err := games.GetGame(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapGameToDto(item))
}

func createGame(c *gin.Context) {
	var model dto.GameEditDto
	if c.MustBindWith(&model, binding.JSON) != nil {
		return
	}

	mapped := dto.MapGameEditDtoToObject(uuid.Nil, &model)
	if err := games.CreateGame(mapped); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.MapGameToDto(mapped))
}

func updateGame(c *gin.Context) {
	id, err := parseUuidFromPath(c)
	if err != nil {
		return
	}
	var model dto.GameEditDto
	if c.MustBindWith(&model, binding.JSON) != nil {
		return
	}

	mapped := dto.MapGameEditDtoToObject(id, &model)
	if err = games.UpdateGame(mapped); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapGameToDto(mapped))
}

func deleteGame(c *gin.Context) {
	id, err := parseUuidFromPath(c)
	if err != nil {
		return
	}

	err = games.DeleteGame(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
