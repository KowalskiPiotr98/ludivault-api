package users

import (
	"github.com/google/uuid"
	"github.com/markbates/goth"
)

type User struct {
	Id           uuid.UUID
	ProviderId   string
	ProviderName string
	Email        string
}

func (u *User) SetId(id uuid.UUID) {
	u.Id = id
}

func NewFromProvider(providerUser *goth.User) *User {
	return &User{
		ProviderId:   providerUser.UserID,
		ProviderName: providerUser.Provider,
		Email:        providerUser.Email,
	}
}
