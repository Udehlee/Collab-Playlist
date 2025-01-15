package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Udehlee/Collab-playlist/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

type OAuth struct {
	config *oauth2.Config
	state  string
}

func NewOAuth() *OAuth {
	return &OAuth{
		config: &oauth2.Config{
			ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"playlist-modify-public", "playlist-modify-private"},
			Endpoint:     spotify.Endpoint,
		},
		state: utils.GenerateRandomState(),
	}
}

func (h *OAuth) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login with spotify"))

}

func (o *OAuth) LoginWithSpotify(w http.ResponseWriter, r *http.Request) {
	url := o.config.AuthCodeURL(o.state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (o *OAuth) HandleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	if state != o.state {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	token, err := o.config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := o.config.Client(context.Background(), token)
	user, err := GetUserProfile(client)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		log.Println("User info error:", err)
		return
	}

	fmt.Fprintf(w, "Logged in as User: %s\n", user)

}
