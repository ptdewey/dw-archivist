package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ptdewey/spotify-tools/internal/api"
	"github.com/ptdewey/spotify-tools/internal/cache"
	"github.com/ptdewey/spotify-tools/internal/picker"
)

var (
	dbPath string = "./data/tracks.db"

	playlistsFile string
	refreshCache  bool
	cacheMode     cache.Mode
)

func parsePlaylists(playlistsFile string) []string {
	data, err := os.ReadFile(playlistsFile)
	if err != nil {
		return nil
	}

	lines := strings.Split(string(data), "\n")
	var playlistNames []string
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name != "" {
			playlistNames = append(playlistNames, name)
		}
	}

	return playlistNames
}

func main() {
	cache.InitDB(dbPath)

	flag.StringVar(&playlistsFile, "playlists-file", "playlists.txt", "--playlists-file=\"playlists.txt\"")
	flag.BoolVar(&refreshCache, "recache", false, "--refresh-cache")
	flag.Var(&cacheMode, "cache-mode", "--cache-mode")
	flag.Parse()

	ctx := context.Background()
	client, user := api.Authorize("token.json")

	// TODO: Time-based re-cache (for sqlite)
	if refreshCache {
		if err := cache.Clear(); err != nil {
			fmt.Println(err)
			return
		}

		playlistNames := parsePlaylists(playlistsFile)
		if err := api.SavePlaylistsByName(ctx, client, user, playlistNames, cacheMode); err != nil {
			fmt.Println(err)
			return
		}
	}

	tracks, err := cache.ReadCachedTracks(cacheMode)
	if err != nil {
		fmt.Println(err)
		return
	}

	pickedTrack, err := picker.Fzf(tracks)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := api.QueueSong(ctx, client, pickedTrack); err != nil {
		fmt.Println(err)
		return
	}
}
