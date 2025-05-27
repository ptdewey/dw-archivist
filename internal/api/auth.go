package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserModifyPlaybackState,
			spotifyauth.ScopePlaylistReadPrivate,
		),
	)
	clientChan = make(chan *spotify.Client)
	state      = "abc123"
)

type storedToken struct {
	Token *oauth2.Token `json:"token"`
}

func Authorize(filename string) (*spotify.Client, *spotify.PrivateUser) {
	ctx := context.Background()

	client, user, err := authenticateWithFile(ctx, filename)
	if err == nil {
		return client, user
	}
	fmt.Println("Error:", err, "\nAttempting authentication with browser.")

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
	client = <-clientChan
	user, err = client.CurrentUser(ctx)
	if err != nil {
		log.Fatal(err)
	}

	tok, err := client.Token()
	if err != nil {
		panic("Failed to extract token from client.")
	}

	if err := storeToken(tok, filename); err != nil {
		fmt.Println("Failed to store token: ", err)
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

func authenticateWithFile(ctx context.Context, filename string) (*spotify.Client, *spotify.PrivateUser, error) {
	tok, err := readToken(filename)
	if err != nil {
		fmt.Println("Failed to authenticate with stored token.")
		return nil, nil, err
	}

	client := spotify.New(auth.Client(ctx, tok))
	user, err := client.CurrentUser(ctx)
	if err != nil {
		fmt.Println("Failed to authenticate user with stored token.")
		return nil, nil, err
	}

	return client, user, nil
}

func readToken(filename string) (*oauth2.Token, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var s storedToken
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	return s.Token, nil
}

func storeToken(tok *oauth2.Token, filename string) error {
	data, err := json.Marshal(storedToken{tok})
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	return nil
}
