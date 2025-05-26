package web

import (
	"log"
	"net/http"
	"text/template"

	"github.com/zmb3/spotify/v2"
)

var playlists []spotify.SimplePlaylist

// REFACTOR: make this actually good
func InitServer(p []spotify.SimplePlaylist) {
	playlists = p

	mux := http.NewServeMux()
	mux.HandleFunc("GET /playlists", HandleServePlaylistsPage)

	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Println(err)
		return
	}
}

func HandleServePlaylistsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/playlists.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = tmpl.Execute(w, playlists)
	if err != nil {
		http.Error(w, "Could not render template", http.StatusInternalServerError)
		log.Println(err)
	}
}
