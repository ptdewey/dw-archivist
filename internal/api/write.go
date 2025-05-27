package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ptdewey/spotify-tools/internal/cache"
	"github.com/zmb3/spotify/v2"
)

func SavePlaylistsByName(ctx context.Context, client *spotify.Client, user *spotify.PrivateUser, playlistNames []string, cacheMode cache.Mode) error {
	playlists, err := GetUserPlaylists(ctx, client, user.ID)
	if err != nil {
		return err
	}

	// REFACTOR: move some of this to the caching package
	for _, playlistName := range playlistNames {
		playlistID, err := getPlaylistID(playlists, playlistName)
		if err != nil {
			return fmt.Errorf("%v %s", err, playlistName)
		}

		items, err := getPlaylistTracks(ctx, client, playlistID)
		if err != nil {
			return err
		}

		if _, err := os.Stat("data"); err != nil {
			if !os.IsNotExist(err) {
				return err
			}

			if err := os.Mkdir("data", 0755); err != nil {
				return err
			}
		}

		switch cacheMode {
		case cache.JSON:
			filename := filepath.Join("data", playlistName+".json")
			tracks := simplifyTracks(items)

			data, err := json.Marshal(tracks)
			if err != nil {
				return err
			}

			if err := os.WriteFile(filename, data, 0644); err != nil {
				return err
			}
		case cache.SQLite:
			// TODO: allow storing playlists separately (or just use a field for that)
			tracks := simplifyTracks(items)
			return cache.InsertTracks(tracks)
		}
	}

	return nil
}
