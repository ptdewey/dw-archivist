package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopePlaylistModifyPublic),
	)
	clientChan = make(chan *spotify.Client)
	state      = "abc123"
)

// TODO: store credential token to file, remove if it expires

func Authorize() (*spotify.Client, *spotify.PrivateUser) {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Please log into to Spotify by visiting the following page in your browser:", url)

	// Wait for auth to complete
	client := <-clientChan

	ctx := context.Background()

	user, err := client.CurrentUser(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You are logged in as:", user.ID)

	return client, user
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		fmt.Println("Couldn't get token:", http.StatusForbidden)
		log.Fatal(err)
	}

	if st := r.FormValue("state"); st != state {
		log.Fatalf("State mismatch: %s != %s\n", st, state)
		http.NotFound(w, r)
	}

	client := spotify.New(auth.Client(r.Context(), tok))
	clientChan <- client
	fmt.Println(w, "Login Completed!")
}
