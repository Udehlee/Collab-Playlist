package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Udehlee/Collab-playlist/models"
)

func GetUserProfile(client *http.Client) (models.User, error) {
	resp, err := client.Get("https://api.spotify.com/v1/me")
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.User{}, fmt.Errorf("error response from Spotify: %s", resp.Status)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return models.User{}, fmt.Errorf("failed to decode user info: %w", err)
	}

	return user, nil
}
