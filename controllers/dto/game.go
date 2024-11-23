package dto

import (
	"github.com/KowalskiPiotr98/ludivault/games"
	"github.com/google/uuid"
	"time"
)

type GameDto struct {
	Id          uuid.UUID  `json:"id"`
	PlatformId  uuid.UUID  `json:"platformId"`
	Title       string     `json:"title"`
	Owned       bool       `json:"owned"`
	ReleaseDate *time.Time `json:"releaseDate,omitempty"`
	Released    bool       `json:"released"`
}

func MapGameToDto(game *games.Game) *GameDto {
	return &GameDto{
		Id:          game.Id,
		PlatformId:  game.PlatformId,
		Title:       game.Title,
		Owned:       game.Owned,
		ReleaseDate: makePointerFromNullTime(game.ReleaseDate),
		Released:    game.Released,
	}
}

type GameEditDto struct {
	PlatformId  uuid.UUID  `json:"platformId" binding:"required"`
	Title       string     `json:"title" binding:"required,max=500"`
	Owned       bool       `json:"owned"`
	ReleaseDate *time.Time `json:"releaseDate"`
	Released    bool       `json:"released"`
}

func MapGameEditDtoToObject(id uuid.UUID, game *GameEditDto) *games.Game {
	return &games.Game{
		Id:          id,
		PlatformId:  game.PlatformId,
		Title:       game.Title,
		Owned:       game.Owned,
		ReleaseDate: makeNullTimeFromPointer(game.ReleaseDate),
		Released:    game.Released,
	}
}
