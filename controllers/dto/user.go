package dto

import (
	"github.com/KowalskiPiotr98/ludivault/users"
	"github.com/google/uuid"
)

type UserDto struct {
	Id           uuid.UUID `json:"id"`
	ProviderId   string    `json:"providerId"`
	ProviderName string    `json:"providerName"`
	Email        string    `json:"email"`
}

func MapUserToDto(user *users.User) *UserDto {
	return &UserDto{
		Id:           user.Id,
		ProviderId:   user.ProviderId,
		ProviderName: user.ProviderName,
		Email:        user.Email,
	}
}
