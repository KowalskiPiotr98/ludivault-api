package games

import "github.com/KowalskiPiotr98/gotabase"

var (
	getDatabase = func() gotabase.Connector { return gotabase.GetConnection() }
)
