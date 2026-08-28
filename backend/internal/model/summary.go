package model

import "time"

type MonthlySummary struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	SummaryMonth    string    `json:"summary_month"` // YYYY-MM
	TotalPlays      int       `json:"total_plays"`
	UniqueTracks    int       `json:"unique_tracks"`
	UniqueArtists   int       `json:"unique_artists"`
	TotalDurationMs int64     `json:"total_duration_ms"`
	TopTracks       string    `json:"top_tracks"`      // JSON array
	TopArtists      string    `json:"top_artists"`     // JSON array
	GenreBreakdown  string    `json:"genre_breakdown"` // JSON object
	PeriodStart     *time.Time `json:"period_start"`   // DATE
	PeriodEnd       *time.Time `json:"period_end"`     // DATE
	Source          string    `json:"source"`          // "spotify_top_tracks" or "listening_history"
	GeneratedAt     time.Time `json:"generated_at"`
}

type HalfYearRecap struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Period        string    `json:"period"` // e.g. "2026-H1", "2026-H2"
	Rank          int       `json:"rank"`   // 1-5 for top, 6-8 for honorable mentions
	TrackID       string    `json:"track_id"`
	TrackName     string    `json:"track_name"`
	ArtistName    string    `json:"artist_name"`
	AlbumCoverURL string    `json:"album_cover_url"`
	PreviewURL    string    `json:"preview_url"`
	Description   string    `json:"description"` // user's personal note
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GenrePlaylist struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	GenreCategory string `json:"genre_category"`
	PlaylistID    string `json:"playlist_id"`
	PlaylistName  string `json:"playlist_name"`
}

// TopTrackItem represents a track in the top_tracks JSON array
type TopTrackItem struct {
	Rank          int    `json:"rank"`
	TrackID       string `json:"track_id"`
	TrackName     string `json:"track_name"`
	ArtistName    string `json:"artist_name"`
	AlbumCoverURL string `json:"album_cover_url"`
	PreviewURL    string `json:"preview_url"`
	PlayCount     int    `json:"play_count"` // 0 if from spotify_top_tracks (no play count available)
}
