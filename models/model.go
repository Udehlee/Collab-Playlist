package models

import "time"

type User struct {
	SpotifyID    string    `json:"spotify_id"`
	DisplayName  string    `json:"display_name"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenExpiry  time.Time `json:"token_expiry"`
	LoggedAt     time.Time `json:"logged_at"`
}
