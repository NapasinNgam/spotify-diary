CREATE TABLE IF NOT EXISTS listening_history (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id        VARCHAR(62) NOT NULL,
    track_name      VARCHAR(500) NOT NULL,
    artist_id       VARCHAR(62) NOT NULL,
    artist_name     VARCHAR(500) NOT NULL,
    album_name      VARCHAR(500),
    album_cover_url VARCHAR(1000),
    preview_url     VARCHAR(1000),
    duration_ms     INTEGER NOT NULL,
    played_at       TIMESTAMP WITH TIME ZONE NOT NULL,
    played_date     DATE NOT NULL,
    played_month    VARCHAR(7) NOT NULL,
    genre_category  VARCHAR(20) DEFAULT 'other',

    UNIQUE(user_id, played_at, track_id)
);

CREATE INDEX idx_history_user_date ON listening_history(user_id, played_date);
CREATE INDEX idx_history_user_month ON listening_history(user_id, played_month);
CREATE INDEX idx_history_user_genre ON listening_history(user_id, genre_category, played_date);

-- Sync cursor table
CREATE TABLE IF NOT EXISTS sync_cursors (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cursor_type     VARCHAR(50) NOT NULL DEFAULT 'recently_played',
    last_cursor_ms  BIGINT NOT NULL DEFAULT 0,
    last_sync_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    records_fetched INTEGER DEFAULT 0,
    total_synced    INTEGER DEFAULT 0,

    UNIQUE(user_id, cursor_type)
);
