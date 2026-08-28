package handler

import (
	"context"

	"github.com/NapasinNgam/spotify-diary/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type MonthlyHandler struct {
	summaryRepo *repository.SummaryRepository
}

func NewMonthlyHandler(summaryRepo *repository.SummaryRepository) *MonthlyHandler {
	return &MonthlyHandler{summaryRepo: summaryRepo}
}

// GetMonthlyRecords returns the monthly summary (top 10 tracks)
func (h *MonthlyHandler) GetMonthlyRecords(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	month := c.Query("month") // format: YYYY-MM

	if month == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "month parameter is required (format: YYYY-MM)",
		})
	}

	summary, err := h.summaryRepo.GetMonthlySummary(context.Background(), userID, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch monthly records",
		})
	}

	return c.JSON(summary)
}

// ListMonths returns available monthly summaries
func (h *MonthlyHandler) ListMonths(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	months, err := h.summaryRepo.ListAvailableMonths(context.Background(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list months",
		})
	}

	return c.JSON(fiber.Map{
		"months": months,
	})
}
