package repository

import (
	"context"
	"time"

	"github.com/NapasinNgam/spotify-diary/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DiaryRepository struct {
	db *pgxpool.Pool
}

func NewDiaryRepository(db *pgxpool.Pool) *DiaryRepository {
	return &DiaryRepository{db: db}
}

type UpsertDiaryParams struct {
	UserID        int
	LogDate       string
	TrackID       string
	TrackName     string
	ArtistName    string
	AlbumName     string
	AlbumCoverURL string
	PreviewURL    string
}

func (r *DiaryRepository) Upsert(ctx context.Context, params UpsertDiaryParams) (*model.DailyLog, error) {
	// Parse date string to time.Time
	logDate, err := time.Parse("2006-01-02", params.LogDate)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO daily_logs (user_id, log_date, track_id, track_name, artist_name, album_name, album_cover_url, preview_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, log_date) DO UPDATE SET
			track_id = EXCLUDED.track_id,
			track_name = EXCLUDED.track_name,
			artist_name = EXCLUDED.artist_name,
			album_name = EXCLUDED.album_name,
			album_cover_url = EXCLUDED.album_cover_url,
			preview_url = EXCLUDED.preview_url
		RETURNING id, user_id, log_date, track_id, track_name, artist_name, album_name, album_cover_url, preview_url, created_at
	`

	entry := &model.DailyLog{}
	err = r.db.QueryRow(ctx, query,
		params.UserID, logDate, params.TrackID, params.TrackName,
		params.ArtistName, params.AlbumName, params.AlbumCoverURL, params.PreviewURL,
	).Scan(
		&entry.ID, &entry.UserID, &entry.LogDate, &entry.TrackID,
		&entry.TrackName, &entry.ArtistName, &entry.AlbumName,
		&entry.AlbumCoverURL, &entry.PreviewURL, &entry.CreatedAt,
	)

	return entry, err
}

func (r *DiaryRepository) GetByMonth(ctx context.Context, userID int, month string) ([]model.DailyLog, error) {
	query := `
		SELECT id, user_id, log_date, track_id, track_name, artist_name, album_name, album_cover_url, preview_url, created_at
		FROM daily_logs
		WHERE user_id = $1 AND to_char(log_date, 'YYYY-MM') = $2
		ORDER BY log_date
	`

	rows, err := r.db.Query(ctx, query, userID, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []model.DailyLog
	for rows.Next() {
		var e model.DailyLog
		err := rows.Scan(&e.ID, &e.UserID, &e.LogDate, &e.TrackID,
			&e.TrackName, &e.ArtistName, &e.AlbumName,
			&e.AlbumCoverURL, &e.PreviewURL, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, nil
}

func (r *DiaryRepository) Delete(ctx context.Context, userID int, date string) error {
	logDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return err
	}
	query := `DELETE FROM daily_logs WHERE user_id = $1 AND log_date = $2`
	_, err = r.db.Exec(ctx, query, userID, logDate)
	return err
}
