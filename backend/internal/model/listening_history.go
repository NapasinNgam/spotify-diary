package model

import "time"

type ListeningHistory struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	TrackID       string    `json:"track_id"`
	TrackName     string    `json:"track_name"`
	ArtistID      string    `json:"artist_id"`
	ArtistName    string    `json:"artist_name"`
	AlbumName     string    `json:"album_name"`
	AlbumCoverURL string    `json:"album_cover_url"`
	PreviewURL    string    `json:"preview_url"`
	DurationMs    int       `json:"duration_ms"`
	PlayedAt      time.Time `json:"played_at"`
	PlayedDate    string    `json:"played_date"`  // YYYY-MM-DD
	PlayedMonth   string    `json:"played_month"` // YYYY-MM
	GenreCategory string    `json:"genre_category"`
}

type SyncCursor struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	CursorType     string    `json:"cursor_type"`
	LastCursorMs   int64     `json:"last_cursor_ms"`
	LastSyncAt     time.Time `json:"last_sync_at"`
	RecordsFetched int       `json:"records_fetched"`
	TotalSynced    int       `json:"total_synced"`
}
