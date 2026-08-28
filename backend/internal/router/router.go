package router

import (
	"github.com/NapasinNgam/spotify-diary/internal/config"
	"github.com/NapasinNgam/spotify-diary/internal/handler"
	"github.com/NapasinNgam/spotify-diary/internal/middleware"
	"github.com/NapasinNgam/spotify-diary/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(app *fiber.App, cfg *config.Config, db *pgxpool.Pool) {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	historyRepo := repository.NewHistoryRepository(db)
	diaryRepo := repository.NewDiaryRepository(db)
	summaryRepo := repository.NewSummaryRepository(db)
	recapRepo := repository.NewRecapRepository(db)

	// Handlers
	authHandler := handler.NewAuthHandler(cfg, userRepo)
	spotifyHandler := handler.NewSpotifyHandler(cfg, userRepo)
	diaryHandler := handler.NewDiaryHandler(diaryRepo)
	newsHandler := handler.NewNewsHandler(historyRepo, diaryRepo)
	monthlyHandler := handler.NewMonthlyHandler(cfg, summaryRepo, userRepo)
	recapHandler := handler.NewRecapHandler(recapRepo)

	// CORS
	app.Use(middleware.SetupCORS(cfg.FrontendURL))

	// API routes
	api := app.Group("/api")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "spotify-music-diary"})
	})

	// Auth (public)
	auth := api.Group("/auth")
	auth.Get("/login", authHandler.Login)
	auth.Get("/callback", authHandler.Callback)

	// Protected routes
	protected := api.Group("", middleware.AuthRequired(cfg.JWTSecret))

	// User
	protected.Get("/me", authHandler.Me)

	// Spotify
	spotifyGroup := protected.Group("/spotify")
	spotifyGroup.Get("/search", spotifyHandler.Search)

	// Diary
	diary := protected.Group("/diary")
	diary.Get("/calendar", diaryHandler.GetCalendar)
	diary.Post("/log", diaryHandler.LogSong)
	diary.Delete("/log/:date", diaryHandler.DeleteLog)

	// News
	news := protected.Group("/news")
	news.Get("/yesterday", newsHandler.GetYesterdayStats)
	news.Get("/top-genres", newsHandler.GetTopByGenre)
	news.Get("/suggestion", newsHandler.GetSuggestion)

	// Monthly
	monthly := protected.Group("/monthly")
	monthly.Get("/records", monthlyHandler.GetMonthlyRecords)
	monthly.Post("/generate", monthlyHandler.GenerateMonthly)
	monthly.Get("/list", monthlyHandler.ListMonths)

	// Recap
	recap := protected.Group("/recap")
	recap.Get("/", recapHandler.GetRecap)
	recap.Post("/", recapHandler.SaveRecap)
	recap.Get("/periods", recapHandler.ListPeriods)
}
