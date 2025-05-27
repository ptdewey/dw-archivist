package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ptdewey/spotify-tools/internal/api"
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

	fmt.Println("target-id", targetID)
	fmt.Println("source-id", discoverWeeklyID)

	client, user := api.Authorize("token.json")
	_ = user // FIX: replace with website serving stuff

	// p, err := playlists.GetUserPlaylists(context.Background(), client, user.ID)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// fmt.Printf("Found %d Playlists.\n", len(p))
	//
	// web.InitServer(p) // REFACTOR: probably move either this or the other code somewhere else

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatalf("failed to load EST location: %v", err)
	}

	if debug {
		err := api.CopyPlaylist(client, discoverWeeklyID, targetID)
		if err != nil {
			log.Printf("Error copying playlist: %v", err)
		} else {
			log.Println("Successfully copied Discover Weekly to target playlist")
		}
	} else {
		fmt.Println("Initializing scheduled jobs...")
		c := cron.New(cron.WithLocation(loc))
		_, err = c.AddFunc("0 12 * * 1", func() {
			err := api.CopyPlaylist(client, discoverWeeklyID, targetID)
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
		fmt.Println("Started cron scheduler. Waiting for jobs to execute...")
		select {}
	}
}
