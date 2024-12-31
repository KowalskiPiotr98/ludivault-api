package users

import (
	"github.com/KowalskiPiotr98/gotabase/operations"
	"github.com/google/uuid"
)

// GetOrCreate sets the id of the provider User object.
// If the User does not yet exist in the database, it will be created.
//
// If the user does exist, the email address will be updated to match the one provided from identity provider.
func GetOrCreate(user *User) error {
	query := `insert into users (provider_id, provider_name, email) values ($1, $2, $3) on conflict (provider_id, provider_name) do update set email = $3 returning id`
	return operations.CreateRowWithId(getDatabase(), user, query, user.ProviderId, user.ProviderName, user.Email)
}

// GetById returns all User data by the user id.
func GetById(id uuid.UUID) (*User, error) {
	query := `select id, provider_id, provider_name, email from users where id = $1`
	return operations.QueryRow[User](getDatabase(), scanUser, query, id)
}
