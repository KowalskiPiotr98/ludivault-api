package users

import "github.com/KowalskiPiotr98/gotabase"

func scanUser(row gotabase.Row) (*User, error) {
	var user User
	if err := row.Scan(&user.Id, &user.ProviderId, &user.ProviderName, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}
