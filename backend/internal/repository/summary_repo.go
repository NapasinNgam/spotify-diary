package repository

import (
	"context"

	"github.com/NapasinNgam/spotify-diary/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SummaryRepository struct {
	db *pgxpool.Pool
}

func NewSummaryRepository(db *pgxpool.Pool) *SummaryRepository {
	return &SummaryRepository{db: db}
}

func (r *SummaryRepository) GetMonthlySummary(ctx context.Context, userID int, month string) (*model.MonthlySummary, error) {
	query := `
		SELECT id, user_id, summary_month, total_plays, unique_tracks, unique_artists, 
			total_duration_ms, top_tracks, top_artists, genre_breakdown, generated_at
		FROM monthly_summaries
		WHERE user_id = $1 AND summary_month = $2
	`

	s := &model.MonthlySummary{}
	err := r.db.QueryRow(ctx, query, userID, month).Scan(
		&s.ID, &s.UserID, &s.SummaryMonth, &s.TotalPlays,
		&s.UniqueTracks, &s.UniqueArtists, &s.TotalDurationMs,
		&s.TopTracks, &s.TopArtists, &s.GenreBreakdown, &s.GeneratedAt,
	)

	return s, err
}

func (r *SummaryRepository) ListAvailableMonths(ctx context.Context, userID int) ([]string, error) {
	query := `SELECT summary_month FROM monthly_summaries WHERE user_id = $1 ORDER BY summary_month DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var months []string
	for rows.Next() {
		var m string
		if err := rows.Scan(&m); err != nil {
			return nil, err
		}
		months = append(months, m)
	}

	return months, nil
}

type RecapRepository struct {
	db *pgxpool.Pool
}

func NewRecapRepository(db *pgxpool.Pool) *RecapRepository {
	return &RecapRepository{db: db}
}

type UpsertRecapParams struct {
	UserID        int
	Period        string
	Rank          int
	TrackID       string
	TrackName     string
	ArtistName    string
	AlbumCoverURL string
	PreviewURL    string
	Description   string
}

func (r *RecapRepository) Upsert(ctx context.Context, params UpsertRecapParams) error {
	query := `
		INSERT INTO half_year_recaps (user_id, period, rank, track_id, track_name, artist_name, album_cover_url, preview_url, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id, period, rank) DO UPDATE SET
			track_id = EXCLUDED.track_id,
			track_name = EXCLUDED.track_name,
			artist_name = EXCLUDED.artist_name,
			album_cover_url = EXCLUDED.album_cover_url,
			preview_url = EXCLUDED.preview_url,
			description = EXCLUDED.description,
			updated_at = NOW()
	`

	_, err := r.db.Exec(ctx, query,
		params.UserID, params.Period, params.Rank, params.TrackID,
		params.TrackName, params.ArtistName, params.AlbumCoverURL,
		params.PreviewURL, params.Description,
	)
	return err
}

func (r *RecapRepository) GetByPeriod(ctx context.Context, userID int, period string) ([]model.HalfYearRecap, error) {
	query := `
		SELECT id, user_id, period, rank, track_id, track_name, artist_name, album_cover_url, preview_url, description, created_at, updated_at
		FROM half_year_recaps
		WHERE user_id = $1 AND period = $2
		ORDER BY rank
	`

	rows, err := r.db.Query(ctx, query, userID, period)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recaps []model.HalfYearRecap
	for rows.Next() {
		var rc model.HalfYearRecap
		err := rows.Scan(&rc.ID, &rc.UserID, &rc.Period, &rc.Rank, &rc.TrackID,
			&rc.TrackName, &rc.ArtistName, &rc.AlbumCoverURL, &rc.PreviewURL,
			&rc.Description, &rc.CreatedAt, &rc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		recaps = append(recaps, rc)
	}

	return recaps, nil
}

func (r *RecapRepository) ListPeriods(ctx context.Context, userID int) ([]string, error) {
	query := `SELECT DISTINCT period FROM half_year_recaps WHERE user_id = $1 ORDER BY period DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periods []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		periods = append(periods, p)
	}

	return periods, nil
}
