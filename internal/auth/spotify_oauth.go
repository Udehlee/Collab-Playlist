package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Udehlee/Collab-playlist/db/db"
	"github.com/Udehlee/Collab-playlist/models"
	"github.com/Udehlee/Collab-playlist/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

type OAuth struct {
	config *oauth2.Config
	state  string
	db     *db.PgDB
}

func NewOAuth(db *db.PgDB) *OAuth {
	return &OAuth{
		config: &oauth2.Config{
			ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RedirectURL:  "http://localhost:8080/callback/spotify",
			Scopes:       []string{"user-read-private", "user-read-email"},
			Endpoint:     spotify.Endpoint,
		},
		state: utils.GenerateRandomState(),
		db:    db,
	}
}

func (o *OAuth) Index(w http.ResponseWriter, r *http.Request) {
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

	if err = o.db.SaveUser(user); err != nil {
		log.Printf("Failed to save user: %v", err)
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Logged in as User: %s\n", user)
	fmt.Printf("logged in a user: %s", user)

}

// Check if the token is expired
// Use the refresh token to get a new access token
// Update the user's token
// Token is still valid, no need to refresh
func (o *OAuth) RefreshToken(user *models.User) (*oauth2.Token, error) {
	if time.Now().After(user.TokenExpiry) {
		tokenSource := o.config.TokenSource(context.Background(), &oauth2.Token{
			AccessToken:  user.AccessToken,
			RefreshToken: user.RefreshToken,
		})

		newToken, err := tokenSource.Token()
		if err != nil {
			return nil, err
		}

		user.AccessToken = newToken.AccessToken
		user.RefreshToken = newToken.RefreshToken
		user.TokenExpiry = time.Now().Add(time.Until(newToken.Expiry))

		err = o.db.SaveUser(*user)
		if err != nil {
			return nil, err
		}

		return newToken, nil
	}

	return &oauth2.Token{
		AccessToken: user.AccessToken,
	}, nil
}
