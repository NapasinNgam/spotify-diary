package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "[REDACTED_CONN_STRING]://127.0.0.1:5432/music_diary?sslmode=disable"
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close(ctx)

	fmt.Println("Connected to database!")

	// Find all .up.sql files
	files, err := filepath.Glob("migrations/*.up.sql")
	if err != nil {
		log.Fatalf("Failed to find migration files: %v", err)
	}
	sort.Strings(files)

	for _, file := range files {
		fmt.Printf("Running: %s ... ", filepath.Base(file))

		sql, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read %s: %v", file, err)
		}

		_, err = conn.Exec(ctx, string(sql))
		if err != nil {
			fmt.Printf("SKIPPED (already exists or error: %v)\n", err)
			continue
		}
		fmt.Println("OK")
	}

	fmt.Println("\nAll migrations completed!")
}
