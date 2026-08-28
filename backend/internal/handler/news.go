package handler

import (
	"context"

	"github.com/NapasinNgam/spotify-diary/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type NewsHandler struct {
	historyRepo *repository.HistoryRepository
	diaryRepo   *repository.DiaryRepository
}

func NewNewsHandler(historyRepo *repository.HistoryRepository, diaryRepo *repository.DiaryRepository) *NewsHandler {
	return &NewsHandler{
		historyRepo: historyRepo,
		diaryRepo:   diaryRepo,
	}
}

// GetYesterdayStats returns listening stats for yesterday
func (h *NewsHandler) GetYesterdayStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	stats, err := h.historyRepo.GetDailyStats(context.Background(), userID, "yesterday")
	if err != nil {
		// Return empty stats if no data
		return c.JSON(fiber.Map{
			"date":              "",
			"total_tracks":      0,
			"unique_tracks":     0,
			"unique_artists":    0,
			"total_duration_ms": 0,
		})
	}

	return c.JSON(stats)
}

// GetTopByGenre returns top 3 tracks per genre for yesterday
func (h *NewsHandler) GetTopByGenre(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	topTracks, err := h.historyRepo.GetTopByGenre(context.Background(), userID, "yesterday", 3)
	if err != nil {
		return c.JSON(fiber.Map{
			"genres": map[string]interface{}{},
		})
	}

	return c.JSON(fiber.Map{
		"genres": topTracks,
	})
}

// GetSuggestion returns a random song from configured genre playlists
func (h *NewsHandler) GetSuggestion(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"suggestion": nil,
		"message":    "Configure your genre playlists first",
	})
}
