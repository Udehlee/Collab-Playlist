package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Udehlee/Collab-playlist/models"
)

func GetUserProfile(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.spotify.com/v1/me")
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error response from Spotify: %s", resp.Status)
	}

	var User models.User

	if err := json.NewDecoder(resp.Body).Decode(&User); err != nil {
		return "", fmt.Errorf("failed to decode user info: %w", err)
	}

	return User.Name, nil
}
