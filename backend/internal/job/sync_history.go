package job

import (
	"context"
	"time"

	"github.com/NapasinNgam/spotify-diary/internal/model"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func (s *Scheduler) SyncAllUsers() {
	ctx := context.Background()

	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		s.logger.Error("Failed to get users for sync", zap.Error(err))
		return
	}

	for _, user := range users {
		s.syncUser(ctx, user)
	}
}

func (s *Scheduler) syncUser(ctx context.Context, user model.User) {
	logger := s.logger.With(zap.Int("user_id", user.ID), zap.String("spotify_id", user.SpotifyID))

	// Create authenticator for token refresh
	authenticator := spotifyauth.New(
		spotifyauth.WithClientID(s.cfg.SpotifyClientID),
		spotifyauth.WithClientSecret(s.cfg.SpotifyClientSecret),
	)

	token := &oauth2.Token{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		Expiry:       user.TokenExpiresAt,
	}

	// Create HTTP client with auto-refresh token source
	httpClient := authenticator.Client(ctx, token)
	client := spotify.New(httpClient)

	// Get cursor
	cursor, err := s.historyRepo.GetSyncCursor(ctx, user.ID)
	var afterMs int64
	if err == nil {
		afterMs = cursor.LastCursorMs
	}

	// Fetch recently played
	opts := &spotify.RecentlyPlayedOptions{
		Limit: 50,
	}
	if afterMs > 0 {
		afterTime := time.UnixMilli(afterMs)
		opts.AfterEpochMs = afterTime.UnixMilli()
	}

	results, err := client.PlayerRecentlyPlayedOpt(ctx, opts)
	if err != nil {
		logger.Error("Failed to fetch recently played", zap.Error(err))
		return
	}

	if len(results) == 0 {
		logger.Debug("No new tracks to sync")
		return
	}

	// Insert tracks
	var newestCursorMs int64
	inserted := 0

	for _, item := range results {
		playedAt := item.PlayedAt

		record := &model.ListeningHistory{
			UserID:        user.ID,
			TrackID:       string(item.Track.ID),
			TrackName:     item.Track.Name,
			ArtistID:      string(item.Track.Artists[0].ID),
			ArtistName:    item.Track.Artists[0].Name,
			AlbumName:     item.Track.Album.Name,
			AlbumCoverURL: "",
			PreviewURL:    "",
			DurationMs:    int(item.Track.Duration),
			PlayedAt:      playedAt,
			PlayedDate:    playedAt.Format("2006-01-02"),
			PlayedMonth:   playedAt.Format("2006-01"),
			GenreCategory: "other", // TODO: resolve genre from artist
		}

		// Get album cover
		if len(item.Track.Album.Images) > 0 {
			record.AlbumCoverURL = item.Track.Album.Images[0].URL
		}

		// Get preview URL
		if item.Track.PreviewURL != "" {
			record.PreviewURL = item.Track.PreviewURL
		}

		err := s.historyRepo.InsertHistory(ctx, record)
		if err != nil {
			logger.Error("Failed to insert history record", zap.Error(err))
			continue
		}
		inserted++

		// Track newest cursor
		cursorMs := playedAt.UnixMilli()
		if cursorMs > newestCursorMs {
			newestCursorMs = cursorMs
		}
	}

	// Update cursor
	if newestCursorMs > 0 {
		err = s.historyRepo.UpsertSyncCursor(ctx, user.ID, newestCursorMs, inserted)
		if err != nil {
			logger.Error("Failed to update sync cursor", zap.Error(err))
		}
	}

	logger.Info("Sync completed", zap.Int("inserted", inserted), zap.Int("total_fetched", len(results)))
}
