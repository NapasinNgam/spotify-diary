CREATE TABLE IF NOT EXISTS users (
    id              SERIAL PRIMARY KEY,
    spotify_id      VARCHAR(255) NOT NULL UNIQUE,
    display_name    VARCHAR(255) NOT NULL,
    email           VARCHAR(255),
    avatar_url      VARCHAR(1000),
    access_token    TEXT NOT NULL,
    refresh_token   TEXT NOT NULL,
    token_expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_spotify_id ON users(spotify_id);
