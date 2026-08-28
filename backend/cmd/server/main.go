package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NapasinNgam/spotify-diary/internal/config"
	"github.com/NapasinNgam/spotify-diary/internal/job"
	"github.com/NapasinNgam/spotify-diary/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	// Load config
	cfg := config.Load()

	// Setup logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Connect to database
	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	// Verify DB connection
	if err := dbPool.Ping(context.Background()); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}
	logger.Info("Connected to database")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Spotify Music Diary",
	})

	// Setup routes
	router.Setup(app, cfg, dbPool)

	// Start scheduler (cron jobs)
	scheduler := job.NewScheduler(cfg, dbPool, logger)
	scheduler.Start()
	defer scheduler.Stop()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("Shutting down server...")
		app.Shutdown()
	}()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.Info("Starting server", zap.String("address", addr))

	if err := app.Listen(addr); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
