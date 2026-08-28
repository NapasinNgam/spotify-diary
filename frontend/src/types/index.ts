export interface User {
  id: number
  spotify_id: string
  display_name: string
  email: string
  avatar_url: string
  created_at: string
  updated_at: string
}

export interface DailyLog {
  id: number
  user_id: number
  log_date: string
  track_id: string
  track_name: string
  artist_name: string
  album_name: string
  album_cover_url: string
  preview_url: string
  created_at: string
}

export interface ListeningHistory {
  id: number
  user_id: number
  track_id: string
  track_name: string
  artist_id: string
  artist_name: string
  album_name: string
  album_cover_url: string
  preview_url: string
  duration_ms: number
  played_at: string
  played_date: string
  played_month: string
  genre_category: string
}

export interface MonthlySummary {
  id: number
  user_id: number
  summary_month: string
  total_plays: number
  unique_tracks: number
  unique_artists: number
  total_duration_ms: number
  top_tracks: TopTrack[]
  generated_at: string
}

export interface TopTrack {
  rank: number
  track_id: string
  track_name: string
  artist_name: string
  album_cover_url: string
  play_count: number
}

export interface HalfYearRecap {
  id: number
  user_id: number
  period: string
  rank: number
  track_id: string
  track_name: string
  artist_name: string
  album_cover_url: string
  preview_url: string
  description: string
  created_at: string
  updated_at: string
}

export interface DailyStats {
  date: string
  total_tracks: number
  unique_tracks: number
  unique_artists: number
  total_duration_ms: number
}
