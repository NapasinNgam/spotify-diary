package handler

import (
	"context"

	"github.com/NapasinNgam/spotify-diary/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type RecapHandler struct {
	recapRepo *repository.RecapRepository
}

func NewRecapHandler(recapRepo *repository.RecapRepository) *RecapHandler {
	return &RecapHandler{recapRepo: recapRepo}
}

// GetRecap returns the half-year recap (top 5 + honorable mentions)
func (h *RecapHandler) GetRecap(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	period := c.Query("period") // e.g. "2026-H1"

	if period == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "period parameter is required (e.g. 2026-H1)",
		})
	}

	recap, err := h.recapRepo.GetByPeriod(context.Background(), userID, period)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch recap",
		})
	}

	return c.JSON(fiber.Map{
		"period": period,
		"tracks": recap,
	})
}

// SaveRecap saves/updates a half-year recap entry
func (h *RecapHandler) SaveRecap(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var body struct {
		Period        string `json:"period"`
		Rank          int    `json:"rank"`
		TrackID       string `json:"track_id"`
		TrackName     string `json:"track_name"`
		ArtistName    string `json:"artist_name"`
		AlbumCoverURL string `json:"album_cover_url"`
		PreviewURL    string `json:"preview_url"`
		Description   string `json:"description"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.recapRepo.Upsert(context.Background(), repository.UpsertRecapParams{
		UserID:        userID,
		Period:        body.Period,
		Rank:          body.Rank,
		TrackID:       body.TrackID,
		TrackName:     body.TrackName,
		ArtistName:    body.ArtistName,
		AlbumCoverURL: body.AlbumCoverURL,
		PreviewURL:    body.PreviewURL,
		Description:   body.Description,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save recap",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Recap saved successfully",
	})
}

// ListPeriods returns all available recap periods
func (h *RecapHandler) ListPeriods(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	periods, err := h.recapRepo.ListPeriods(context.Background(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list periods",
		})
	}

	return c.JSON(fiber.Map{
		"periods": periods,
	})
}
