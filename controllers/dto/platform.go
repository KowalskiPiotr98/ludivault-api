package dto

import (
	"github.com/KowalskiPiotr98/ludivault/platforms"
	"github.com/google/uuid"
)

type PlatformDto struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ShortName string    `json:"shortName"`
}

func MapPlatformToDto(platform *platforms.Platform) *PlatformDto {
	return &PlatformDto{
		Id:        platform.Id,
		Name:      platform.Name,
		ShortName: platform.ShortName,
	}
}

type PlatformEditDto struct {
	Name      string `json:"name" binding:"required,max=200"`
	ShortName string `json:"shortName" binding:"required,max=5"`
}

func MapPlatformEditDtoToObject(id uuid.UUID, platform *PlatformEditDto) *platforms.Platform {
	return &platforms.Platform{
		Id:        id,
		Name:      platform.Name,
		ShortName: platform.ShortName,
	}
}
