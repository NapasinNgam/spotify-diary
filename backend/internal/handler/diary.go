package handler

import (
	"context"

	"github.com/NapasinNgam/spotify-diary/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type DiaryHandler struct {
	diaryRepo *repository.DiaryRepository
}

func NewDiaryHandler(diaryRepo *repository.DiaryRepository) *DiaryHandler {
	return &DiaryHandler{diaryRepo: diaryRepo}
}

// GetCalendar returns diary entries for a given month
func (h *DiaryHandler) GetCalendar(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	month := c.Query("month") // format: YYYY-MM

	if month == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "month parameter is required (format: YYYY-MM)",
		})
	}

	entries, err := h.diaryRepo.GetByMonth(context.Background(), userID, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch diary entries",
		})
	}

	return c.JSON(fiber.Map{
		"entries": entries,
		"month":   month,
	})
}

// LogSong saves today's song
func (h *DiaryHandler) LogSong(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var body struct {
		Date          string `json:"date"` // YYYY-MM-DD
		TrackID       string `json:"track_id"`
		TrackName     string `json:"track_name"`
		ArtistName    string `json:"artist_name"`
		AlbumName     string `json:"album_name"`
		AlbumCoverURL string `json:"album_cover_url"`
		PreviewURL    string `json:"preview_url"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if body.TrackID == "" || body.Date == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "track_id and date are required",
		})
	}

	entry, err := h.diaryRepo.Upsert(context.Background(), repository.UpsertDiaryParams{
		UserID:        userID,
		LogDate:       body.Date,
		TrackID:       body.TrackID,
		TrackName:     body.TrackName,
		ArtistName:    body.ArtistName,
		AlbumName:     body.AlbumName,
		AlbumCoverURL: body.AlbumCoverURL,
		PreviewURL:    body.PreviewURL,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save diary entry",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"entry": entry,
	})
}

// DeleteLog removes a diary entry
func (h *DiaryHandler) DeleteLog(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	date := c.Params("date") // YYYY-MM-DD

	err := h.diaryRepo.Delete(context.Background(), userID, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete diary entry",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Deleted successfully",
	})
}
