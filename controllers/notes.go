package controllers

import (
	"github.com/KowalskiPiotr98/ludivault/auth"
	"github.com/KowalskiPiotr98/ludivault/controllers/dto"
	"github.com/KowalskiPiotr98/ludivault/notes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func getNoteTitles(c *gin.Context) {
	gameId, err := parseUuidFromPath(c)
	if err != nil {
		return
	}

	items, err := notes.GetNoteTitles(gameId, auth.GetUserId(c))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapMany(items, dto.MapNoteToDto))
}

func getNote(c *gin.Context) {
	noteId, err := parseUuidFromPath(c)
	if err != nil {
		return
	}

	note, err := notes.GetNote(noteId, auth.GetUserId(c))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapNoteToDto(note))
}

func createNote(c *gin.Context) {
	var model dto.ModifyNoteDto
	if c.MustBindWith(&model, binding.JSON) != nil {
		return
	}

	mapped := dto.MapCreateDtoToNote(&model)
	if err := notes.CreateNote(mapped, auth.GetUserId(c)); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.MapNoteToDto(mapped))
}

func updateNote(c *gin.Context) {
	var model dto.ModifyNoteDto
	if c.MustBindWith(&model, binding.JSON) != nil {
		return
	}
	id, err := parseUuidFromPath(c)
	if err != nil {
		return
	}

	mapped := dto.MapModifyDtoToNote(id, &model)
	if err := notes.UpdateNote(mapped, auth.GetUserId(c)); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapNoteToDto(mapped))
}

func deleteNote(c *gin.Context) {
	id, err := parseUuidFromPath(c)
	if err != nil {
		return
	}

	err = notes.DeleteNote(id, auth.GetUserId(c))
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
