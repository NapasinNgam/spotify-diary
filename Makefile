.PHONY: dev dev-backend dev-frontend db-up db-down migrate-up migrate-down build

# ═══ Development ═══
dev-backend:
	cd backend && air

dev-frontend:
	cd frontend && npm run dev

# ═══ Database ═══
db-up:
	docker compose up -d

db-down:
	docker compose down

# ═══ Migrations ═══
migrate-up:
	cd backend && migrate -path migrations -database "postgres://diary:diary_secret@127.0.0.1:5432/music_diary?sslmode=disable" up

migrate-down:
	cd backend && migrate -path migrations -database "postgres://diary:diary_secret@127.0.0.1:5432/music_diary?sslmode=disable" down 1

# ═══ Build ═══
build-backend:
	cd backend && go build -o bin/server ./cmd/server

build-frontend:
	cd frontend && npm run build

# ═══ Install ═══
install:
	cd backend && go mod download
	cd frontend && npm install
