package repository

import (
	"context"
	"time"

	"github.com/NapasinNgam/spotify-diary/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HistoryRepository struct {
	db *pgxpool.Pool
}

func NewHistoryRepository(db *pgxpool.Pool) *HistoryRepository {
	return &HistoryRepository{db: db}
}

func (r *HistoryRepository) InsertHistory(ctx context.Context, h *model.ListeningHistory) error {
	query := `
		INSERT INTO listening_history 
			(user_id, track_id, track_name, artist_id, artist_name, album_name, album_cover_url, preview_url, duration_ms, played_at, played_date, played_month, genre_category)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (user_id, played_at, track_id) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query,
		h.UserID, h.TrackID, h.TrackName, h.ArtistID, h.ArtistName,
		h.AlbumName, h.AlbumCoverURL, h.PreviewURL, h.DurationMs,
		h.PlayedAt, h.PlayedDate, h.PlayedMonth, h.GenreCategory,
	)
	return err
}

func (r *HistoryRepository) GetSyncCursor(ctx context.Context, userID int) (*model.SyncCursor, error) {
	query := `SELECT id, user_id, cursor_type, last_cursor_ms, last_sync_at, records_fetched, total_synced FROM sync_cursors WHERE user_id = $1 AND cursor_type = 'recently_played'`

	cursor := &model.SyncCursor{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&cursor.ID, &cursor.UserID, &cursor.CursorType,
		&cursor.LastCursorMs, &cursor.LastSyncAt,
		&cursor.RecordsFetched, &cursor.TotalSynced,
	)
	return cursor, err
}

func (r *HistoryRepository) UpsertSyncCursor(ctx context.Context, userID int, lastCursorMs int64, recordsFetched int) error {
	query := `
		INSERT INTO sync_cursors (user_id, cursor_type, last_cursor_ms, last_sync_at, records_fetched, total_synced)
		VALUES ($1, 'recently_played', $2, $3, $4, $4)
		ON CONFLICT (user_id, cursor_type) DO UPDATE SET
			last_cursor_ms = EXCLUDED.last_cursor_ms,
			last_sync_at = EXCLUDED.last_sync_at,
			records_fetched = EXCLUDED.records_fetched,
			total_synced = sync_cursors.total_synced + EXCLUDED.records_fetched
	`

	_, err := r.db.Exec(ctx, query, userID, lastCursorMs, time.Now(), recordsFetched)
	return err
}

type DailyStats struct {
	Date            string `json:"date"`
	TotalTracks     int    `json:"total_tracks"`
	UniqueTracks    int    `json:"unique_tracks"`
	UniqueArtists   int    `json:"unique_artists"`
	TotalDurationMs int64  `json:"total_duration_ms"`
}

func (r *HistoryRepository) GetDailyStats(ctx context.Context, userID int, dateRef string) (*DailyStats, error) {
	// dateRef: "yesterday" or "YYYY-MM-DD"
	dateExpr := "CURRENT_DATE - 1"
	if dateRef != "yesterday" {
		dateExpr = "'" + dateRef + "'::date"
	}

	query := `
		SELECT 
			played_date,
			COUNT(*) as total_tracks,
			COUNT(DISTINCT track_id) as unique_tracks,
			COUNT(DISTINCT artist_id) as unique_artists,
			COALESCE(SUM(duration_ms), 0) as total_duration_ms
		FROM listening_history
		WHERE user_id = $1 AND played_date = ` + dateExpr + `
		GROUP BY played_date
	`

	stats := &DailyStats{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&stats.Date, &stats.TotalTracks, &stats.UniqueTracks,
		&stats.UniqueArtists, &stats.TotalDurationMs,
	)
	return stats, err
}

type GenreTopTrack struct {
	Genre         string `json:"genre"`
	TrackID       string `json:"track_id"`
	TrackName     string `json:"track_name"`
	ArtistName    string `json:"artist_name"`
	AlbumCoverURL string `json:"album_cover_url"`
	PlayCount     int    `json:"play_count"`
}

func (r *HistoryRepository) GetTopByGenre(ctx context.Context, userID int, dateRef string, limit int) (map[string][]GenreTopTrack, error) {
	dateExpr := "CURRENT_DATE - 1"
	if dateRef != "yesterday" {
		dateExpr = "'" + dateRef + "'::date"
	}

	query := `
		SELECT genre_category, track_id, track_name, artist_name, album_cover_url, COUNT(*) as play_count
		FROM listening_history
		WHERE user_id = $1 AND played_date = ` + dateExpr + `
		GROUP BY genre_category, track_id, track_name, artist_name, album_cover_url
		ORDER BY genre_category, play_count DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]GenreTopTrack)
	for rows.Next() {
		var t GenreTopTrack
		err := rows.Scan(&t.Genre, &t.TrackID, &t.TrackName, &t.ArtistName, &t.AlbumCoverURL, &t.PlayCount)
		if err != nil {
			return nil, err
		}
		if len(result[t.Genre]) < limit {
			result[t.Genre] = append(result[t.Genre], t)
		}
	}

	return result, nil
}
