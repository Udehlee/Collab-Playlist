CREATE TABLE users (
    spotify_id VARCHAR(255) PRIMARY KEY,   -- Spotify user ID as primary key
    display_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    token_expiry TIMESTAMP NOT NULL
);
