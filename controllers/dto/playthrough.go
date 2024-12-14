package dto

import (
	"github.com/KowalskiPiotr98/ludivault/playthroughs"
	"github.com/google/uuid"
	"time"
)

type PlaythroughDto struct {
	Id        uuid.UUID  `json:"id"`
	GameId    uuid.UUID  `json:"gameId"`
	StartDate time.Time  `json:"startDate"`
	EndDate   *time.Time `json:"endDate,omitempty"`
	Status    int        `json:"status"`
	Runtime   *int       `json:"runtime,omitempty"`
}

func MapPlaythroughToDto(playthrough *playthroughs.Playthrough) *PlaythroughDto {
	return &PlaythroughDto{
		Id:        playthrough.Id,
		GameId:    playthrough.GameId,
		StartDate: playthrough.StartDate,
		EndDate:   makePointerFromNullTime(playthrough.EndDate),
		Status:    int(playthrough.Status),
		Runtime:   makePointerFromNullInt(playthrough.Runtime),
	}
}

type PlaythroughEditDto struct {
	GameId    uuid.UUID  `json:"gameId" binding:"required"`
	StartDate time.Time  `json:"startDate" binding:"required"`
	EndDate   *time.Time `json:"endDate"`
	Status    int        `json:"status" binding:"min=0,max=4"`
	Runtime   int        `json:"runtime" binding:"min=0"`
}

func MapPlaythroughEditDtoToObject(id uuid.UUID, playthrough *PlaythroughEditDto) *playthroughs.Playthrough {
	return &playthroughs.Playthrough{
		Id:        id,
		GameId:    playthrough.GameId,
		StartDate: playthrough.StartDate,
		EndDate:   makeNullTimeFromPointer(playthrough.EndDate),
		Status:    playthroughs.PlaythroughStatus(playthrough.Status),
		Runtime:   makeNullInt(playthrough.Runtime),
	}
}
