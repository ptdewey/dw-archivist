package api

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

// TODO: create if playlist doesn't exist?
// - will require name search rather than lookup by id

func CopyPlaylist(client *spotify.Client, sourceID string, targetID string) error {
	ctx := context.Background()

	tracks, err := getPlaylistTracks(ctx, client, spotify.ID(sourceID))
	if err != nil {
		return err
	}

	var trackIDs []spotify.ID
	for _, item := range tracks {
		if item.Track.Track != nil {
			trackIDs = append(trackIDs, item.Track.Track.ID)
		}
	}

	const batchSize = 100
	for i := 0; i < len(trackIDs); i += batchSize {
		end := min(i+batchSize, len(trackIDs))
		_, err := client.AddTracksToPlaylist(ctx, spotify.ID(targetID), trackIDs[i:end]...)
		if err != nil {
			return err
		}
	}

	return nil
}

func getPlaylistTracks(ctx context.Context, client *spotify.Client, playlistID spotify.ID) ([]spotify.PlaylistItem, error) {
	var tracks []spotify.PlaylistItem
	offset := 0
	limit := 100
	for {
		page, err := client.GetPlaylistItems(ctx, playlistID, spotify.Offset(offset), spotify.Limit(limit))
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, page.Items...)

		if len(page.Items) < limit {
			break
		}
		offset += limit
	}

	return tracks, nil
}

func GetUserPlaylists(ctx context.Context, client *spotify.Client, userID string) ([]spotify.SimplePlaylist, error) {
	var playlists []spotify.SimplePlaylist
	offset := 0
	limit := 50
	for {
		page, err := client.GetPlaylistsForUser(ctx, userID, spotify.Offset(offset), spotify.Limit(50))
		if err != nil {
			return nil, err
		}

		playlists = append(playlists, page.Playlists...)

		if len(page.Playlists) < limit {
			break
		}

		offset += limit
	}

	return playlists, nil
}
