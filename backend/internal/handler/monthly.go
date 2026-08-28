package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NapasinNgam/spotify-diary/internal/config"
	"github.com/NapasinNgam/spotify-diary/internal/model"
	"github.com/NapasinNgam/spotify-diary/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type MonthlyHandler struct {
	cfg         *config.Config
	summaryRepo *repository.SummaryRepository
	userRepo    *repository.UserRepository
}

func NewMonthlyHandler(cfg *config.Config, summaryRepo *repository.SummaryRepository, userRepo *repository.UserRepository) *MonthlyHandler {
	return &MonthlyHandler{
		cfg:         cfg,
		summaryRepo: summaryRepo,
		userRepo:    userRepo,
	}
}

// GetMonthlyRecords returns the monthly summary (top 10 tracks)
func (h *MonthlyHandler) GetMonthlyRecords(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	month := c.Query("month") // format: YYYY-MM

	if month == "" {
		// Default to current month
		month = time.Now().Format("2006-01")
	}

	summary, err := h.summaryRepo.GetMonthlySummary(context.Background(), userID, month)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(fiber.Map{
				"exists": false,
				"month":  month,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch monthly records",
		})
	}

	// Parse top_tracks JSON for frontend
	var topTracks []model.TopTrackItem
	if summary.TopTracks != "" && summary.TopTracks != "[]" {
		_ = json.Unmarshal([]byte(summary.TopTracks), &topTracks)
	}

	// Format period dates as strings for frontend
	var periodStartStr, periodEndStr string
	if summary.PeriodStart != nil {
		periodStartStr = summary.PeriodStart.Format("2006-01-02")
	}
	if summary.PeriodEnd != nil {
		periodEndStr = summary.PeriodEnd.Format("2006-01-02")
	}

	return c.JSON(fiber.Map{
		"exists":            true,
		"month":             summary.SummaryMonth,
		"total_plays":       summary.TotalPlays,
		"unique_tracks":     summary.UniqueTracks,
		"unique_artists":    summary.UniqueArtists,
		"total_duration_ms": summary.TotalDurationMs,
		"top_tracks":        topTracks,
		"period_start":      periodStartStr,
		"period_end":        periodEndStr,
		"source":            summary.Source,
		"generated_at":      summary.GeneratedAt,
	})
}

// GenerateMonthly manually triggers monthly summary generation (first-time bootstrap)
func (h *MonthlyHandler) GenerateMonthly(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	spotifyID := c.Locals("spotify_id").(string)

	// Get user tokens
	user, err := h.userRepo.GetBySpotifyID(context.Background(), spotifyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	// Create Spotify client
	auth := spotifyauth.New(
		spotifyauth.WithClientID(h.cfg.SpotifyClientID),
		spotifyauth.WithClientSecret(h.cfg.SpotifyClientSecret),
	)

	token := &oauth2.Token{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		Expiry:       user.TokenExpiresAt,
	}

	httpClient := auth.Client(context.Background(), token)
	client := spotify.New(httpClient)

	// Fetch top tracks (short_term = ~4 weeks)
	topTracks, err := client.CurrentUsersTopTracks(
		context.Background(),
		spotify.Limit(10),
		spotify.Timerange("short_term"),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch top tracks from Spotify: " + err.Error(),
		})
	}

	// Build top tracks list
	var tracks []model.TopTrackItem
	uniqueArtists := make(map[string]bool)

	for i, t := range topTracks.Tracks {
		albumCover := ""
		if len(t.Album.Images) > 0 {
			albumCover = t.Album.Images[0].URL
		}
		artistName := ""
		if len(t.Artists) > 0 {
			artistName = t.Artists[0].Name
			uniqueArtists[string(t.Artists[0].ID)] = true
		}

		tracks = append(tracks, model.TopTrackItem{
			Rank:          i + 1,
			TrackID:       string(t.ID),
			TrackName:     t.Name,
			ArtistName:    artistName,
			AlbumCoverURL: albumCover,
			PreviewURL:    t.PreviewURL,
			PlayCount:     0, // Spotify doesn't expose play count from top tracks
		})
	}

	// Calculate period (today - 28 days to today)
	now := time.Now()
	periodEnd := now
	periodStart := now.AddDate(0, 0, -28)
	summaryMonth := now.Format("2006-01")

	// JSON encode tracks
	tracksJSON, _ := json.Marshal(tracks)

	// Save to DB
	err = h.summaryRepo.UpsertMonthlySummary(context.Background(), repository.UpsertMonthlySummaryParams{
		UserID:          userID,
		SummaryMonth:    summaryMonth,
		TotalPlays:      0, // Not available from top tracks endpoint
		UniqueTracks:    len(topTracks.Tracks),
		UniqueArtists:   len(uniqueArtists),
		TotalDurationMs: 0, // Not available from top tracks endpoint
		TopTracks:       string(tracksJSON),
		TopArtists:      "[]",
		GenreBreakdown:  "{}",
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
		Source:          "spotify_top_tracks",
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save monthly summary: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":      "Monthly records generated successfully",
		"month":        summaryMonth,
		"period_start": periodStart.Format("2006-01-02"),
		"period_end":   periodEnd.Format("2006-01-02"),
		"tracks_count": len(tracks),
		"source":       "spotify_top_tracks",
		"note":         fmt.Sprintf("Data from Spotify Top Tracks (short_term ~4 weeks: %s to %s)", periodStart.Format("2006-01-02"), periodEnd.Format("2006-01-02")),
	})
}

// ListMonths returns available monthly summaries
func (h *MonthlyHandler) ListMonths(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	months, err := h.summaryRepo.ListAvailableMonths(context.Background(), userID)
	if err != nil {
		return c.JSON(fiber.Map{
			"months": []string{},
		})
	}

	return c.JSON(fiber.Map{
		"months": months,
	})
}
