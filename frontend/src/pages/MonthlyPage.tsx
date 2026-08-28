import { useState, useEffect } from 'react'
import { api } from '../lib/api'
import { format, subMonths } from 'date-fns'

interface TopTrack {
  rank: number
  track_id: string
  track_name: string
  artist_name: string
  album_cover_url: string
  preview_url: string
  play_count: number
}

interface MonthlyData {
  exists: boolean
  month: string
  total_plays: number
  unique_tracks: number
  unique_artists: number
  total_duration_ms: number
  top_tracks: TopTrack[]
  period_start: string | null
  period_end: string | null
  source: string
  generated_at: string
}

export default function MonthlyPage() {
  const [data, setData] = useState<MonthlyData | null>(null)
  const [loading, setLoading] = useState(true)
  const [generating, setGenerating] = useState(false)
  const [currentMonth, setCurrentMonth] = useState(format(new Date(), 'yyyy-MM'))

  useEffect(() => {
    fetchRecords()
  }, [currentMonth])

  async function fetchRecords() {
    setLoading(true)
    try {
      const res = await api.get(`/api/monthly/records?month=${currentMonth}`)
      setData(res.data)
    } catch {
      setData(null)
    }
    setLoading(false)
  }

  async function generateNow() {
    setGenerating(true)
    try {
      await api.post('/api/monthly/generate')
      await fetchRecords()
    } catch (err) {
      console.error('Failed to generate:', err)
    }
    setGenerating(false)
  }

  function prevMonth() {
    const d = new Date(currentMonth + '-01')
    setCurrentMonth(format(subMonths(d, 1), 'yyyy-MM'))
  }
  function nextMonth() {
    const d = new Date(currentMonth + '-01')
    const next = new Date(d.getFullYear(), d.getMonth() + 1, 1)
    setCurrentMonth(format(next, 'yyyy-MM'))
  }

  function formatDuration(ms: number): string {
    const hours = Math.floor(ms / 3600000)
    const minutes = Math.floor((ms % 3600000) / 60000)
    if (hours > 0) return `${hours}h ${minutes}m`
    return `${minutes}m`
  }

  function formatPeriod(): string {
    if (!data?.period_start || !data?.period_end) return ''
    const start = new Date(data.period_start)
    const end = new Date(data.period_end)
    return `${format(start, 'd MMM')} – ${format(end, 'd MMM yyyy')}`
  }

  const maxPlayCount = data?.top_tracks?.[0]?.play_count || 1

  return (
    <div className="space-y-5">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">📊 Monthly Records</h1>
          <p className="text-sm text-[var(--text-muted)]">
            Top 10 most played tracks
          </p>
        </div>
        <button className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-gradient-to-r from-purple-600 via-red-500 to-orange-400 text-white text-xs font-semibold">
          📷 Share to IG
        </button>
      </div>

      {/* Month Navigation */}
      <div className="flex items-center gap-4">
        <button onClick={prevMonth} className="px-3 py-1 bg-[var(--bg-card)] border border-[var(--border)] rounded-md text-sm hover:border-[var(--spotify-green)] transition-colors">
          ← Prev
        </button>
        <span className="font-semibold text-lg">
          {format(new Date(currentMonth + '-01'), 'MMMM yyyy')}
        </span>
        <button onClick={nextMonth} className="px-3 py-1 bg-[var(--bg-card)] border border-[var(--border)] rounded-md text-sm hover:border-[var(--spotify-green)] transition-colors">
          Next →
        </button>
      </div>

      {loading ? (
        <div className="text-center py-12 text-[var(--text-muted)]">Loading...</div>
      ) : data?.exists ? (
        <>
          {/* Period info */}
          <div className="flex items-center gap-2 text-xs text-[var(--text-muted)]">
            <span className="px-2 py-1 bg-[rgba(29,185,84,0.1)] text-[var(--spotify-green)] rounded font-medium">
              {data.source === 'spotify_top_tracks' ? '📡 From Spotify Top Tracks' : '📊 From Listening History'}
            </span>
            {data.period_start && (
              <span>• Data period: {formatPeriod()}</span>
            )}
          </div>

          {/* Stats */}
          {data.source !== 'spotify_top_tracks' && (
            <div className="grid grid-cols-3 gap-3">
              <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-lg p-4 text-center">
                <div className="text-2xl font-bold text-[var(--spotify-green)] font-mono">{data.total_plays}</div>
                <div className="text-xs text-[var(--text-muted)] mt-1">Total Plays</div>
              </div>
              <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-lg p-4 text-center">
                <div className="text-2xl font-bold text-[var(--spotify-green)] font-mono">
                  {data.total_duration_ms > 0 ? formatDuration(data.total_duration_ms) : `${data.unique_tracks} tracks`}
                </div>
                <div className="text-xs text-[var(--text-muted)] mt-1">
                  {data.total_duration_ms > 0 ? 'Listening Time' : 'Unique Tracks'}
                </div>
              </div>
              <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-lg p-4 text-center">
                <div className="text-2xl font-bold text-[var(--spotify-green)] font-mono">{data.unique_artists}</div>
                <div className="text-xs text-[var(--text-muted)] mt-1">Artists</div>
              </div>
            </div>
          )}

          {/* Top 10 Tracks */}
          <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-xl overflow-hidden">
            <div className="p-4 border-b border-[var(--border)]">
              <h3 className="text-sm font-semibold flex items-center gap-2">
                <span className="text-orange-400">🔥</span> Top Tracks
              </h3>
            </div>
            <div className="divide-y divide-[var(--border)]">
              {data.top_tracks?.map((track, i) => (
                <div key={track.track_id + i} className="flex items-center gap-3 px-4 py-3 hover:bg-[var(--bg-hover)] transition-colors">
                  {/* Rank */}
                  <span className={`text-lg font-bold w-7 text-center ${
                    i === 0 ? 'text-yellow-400' :
                    i === 1 ? 'text-gray-400' :
                    i === 2 ? 'text-amber-600' :
                    'text-[var(--text-muted)]'
                  }`}>
                    {track.rank}
                  </span>

                  {/* Album cover */}
                  <img
                    src={track.album_cover_url}
                    alt={track.track_name}
                    className="w-11 h-11 rounded-md object-cover"
                  />

                  {/* Track info */}
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">{track.track_name}</p>
                    <p className="text-xs text-[var(--text-muted)] truncate">{track.artist_name}</p>
                  </div>

                  {/* Play count */}
                  {track.play_count > 0 && (
                    <span className="text-xs text-[var(--text-secondary)] font-mono whitespace-nowrap">
                      {track.play_count}x
                    </span>
                  )}

                  {/* Progress bar */}
                  {track.play_count > 0 && (
                    <div className="w-20 h-1 bg-[var(--border)] rounded-full overflow-hidden">
                      <div
                        className="h-full bg-[var(--spotify-green)] rounded-full"
                        style={{ width: `${(track.play_count / maxPlayCount) * 100}%` }}
                      />
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        </>
      ) : (
        /* No data — show generate button */
        <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-xl p-8 text-center">
          <p className="text-4xl mb-4">📊</p>
          <h3 className="text-lg font-semibold mb-2">No monthly records yet</h3>
          <p className="text-sm text-[var(--text-muted)] mb-6">
            Generate your first monthly recap by pulling your top tracks from Spotify (last 4 weeks).
          </p>
          <button
            onClick={generateNow}
            disabled={generating}
            className="px-6 py-3 bg-[var(--spotify-green)] text-black font-semibold rounded-full hover:bg-[var(--spotify-green-hover)] transition-colors disabled:opacity-50"
          >
            {generating ? 'Generating...' : '✨ Generate Now'}
          </button>
          <p className="text-xs text-[var(--text-muted)] mt-4">
            After this, records will auto-generate on the 1st of each month.
          </p>
        </div>
      )}
    </div>
  )
}
