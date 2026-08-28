-- Monthly summaries
CREATE TABLE IF NOT EXISTS monthly_summaries (
    id                SERIAL PRIMARY KEY,
    user_id           INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    summary_month     VARCHAR(7) NOT NULL,
    total_plays       INTEGER DEFAULT 0,
    unique_tracks     INTEGER DEFAULT 0,
    unique_artists    INTEGER DEFAULT 0,
    total_duration_ms BIGINT DEFAULT 0,
    top_tracks        JSONB DEFAULT '[]',
    top_artists       JSONB DEFAULT '[]',
    genre_breakdown   JSONB DEFAULT '{}',
    generated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(user_id, summary_month)
);

-- Half-year recaps
CREATE TABLE IF NOT EXISTS half_year_recaps (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    period          VARCHAR(10) NOT NULL,
    rank            INTEGER NOT NULL CHECK (rank BETWEEN 1 AND 8),
    track_id        VARCHAR(62) NOT NULL,
    track_name      VARCHAR(500) NOT NULL,
    artist_name     VARCHAR(500) NOT NULL,
    album_cover_url VARCHAR(1000),
    preview_url     VARCHAR(1000),
    description     TEXT DEFAULT '',
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(user_id, period, rank)
);

-- Genre playlist configuration
CREATE TABLE IF NOT EXISTS genre_playlists (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    genre_category  VARCHAR(20) NOT NULL,
    playlist_id     VARCHAR(62) NOT NULL,
    playlist_name   VARCHAR(500),

    UNIQUE(user_id, genre_category)
);
