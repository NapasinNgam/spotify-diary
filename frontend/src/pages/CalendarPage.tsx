import { useState, useEffect } from 'react'
import { api } from '../lib/api'
import { format, startOfMonth, endOfMonth, eachDayOfInterval, getDay, isSameDay, isToday } from 'date-fns'

interface DiaryEntry {
  id: number
  log_date: string
  track_id: string
  track_name: string
  artist_name: string
  album_name: string
  album_cover_url: string
  preview_url: string
}

interface SearchResult {
  id: string
  name: string
  artist: string
  album: string
  album_cover: string
  preview_url: string
}

export default function CalendarPage() {
  const [currentDate, setCurrentDate] = useState(new Date())
  const [entries, setEntries] = useState<DiaryEntry[]>([])
  const [searchQuery, setSearchQuery] = useState('')
  const [searchResults, setSearchResults] = useState<SearchResult[]>([])
  const [searching, setSearching] = useState(false)
  const [selectedDate, setSelectedDate] = useState<string>(format(new Date(), 'yyyy-MM-dd'))

  const currentMonth = format(currentDate, 'yyyy-MM')

  useEffect(() => {
    fetchEntries()
  }, [currentMonth])

  async function fetchEntries() {
    try {
      const res = await api.get(`/api/diary/calendar?month=${currentMonth}`)
      setEntries(res.data.entries || [])
    } catch (err) {
      // No entries yet — that's OK
      setEntries([])
    }
  }

  async function searchSpotify() {
    if (!searchQuery.trim()) return
    setSearching(true)
    try {
      const res = await api.get(`/api/spotify/search?q=${encodeURIComponent(searchQuery)}`)
      setSearchResults(res.data.tracks || [])
    } catch (err) {
      // Spotify search endpoint not yet implemented — show mock
      setSearchResults([])
    }
    setSearching(false)
  }

  async function saveSong(track: SearchResult) {
    try {
      await api.post('/api/diary/log', {
        date: selectedDate,
        track_id: track.id,
        track_name: track.name,
        artist_name: track.artist,
        album_name: track.album,
        album_cover_url: track.album_cover,
        preview_url: track.preview_url,
      })
      setSearchResults([])
      setSearchQuery('')
      fetchEntries()
    } catch (err) {
      console.error('Failed to save:', err)
    }
  }

  function prevMonth() {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() - 1, 1))
  }
  function nextMonth() {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 1))
  }

  // Calendar grid calculations
  const monthStart = startOfMonth(currentDate)
  const monthEnd = endOfMonth(currentDate)
  const days = eachDayOfInterval({ start: monthStart, end: monthEnd })
  const startDayOfWeek = getDay(monthStart) // 0=Sun, adjust for Mon start
  const emptySlots = startDayOfWeek === 0 ? 6 : startDayOfWeek - 1

  function getEntryForDate(date: Date): DiaryEntry | undefined {
    const dateStr = format(date, 'yyyy-MM-dd')
    // log_date from backend might be ISO timestamp or date string
    return entries.find(e => e.log_date.substring(0, 10) === dateStr)
  }

  return (
    <div className="space-y-5">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">📅 Diary Calendar</h1>
          <p className="text-sm text-[var(--text-muted)]">
            {format(currentDate, 'MMMM yyyy')} • Log 1 song per day
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
        <span className="font-semibold text-lg">{format(currentDate, 'MMMM yyyy')}</span>
        <button onClick={nextMonth} className="px-3 py-1 bg-[var(--bg-card)] border border-[var(--border)] rounded-md text-sm hover:border-[var(--spotify-green)] transition-colors">
          Next →
        </button>
      </div>

      {/* Search + Log Song */}
      <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-xl p-4">
        <div className="flex items-center gap-2 mb-3">
          <span className="text-[var(--spotify-green)] text-sm">＋</span>
          <span className="text-sm text-[var(--text-secondary)]">
            Log song for <strong className="text-[var(--text-primary)]">{selectedDate}</strong>
          </span>
        </div>
        <div className="flex gap-2">
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && searchSpotify()}
            placeholder="Search Spotify for a song..."
            className="flex-1 px-4 py-2 bg-[var(--bg-primary)] border border-[var(--border)] rounded-lg text-sm text-[var(--text-primary)] placeholder:text-[var(--text-muted)] outline-none focus:border-[var(--spotify-green)] transition-colors"
          />
          <button
            onClick={searchSpotify}
            disabled={searching}
            className="px-4 py-2 bg-[var(--spotify-green)] text-black font-semibold text-sm rounded-lg hover:bg-[var(--spotify-green-hover)] transition-colors disabled:opacity-50"
          >
            {searching ? '...' : 'Search'}
          </button>
        </div>

        {/* Search Results */}
        {searchResults.length > 0 && (
          <div className="mt-3 space-y-1 max-h-48 overflow-y-auto">
            {searchResults.map((track) => (
              <div key={track.id} className="flex items-center gap-3 p-2 rounded-lg hover:bg-[var(--bg-hover)] transition-colors">
                <img src={track.album_cover} alt="" className="w-10 h-10 rounded" />
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium truncate">{track.name}</p>
                  <p className="text-xs text-[var(--text-muted)] truncate">{track.artist} • {track.album}</p>
                </div>
                <button
                  onClick={() => saveSong(track)}
                  className="px-3 py-1 bg-[var(--spotify-green)] text-black text-xs font-semibold rounded-full hover:bg-[var(--spotify-green-hover)] transition-colors whitespace-nowrap"
                >
                  Save
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Calendar Grid */}
      <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-xl p-4">
        {/* Day headers */}
        <div className="grid grid-cols-7 gap-1 mb-2">
          {['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'].map(day => (
            <div key={day} className="text-center text-xs font-semibold text-[var(--text-muted)] py-1">
              {day}
            </div>
          ))}
        </div>

        {/* Calendar cells */}
        <div className="grid grid-cols-7 gap-1">
          {/* Empty slots */}
          {Array.from({ length: emptySlots }).map((_, i) => (
            <div key={`empty-${i}`} className="aspect-square" />
          ))}

          {/* Day cells */}
          {days.map((day) => {
            const entry = getEntryForDate(day)
            const dateStr = format(day, 'yyyy-MM-dd')
            const isSelected = dateStr === selectedDate
            const today = isToday(day)

            return (
              <div
                key={dateStr}
                onClick={() => setSelectedDate(dateStr)}
                className={`aspect-square rounded-lg border flex flex-col items-center justify-center p-1 cursor-pointer transition-all hover:border-[var(--spotify-green)] ${
                  today ? 'border-[var(--spotify-green)] shadow-[0_0_0_1px_var(--spotify-green)]' :
                  isSelected ? 'border-[var(--spotify-green)] bg-[var(--bg-hover)]' :
                  'border-[var(--border)] bg-[var(--bg-secondary)]'
                }`}
              >
                <span className={`text-[10px] mb-0.5 ${today ? 'text-[var(--spotify-green)] font-bold' : 'text-[var(--text-muted)]'}`}>
                  {format(day, 'd')}
                </span>

                {entry ? (
                  <>
                    <img
                      src={entry.album_cover_url}
                      alt={entry.track_name}
                      className="w-8 h-8 rounded object-cover"
                    />
                    <span className="text-[8px] text-[var(--text-secondary)] text-center truncate w-full mt-0.5">
                      {entry.track_name}
                    </span>
                  </>
                ) : (
                  <span className="text-lg text-[var(--text-muted)] opacity-30">+</span>
                )}
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}
