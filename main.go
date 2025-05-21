package main

import (
	"flag"
	"log"
	"time"

	"github.com/ptdewey/dw-archivist/internal/auth"
	"github.com/ptdewey/dw-archivist/internal/playlists"
	"github.com/robfig/cron/v3"
)

var (
	debug            bool
	targetID         string
	discoverWeeklyID string
	tz               string
)

func main() {
	// TODO: add options for using playlist names rather than IDs (and/or URIs)
	flag.BoolVar(&debug, "debug", false, "--debug")
	flag.StringVar(&targetID, "target-id", "", "--targetID <target-playlist-id>")
	flag.StringVar(&discoverWeeklyID, "discover-weekly-id", "", "--discover-weekly-id <source-playlist-id>")
	flag.StringVar(&tz, "tz", "America/New_York", "--tz <go-timezone-string>")
	flag.Parse()

	client, _ := auth.Authorize()

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatalf("failed to load EST location: %v", err)
	}

	if debug {
		err := playlists.CopyPlaylist(client, discoverWeeklyID, targetID)
		if err != nil {
			log.Printf("Error copying playlist: %v", err)
		} else {
			log.Println("Successfully copied Discover Weekly to target playlist")
		}
	} else {
		c := cron.New(cron.WithLocation(loc))
		_, err = c.AddFunc("0 12 * * 1", func() {
			err := playlists.CopyPlaylist(client, discoverWeeklyID, targetID)
			if err != nil {
				log.Printf("Error copying playlist: %v", err)
			} else {
				log.Println("Successfully copied Discover Weekly to playlist with ID: ", targetID)
			}
		})
		if err != nil {
			log.Fatalf("failed to schedule task: %v", err)
		}

		c.Start()
		select {}
	}
}
