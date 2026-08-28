import { useState, useEffect } from 'react'
import { api } from '../lib/api'

interface DailyStats {
  date: string
  total_tracks: number
  unique_tracks: number
  unique_artists: number
  total_duration_ms: number
}

export default function NewsPage() {
  const [stats, setStats] = useState<DailyStats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchStats()
  }, [])

  async function fetchStats() {
    try {
      const res = await api.get('/api/news/yesterday')
      setStats(res.data)
    } catch {
      setStats(null)
    }
    setLoading(false)
  }

  function formatDuration(ms: number): string {
    const hours = Math.floor(ms / 3600000)
    const minutes = Math.floor((ms % 3600000) / 60000)
    if (hours > 0) return `${hours}h ${minutes}m`
    return `${minutes}m`
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">📰 Daily News</h1>
          <p className="text-sm text-[var(--text-muted)]">Yesterday's listening stats</p>
        </div>
        <button className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-gradient-to-r from-purple-600 via-red-500 to-orange-400 text-white text-xs font-semibold">
          📷 Share to IG
        </button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-3 gap-3">
        <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-lg p-4 text-center">
          <div className="text-2xl font-bold text-[var(--spotify-green)] font-mono">
            {loading ? '...' : (stats?.total_tracks ?? 0)}
          </div>
          <div className="text-xs text-[var(--text-muted)] mt-1">Tracks Played</div>
        </div>
        <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-lg p-4 text-center">
          <div className="text-2xl font-bold text-[var(--spotify-green)] font-mono">
            {loading ? '...' : (stats ? formatDuration(stats.total_duration_ms) : '0m')}
          </div>
          <div className="text-xs text-[var(--text-muted)] mt-1">Listening Time</div>
        </div>
        <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-lg p-4 text-center">
          <div className="text-2xl font-bold text-[var(--spotify-green)] font-mono">
            {loading ? '...' : (stats?.unique_artists ?? 0)}
          </div>
          <div className="text-xs text-[var(--text-muted)] mt-1">Artists</div>
        </div>
      </div>

      {/* Info message */}
      {!loading && !stats && (
        <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-xl p-6 text-center">
          <p className="text-4xl mb-4">🎵</p>
          <p className="text-[var(--text-secondary)] text-sm">
            No listening data yet for yesterday.
          </p>
          <p className="text-[var(--text-muted)] text-xs mt-2">
            Your listening history syncs automatically every hour.<br/>
            Stats will appear here after your first full day of listening.
          </p>
        </div>
      )}

      {/* Suggestion placeholder */}
      <div className="bg-gradient-to-br from-[#1a2a1a] to-[#0a1a0a] border border-[rgba(29,185,84,0.3)] rounded-xl p-5">
        <div className="flex items-center gap-2 mb-3">
          <span className="text-sm">✨</span>
          <span className="text-xs text-[var(--text-muted)] uppercase tracking-wider font-semibold">Today's Suggested Song</span>
        </div>
        <p className="text-sm text-[var(--text-secondary)]">
          Configure your genre playlists in Settings to get daily song suggestions.
        </p>
      </div>
    </div>
  )
}
