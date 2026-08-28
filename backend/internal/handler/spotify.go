package handler

import (
	"context"

	"github.com/NapasinNgam/spotify-diary/internal/config"
	"github.com/NapasinNgam/spotify-diary/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type SpotifyHandler struct {
	cfg      *config.Config
	userRepo *repository.UserRepository
}

func NewSpotifyHandler(cfg *config.Config, userRepo *repository.UserRepository) *SpotifyHandler {
	return &SpotifyHandler{cfg: cfg, userRepo: userRepo}
}

// Search searches Spotify for tracks
func (h *SpotifyHandler) Search(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	query := c.Query("q")

	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query parameter 'q' is required",
		})
	}

	// Get user's tokens from DB
	user, err := h.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User not found"})
	}

	// Get full user with tokens
	fullUser, err := h.userRepo.GetBySpotifyID(context.Background(), user.SpotifyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user tokens"})
	}

	// Create Spotify client
	auth := spotifyauth.New(
		spotifyauth.WithClientID(h.cfg.SpotifyClientID),
		spotifyauth.WithClientSecret(h.cfg.SpotifyClientSecret),
	)

	token := &oauth2.Token{
		AccessToken:  fullUser.AccessToken,
		RefreshToken: fullUser.RefreshToken,
		Expiry:       fullUser.TokenExpiresAt,
	}

	httpClient := auth.Client(context.Background(), token)
	client := spotify.New(httpClient)

	// Search
	results, err := client.Search(context.Background(), query, spotify.SearchTypeTrack)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Spotify search failed: " + err.Error(),
		})
	}

	// Format response
	var tracks []fiber.Map
	if results.Tracks != nil {
		for _, t := range results.Tracks.Tracks {
			albumCover := ""
			if len(t.Album.Images) > 0 {
				albumCover = t.Album.Images[0].URL
			}
			artistName := ""
			if len(t.Artists) > 0 {
				artistName = t.Artists[0].Name
			}

			tracks = append(tracks, fiber.Map{
				"id":          string(t.ID),
				"name":        t.Name,
				"artist":      artistName,
				"album":       t.Album.Name,
				"album_cover": albumCover,
				"preview_url": t.PreviewURL,
			})
		}
	}

	return c.JSON(fiber.Map{
		"tracks": tracks,
	})
}
