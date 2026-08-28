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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch stats",
		})
	}

	return c.JSON(stats)
}

// GetTopByGenre returns top 3 tracks per genre for yesterday
func (h *NewsHandler) GetTopByGenre(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	topTracks, err := h.historyRepo.GetTopByGenre(context.Background(), userID, "yesterday", 3)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch top tracks by genre",
		})
	}

	return c.JSON(fiber.Map{
		"genres": topTracks,
	})
}

// GetSuggestion returns a random song from configured genre playlists
func (h *NewsHandler) GetSuggestion(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	// TODO: implement playlist-based suggestion
	_ = userID

	return c.JSON(fiber.Map{
		"suggestion": nil,
		"message":    "Configure your genre playlists first",
	})
}
