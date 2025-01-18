package db

import (
	"fmt"

	"github.com/Udehlee/Collab-playlist/models"
)

type PgDB struct {
	Db *PgConn
}

func NewPgDB(db *PgConn) *PgDB {
	return &PgDB{
		Db: db,
	}

}

func (p *PgDB) SaveUser(user models.User) error {
	query := ` INSERT INTO users (spotify_id, display_name, email, access_token, refresh_token, token_expiry,logged_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (spotify_id)
	DO UPDATE SET
		display_name = $2,
		email = $3,
		spotify_token = $4,
		spotify_refresh_token = $5,
		token_expiry = $6,
		logged_at = $7;
`
	_, err := p.Db.Conn.Exec(query, user.SpotifyID, user.DisplayName, user.Email, user.AccessToken, user.RefreshToken, user.TokenExpiry, user.LoggedAt)
	if err != nil {
		return fmt.Errorf("could not save user: %v", err)
	}
	return nil
}
