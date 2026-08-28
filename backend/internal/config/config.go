package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort  string
	FrontendURL string

	// Spotify
	SpotifyClientID     string
	SpotifyClientSecret string
	SpotifyRedirectURI  string

	// Database
	DatabaseURL string

	// JWT
	JWTSecret      string
	JWTExpiryHours int

	// Cron
	SyncIntervalMinutes int

	// Notification
	SMTPHost          string
	SMTPPort          string
	SMTPUser          string
	SMTPPass          string
	NotificationEmail string
}

func Load() *Config {
	// Load .env file (ignore error if not found — use system env)
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "168"))
	syncInterval, _ := strconv.Atoi(getEnv("SYNC_INTERVAL_MINUTES", "60"))

	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://127.0.0.1:5173"),

		SpotifyClientID:     getEnv("SPOTIFY_CLIENT_ID", ""),
		SpotifyClientSecret: getEnv("SPOTIFY_CLIENT_SECRET", ""),
		SpotifyRedirectURI:  getEnv("SPOTIFY_REDIRECT_URI", "http://127.0.0.1:8080/api/auth/callback"),

		DatabaseURL: getEnv("DATABASE_URL", "postgres://diary:diary_secret@127.0.0.1:5432/music_diary?sslmode=disable"),

		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiryHours: jwtExpiry,

		SyncIntervalMinutes: syncInterval,

		SMTPHost:          getEnv("SMTP_HOST", ""),
		SMTPPort:          getEnv("SMTP_PORT", "587"),
		SMTPUser:          getEnv("SMTP_USER", ""),
		SMTPPass:          getEnv("SMTP_PASS", ""),
		NotificationEmail: getEnv("NOTIFICATION_EMAIL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
