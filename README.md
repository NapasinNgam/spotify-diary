# 🎵 Spotify Music Diary

Personal music diary that logs your daily songs, shows monthly listening stats, and creates half-year recaps — powered by Spotify API.

## Features

- 📰 **Daily News** — Yesterday's listening stats + top tracks per genre
- 📅 **Diary Calendar** — Log 1 song per day, displayed on a calendar
- 📊 **Monthly Records** — Top 10 most-played tracks from Spotify history
- 🏆 **Half-Year Recap** — Pick your top 5 songs + honorable mentions
- 📚 **History Bookshelf** — Archive of past records

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.22+ (Fiber) |
| Frontend | React 18 + TypeScript + Vite |
| UI | Shadcn/ui + Tailwind CSS |
| Database | PostgreSQL 16 |
| Auth | Spotify OAuth2 → JWT |
| API | Spotify Web API (zmb3/spotify) |

## Quick Start

```bash
# 1. Copy environment file
cp .env.example .env
# Edit .env with your Spotify credentials

# 2. Start database
docker compose up -d

# 3. Run migrations
cd backend
make migrate-up

# 4. Start backend (with hot-reload)
air

# 5. Start frontend (new terminal)
cd frontend
npm install
npm run dev

# 6. Open browser → http://127.0.0.1:5173
```

## Prerequisites

- Go 1.22+
- Node.js 20+
- PostgreSQL 16+ (or Docker)
- Spotify Developer Account

## Environment Variables

See `.env.example` for all required variables.

## License

MIT
