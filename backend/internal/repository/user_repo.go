package repository

import (
	"context"
	"time"

	"github.com/NapasinNgam/spotify-diary/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

type UpsertUserParams struct {
	SpotifyID      string
	DisplayName    string
	Email          string
	AvatarURL      string
	AccessToken    string
	RefreshToken   string
	TokenExpiresAt time.Time
}

func (r *UserRepository) UpsertUser(ctx context.Context, params UpsertUserParams) (*model.User, error) {
	query := `
		INSERT INTO users (spotify_id, display_name, email, avatar_url, access_token, refresh_token, token_expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (spotify_id) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			email = EXCLUDED.email,
			avatar_url = EXCLUDED.avatar_url,
			access_token = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token,
			token_expires_at = EXCLUDED.token_expires_at,
			updated_at = NOW()
		RETURNING id, spotify_id, display_name, email, avatar_url, created_at, updated_at
	`

	user := &model.User{}
	err := r.db.QueryRow(ctx, query,
		params.SpotifyID,
		params.DisplayName,
		params.Email,
		params.AvatarURL,
		params.AccessToken,
		params.RefreshToken,
		params.TokenExpiresAt,
	).Scan(
		&user.ID,
		&user.SpotifyID,
		&user.DisplayName,
		&user.Email,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `SELECT id, spotify_id, display_name, email, avatar_url, created_at, updated_at FROM users WHERE id = $1`

	user := &model.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.SpotifyID,
		&user.DisplayName,
		&user.Email,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func (r *UserRepository) GetBySpotifyID(ctx context.Context, spotifyID string) (*model.User, error) {
	query := `SELECT id, spotify_id, display_name, email, avatar_url, access_token, refresh_token, token_expires_at, created_at, updated_at FROM users WHERE spotify_id = $1`

	user := &model.User{}
	err := r.db.QueryRow(ctx, query, spotifyID).Scan(
		&user.ID,
		&user.SpotifyID,
		&user.DisplayName,
		&user.Email,
		&user.AvatarURL,
		&user.AccessToken,
		&user.RefreshToken,
		&user.TokenExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func (r *UserRepository) UpdateTokens(ctx context.Context, userID int, accessToken, refreshToken string, expiresAt time.Time) error {
	query := `UPDATE users SET access_token = $2, refresh_token = $3, token_expires_at = $4, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, userID, accessToken, refreshToken, expiresAt)
	return err
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	query := `SELECT id, spotify_id, display_name, email, avatar_url, access_token, refresh_token, token_expires_at FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.SpotifyID, &u.DisplayName, &u.Email, &u.AvatarURL, &u.AccessToken, &u.RefreshToken, &u.TokenExpiresAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
