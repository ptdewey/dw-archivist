package cache

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ptdewey/spotify-tools/internal/types"
)

var db *sqlx.DB

func InitDB(filename string) {
	// TODO: create preceding path to filename if it does not exist
	var err error
	db, err = sqlx.Connect("sqlite3", filename)
	if err != nil {
		panic(err)
	}

	var query = `CREATE TABLE Tracks (
		name TEXT NOT NULL,
		artists TEXT,
		album TEXT,
		id TEXT PRIMARY KEY,
		uri TEXT NOT NULL
	);`

	if _, err = db.Exec(query); err != nil {
		return
	}
}

func InsertTracks(tracks []types.SimplerTrack) error {
	query := `INSERT INTO Tracks (
		name,
		artists,
		album,
		id,
		uri
	) VALUES (
		:name,
		:artists,
		:album,
		:id,
		:uri
	);`

	_, err := db.NamedExec(query, tracks)
	return err
}

func Clear() error {
	query := "DELETE FROM Tracks"
	_, err := db.Exec(query)
	return err
}

func readTracksSQL() ([]types.SimplerTrack, error) {
	var tracks []types.SimplerTrack
	if err := db.Select(&tracks, "SELECT * FROM Tracks"); err != nil {
		return nil, err
	}

	return tracks, nil
}
