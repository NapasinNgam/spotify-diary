package model

import "time"

type DailyLog struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	LogDate       string    `json:"log_date"` // YYYY-MM-DD
	TrackID       string    `json:"track_id"`
	TrackName     string    `json:"track_name"`
	ArtistName    string    `json:"artist_name"`
	AlbumName     string    `json:"album_name"`
	AlbumCoverURL string    `json:"album_cover_url"`
	PreviewURL    string    `json:"preview_url"`
	CreatedAt     time.Time `json:"created_at"`
}

type DailySummary struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	SummaryDate     string    `json:"summary_date"`
	TotalTracks     int       `json:"total_tracks"`
	UniqueTracks    int       `json:"unique_tracks"`
	UniqueArtists   int       `json:"unique_artists"`
	TotalDurationMs int64     `json:"total_duration_ms"`
	TopTrackID      string    `json:"top_track_id"`
	TopTrackName    string    `json:"top_track_name"`
	TopTrackCount   int       `json:"top_track_count"`
	GeneratedAt     time.Time `json:"generated_at"`
}
