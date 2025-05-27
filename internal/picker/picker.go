package picker

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/ptdewey/spotify-tools/internal/types"
)

// TODO: picker option param (fzf vs nvim fzf vs telescope vs etc) and args?
func Fzf(tracks []types.SimplerTrack) (*types.SimplerTrack, error) {
	m := formatTracks(tracks)

	cmd := exec.Command("fzf", "--preview", "")

	var input bytes.Buffer
	for t := range m {
		input.WriteString(t + "\n")
	}
	cmd.Stdin = &input

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	song, exists := m[string(bytes.TrimSpace(out))]
	if !exists {
		return nil, fmt.Errorf("selected song not found")
	}

	return song, nil
}

func formatTracks(tracks []types.SimplerTrack) map[string]*types.SimplerTrack {
	tmap := make(map[string]*types.SimplerTrack, len(tracks))

	for _, track := range tracks {
		s := fmt.Sprintf("%s - %s - %s", track.Name, track.Artists, track.Album)
		tmap[s] = &track
	}

	return tmap
}
