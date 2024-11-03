package controllers

import (
	"github.com/KowalskiPiotr98/ludivault/controllers/dto"
	"github.com/KowalskiPiotr98/ludivault/platforms"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"net/http"
)

func getPlatforms(c *gin.Context) {
	list, err := platforms.GetPlatforms()

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapMany(list, dto.MapPlatformToDto))
}

func getPlatform(c *gin.Context) {
	id, err := parseUuidFromPath(c)
	if err != nil {
		return
	}

	item, err := platforms.GetPlatform(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapPlatformToDto(item))
}

func createPlatform(c *gin.Context) {
	var model dto.PlatformEditDto
	if c.MustBindWith(&model, binding.JSON) != nil {
		return
	}

	mapped := dto.MapPlatformEditDtoToObject(uuid.Nil, &model)
	if err := platforms.CreatePlatform(mapped); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.MapPlatformToDto(mapped))
}

func updatePlatform(c *gin.Context) {
	id, err := parseUuidFromPath(c)
	if err != nil {
		return
	}
	var model dto.PlatformEditDto
	if c.MustBindWith(&model, binding.JSON) != nil {
		return
	}

	mapped := dto.MapPlatformEditDtoToObject(id, &model)
	if err = platforms.UpdatePlatform(mapped); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapPlatformToDto(mapped))
}

func deletePlatform(c *gin.Context) {
	id, err := parseUuidFromPath(c)
	if err != nil {
		return
	}

	err = platforms.DeletePlatform(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
