CREATE TABLE IF NOT EXISTS daily_logs (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    log_date        DATE NOT NULL,
    track_id        VARCHAR(62) NOT NULL,
    track_name      VARCHAR(500) NOT NULL,
    artist_name     VARCHAR(500) NOT NULL,
    album_name      VARCHAR(500),
    album_cover_url VARCHAR(1000),
    preview_url     VARCHAR(1000),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(user_id, log_date)
);

CREATE INDEX idx_daily_logs_user_date ON daily_logs(user_id, log_date);

-- Daily summaries (auto-generated)
CREATE TABLE IF NOT EXISTS daily_summaries (
    id                SERIAL PRIMARY KEY,
    user_id           INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    summary_date      DATE NOT NULL,
    total_tracks      INTEGER DEFAULT 0,
    unique_tracks     INTEGER DEFAULT 0,
    unique_artists    INTEGER DEFAULT 0,
    total_duration_ms BIGINT DEFAULT 0,
    top_track_id      VARCHAR(62),
    top_track_name    VARCHAR(500),
    top_track_count   INTEGER DEFAULT 0,
    generated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(user_id, summary_date)
);
