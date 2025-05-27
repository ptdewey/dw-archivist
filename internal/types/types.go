package types

import "github.com/zmb3/spotify/v2"

type SimplerTrack struct {
	Album   string      `json:"album" db:"album"`
	Artists string      `json:"artists" db:"artists"`
	Name    string      `json:"name" db:"name"`
	ID      spotify.ID  `json:"id" db:"id"`
	URI     spotify.URI `json:"uri" db:"uri"`
}
