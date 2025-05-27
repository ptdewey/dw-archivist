package api

import (
	"context"
	"log"

	"github.com/ptdewey/spotify-tools/internal/types"
	"github.com/zmb3/spotify/v2"
)

func QueueSong(ctx context.Context, client *spotify.Client, track *types.SimplerTrack) error {
	if err := client.QueueSong(ctx, track.ID); err != nil {
		return err
	}

	// TODO: pretty print artists array (i.e. join into string with commas and 'and')
	log.Printf("Added '%s' by %s (from '%s') to the queue.\n", track.Name, track.Artists, track.Album)

	return nil
}
