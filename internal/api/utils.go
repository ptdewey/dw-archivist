package api

import (
	"fmt"
	"strings"

	"github.com/ptdewey/spotify-tools/internal/types"
	"github.com/zmb3/spotify/v2"
)

type TrackDetails struct {
	Name   string
	Album  string
	Artist string
}

func getPlaylistID(playlists []spotify.SimplePlaylist, playlistName string) (spotify.ID, error) {
	for _, p := range playlists {
		if p.Name == playlistName {
			return p.ID, nil
		}
	}

	return "", fmt.Errorf("failed to find matching playlist in user library")
}

func getSongID(tracks []types.SimplerTrack, td TrackDetails) (spotify.ID, error) {
	for _, track := range tracks {
		if track.Name == td.Name && (track.Album == td.Album || strings.Contains(track.Artists, td.Artist)) {
			return track.ID, nil
		}
	}

	return "", fmt.Errorf("failed to find matching matching song")
}

func simplifyTracks(items []spotify.PlaylistItem) []types.SimplerTrack {
	tracks := make([]types.SimplerTrack, 0, len(items))

	for _, item := range items {
		t := item.Track.Track

		artists := make([]string, 0, len(t.Artists))
		for _, artist := range t.Artists {
			artists = append(artists, artist.Name)
		}

		out := types.SimplerTrack{
			Album:   t.Album.Name,
			Artists: strings.Join(artists, ", "),
			Name:    t.Name,
			ID:      t.ID,
			URI:     t.URI,
		}

		tracks = append(tracks, out)
	}

	return tracks
}
