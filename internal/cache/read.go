package cache

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ptdewey/spotify-tools/internal/types"
)

func ReadCachedTracks(mode Mode) ([]types.SimplerTrack, error) {
	var tracks []types.SimplerTrack

	switch mode {
	case JSON:
		if err := filepath.WalkDir("./data", func(path string, d fs.DirEntry, _ error) error {
			if d.IsDir() || !strings.HasSuffix(path, ".json") {
				return nil
			}

			t, err := readTracksJSON(path)
			if err != nil {
				return err
			}
			tracks = append(tracks, t...)

			return nil
		}); err != nil {
			return nil, err
		}
	case SQLite:
		var err error
		tracks, err = readTracksSQL()
		if err != nil {
			return nil, err
		}
	}

	return tracks, nil
}

func readTracksJSON(filename string) ([]types.SimplerTrack, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var tracks []types.SimplerTrack
	if err := json.Unmarshal(data, &tracks); err != nil {
		return nil, err
	}

	return tracks, err
}
