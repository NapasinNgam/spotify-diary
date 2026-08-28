import { useState, useEffect, useRef } from 'react'
import { api } from '../lib/api'

const HALF_YEAR_SLOTS = [
  { key: 'overall', label: 'Overall', display: '🎵', type: 'text' },
  { key: 'hm_1', label: '1st Honorable Mention', display: 'HM', type: 'track' },
  { key: 'place_5', label: '5th Place', display: '#5', type: 'track' },
  { key: 'place_4', label: '4th Place', display: '#4', type: 'track' },
  { key: 'hm_2', label: '2nd Honorable Mention', display: 'HM', type: 'track' },
  { key: 'place_3', label: '3rd Place', display: '#3', type: 'track' },
  { key: 'place_2', label: '2nd Place', display: '#2', type: 'track' },
  { key: 'hm_3', label: '3rd Honorable Mention', display: 'HM', type: 'track' },
  { key: 'place_1', label: '1st Place', display: '#1', type: 'track' },
  { key: 'summary', label: 'Summary', display: '📝', type: 'text' },
]

const FULL_YEAR_SLOTS = [
  { key: 'overall', label: 'Overall', display: '🎵', type: 'text' },
  { key: 'hm_1', label: '1st Honorable Mention', display: 'HM', type: 'track' },
  { key: 'place_10', label: '10th Place', display: '#10', type: 'track' },
  { key: 'place_9', label: '9th Place', display: '#9', type: 'track' },
  { key: 'place_8', label: '8th Place', display: '#8', type: 'track' },
  { key: 'place_7', label: '7th Place', display: '#7', type: 'track' },
  { key: 'hm_2', label: '2nd Honorable Mention', display: 'HM', type: 'track' },
  { key: 'place_6', label: '6th Place', display: '#6', type: 'track' },
  { key: 'place_5', label: '5th Place', display: '#5', type: 'track' },
  { key: 'place_4', label: '4th Place', display: '#4', type: 'track' },
  { key: 'hm_3', label: '3rd Honorable Mention', display: 'HM', type: 'track' },
  { key: 'place_3', label: '3rd Place', display: '#3', type: 'track' },
  { key: 'place_2', label: '2nd Place', display: '#2', type: 'track' },
  { key: 'golden', label: 'Golden Trophy Award', display: '🥇', type: 'track' },
  { key: 'place_1', label: '1st Place', display: '#1', type: 'track' },
  { key: 'summary', label: 'Summary', display: '📝', type: 'text' },
]

interface SearchResult {
  id: string
  name: string
  artist: string
  album_cover: string
  preview_url: string
}

interface TrackData {
  id: string
  name: string
  artist: string
  album_cover: string
  preview_url: string
}

interface SlotState {
  track: TrackData | null
  description: string
}

export default function RecapPage() {
  const currentYear = new Date().getFullYear()
  const [recapType, setRecapType] = useState<'half_year' | 'full_year'>('half_year')
  const [year, setYear] = useState(currentYear)
  const [slotStates, setSlotStates] = useState<Record<string, SlotState>>({})
  const [searchingSlot, setSearchingSlot] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [searchResults, setSearchResults] = useState<SearchResult[]>([])
  const [searching, setSearching] = useState(false)
  const [saving, setSaving] = useState(false)
  const [saved, setSaved] = useState(false)
  const [playingUrl, setPlayingUrl] = useState<string | null>(null)
  const audioRef = useRef<HTMLAudioElement | null>(null)
  const searchTimeout = useRef<ReturnType<typeof setTimeout> | null>(null)

  const period = recapType === 'half_year' ? `${year}-H1` : `${year}-FULL`
  const slots = recapType === 'half_year' ? HALF_YEAR_SLOTS : FULL_YEAR_SLOTS

  useEffect(() => {
    fetchRecap()
  }, [period])

  async function fetchRecap() {
    try {
      const res = await api.get(`/api/recap?period=${period}`)
      const states: Record<string, SlotState> = {}
      for (const entry of res.data.tracks || []) {
        states[entry.slot_key] = {
          track: entry.track_id ? {
            id: entry.track_id,
            name: entry.track_name,
            artist: entry.artist_name,
            album_cover: entry.album_cover_url,
            preview_url: entry.preview_url,
          } : null,
          description: entry.description || '',
        }
      }
      setSlotStates(states)
      setSaved(Object.keys(states).length > 0)
    } catch {
      setSlotStates({})
    }
  }

  function handleSearchInput(q: string) {
    setSearchQuery(q)
    if (searchTimeout.current) clearTimeout(searchTimeout.current)
    searchTimeout.current = setTimeout(() => searchTracks(q), 400)
  }

  async function searchTracks(q: string) {
    if (!q.trim()) { setSearchResults([]); return }
    setSearching(true)
    try {
      const res = await api.get(`/api/spotify/search?q=${encodeURIComponent(q)}`)
      setSearchResults(res.data.tracks || [])
    } catch { setSearchResults([]) }
    setSearching(false)
  }

  function pickTrack(slotKey: string, track: SearchResult) {
    setSlotStates(prev => ({
      ...prev,
      [slotKey]: { track, description: prev[slotKey]?.description || '' }
    }))
    setSearchingSlot(null)
    setSearchQuery('')
    setSearchResults([])
    setSaved(false)
  }

  function updateDescription(slotKey: string, desc: string) {
    setSlotStates(prev => ({
      ...prev,
      [slotKey]: { track: prev[slotKey]?.track || null, description: desc }
    }))
    setSaved(false)
  }

  function clearSlot(slotKey: string) {
    setSlotStates(prev => {
      const next = { ...prev }
      delete next[slotKey]
      return next
    })
    setSaved(false)
  }

  async function saveAll() {
    setSaving(true)
    try {
      for (const slot of slots) {
        const state = slotStates[slot.key]
        if (!state) continue
        // For text-only slots (overall, summary), save with empty track
        if (slot.type === 'text' && state.description) {
          await api.post('/api/recap', {
            period,
            recap_type: recapType,
            slot_key: slot.key,
            track_id: '',
            track_name: '',
            artist_name: '',
            album_cover_url: '',
            preview_url: '',
            description: state.description,
          })
        } else if (slot.type === 'track' && state.track) {
          await api.post('/api/recap', {
            period,
            recap_type: recapType,
            slot_key: slot.key,
            track_id: state.track.id,
            track_name: state.track.name,
            artist_name: state.track.artist,
            album_cover_url: state.track.album_cover,
            preview_url: state.track.preview_url || '',
            description: state.description || '',
          })
        }
      }
      setSaved(true)
    } catch (err) {
      console.error('Failed to save:', err)
    }
    setSaving(false)
  }

  function togglePlay(url: string) {
    if (!url) return
    if (playingUrl === url) {
      audioRef.current?.pause()
      setPlayingUrl(null)
    } else {
      if (audioRef.current) audioRef.current.pause()
      const audio = new Audio(url)
      audio.play()
      audio.onended = () => setPlayingUrl(null)
      audioRef.current = audio
      setPlayingUrl(url)
    }
  }

  function getPeriodTitle() {
    if (recapType === 'half_year') return `My Top 5 — ${year} First Half`
    return `My Top 10 — ${year} Full Year`
  }

  function getPeriodSubtitle() {
    if (recapType === 'half_year') return `January – June ${year} • Your personal picks with stories`
    return `January – December ${year} • Your personal picks with stories`
  }

  return (
    <div className="space-y-5">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">🏆 {getPeriodTitle()}</h1>
          <p className="text-sm text-[var(--text-muted)]">{getPeriodSubtitle()}</p>
        </div>
        <button className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-gradient-to-r from-purple-600 via-red-500 to-orange-400 text-white text-xs font-semibold">
          📷 Share to IG
        </button>
      </div>

      {/* Controls */}
      <div className="flex items-center gap-3 flex-wrap">
        <select
          value={recapType}
          onChange={(e) => setRecapType(e.target.value as 'half_year' | 'full_year')}
          className="px-3 py-2 bg-[var(--bg-card)] border border-[var(--border)] rounded-lg text-sm"
        >
          <option value="half_year">Half-Year (5 อันดับ)</option>
          <option value="full_year">Full-Year (10 อันดับ)</option>
        </select>

        <div className="flex items-center gap-2">
          <button onClick={() => setYear(y => y - 1)} className="px-2 py-1 bg-[var(--bg-card)] border border-[var(--border)] rounded text-sm">←</button>
          <span className="font-semibold">{year}</span>
          <button onClick={() => setYear(y => y + 1)} className="px-2 py-1 bg-[var(--bg-card)] border border-[var(--border)] rounded text-sm">→</button>
        </div>

        <button
          onClick={saveAll}
          disabled={saving}
          className="ml-auto px-5 py-2 bg-[var(--spotify-green)] text-black font-semibold rounded-full text-sm hover:bg-[var(--spotify-green-hover)] transition-colors disabled:opacity-50"
        >
          {saving ? 'Saving...' : saved ? '✓ Saved' : '💾 Save All'}
        </button>
      </div>

      {/* Cards */}
      <div className="space-y-4">
        {slots.map((slot) => {
          const state = slotStates[slot.key]
          const isSearching = searchingSlot === slot.key

          // TEXT-ONLY slot (Overall / Summary)
          if (slot.type === 'text') {
            return (
              <div key={slot.key} className="bg-[var(--bg-card)] border border-[var(--border)] rounded-xl p-5">
                <div className="flex gap-5">
                  <div className="flex-shrink-0 flex items-start justify-center w-14">
                    <span className="text-3xl">{slot.display}</span>
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-xs text-[var(--text-muted)] mb-2 uppercase tracking-wide">{slot.label}</p>
                    <textarea
                      value={state?.description || ''}
                      onChange={(e) => updateDescription(slot.key, e.target.value)}
                      placeholder={slot.key === 'overall' ? 'เขียน overview ของครึ่งปีนี้...' : 'เขียนสรุปความรู้สึกรวมๆ...'}
                      className="w-full px-4 py-3 bg-[var(--bg-main)] border border-[var(--border)] rounded-lg text-sm resize-none focus:outline-none focus:border-[var(--spotify-green)] min-h-[80px]"
                      rows={3}
                    />
                  </div>
                </div>
              </div>
            )
          }

          // TRACK slot
          const hasTrack = !!state?.track
          return (
            <div key={slot.key} className="bg-[var(--bg-card)] border border-[var(--border)] rounded-xl p-5">
              <div className="flex gap-5">
                {/* Rank badge */}
                <div className="flex-shrink-0 flex items-start justify-center w-14">
                  <span className="text-3xl font-black text-[var(--spotify-green)]">
                    {slot.display}
                  </span>
                </div>

                {/* Content */}
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-[var(--text-muted)] mb-2 uppercase tracking-wide">{slot.label}</p>

                  {hasTrack ? (
                    <div className="flex gap-4 items-start">
                      {/* Album cover */}
                      <img
                        src={state.track!.album_cover}
                        alt={state.track!.name}
                        className="w-20 h-20 rounded-lg object-cover shadow-lg flex-shrink-0"
                      />

                      <div className="flex-1 min-w-0">
                        <h3 className="font-bold text-base truncate">{state.track!.name}</h3>
                        <p className="text-sm text-[var(--spotify-green)]">{state.track!.artist}</p>

                        {/* Description */}
                        <div className="mt-2 border-l-2 border-[var(--text-muted)] pl-3">
                          <textarea
                            value={state.description || ''}
                            onChange={(e) => updateDescription(slot.key, e.target.value)}
                            placeholder="เขียนเรื่องราวของเพลงนี้..."
                            className="w-full bg-transparent text-sm text-[var(--text-secondary)] italic resize-none border-none outline-none placeholder:text-[var(--text-muted)]"
                            rows={2}
                          />
                        </div>

                        {/* Preview player */}
                        {state.track!.preview_url && (
                          <div className="mt-3 flex items-center gap-3">
                            <button
                              onClick={() => togglePlay(state.track!.preview_url)}
                              className="w-8 h-8 rounded-full bg-[var(--spotify-green)] flex items-center justify-center text-black flex-shrink-0"
                            >
                              {playingUrl === state.track!.preview_url ? '⏸' : '▶'}
                            </button>
                            <div className="flex-1 h-1 bg-[var(--border)] rounded-full overflow-hidden">
                              <div className={`h-full bg-[var(--spotify-green)] rounded-full transition-all ${
                                playingUrl === state.track!.preview_url ? 'w-1/3' : 'w-0'
                              }`} />
                            </div>
                            <span className="text-xs text-[var(--text-muted)]">0:30</span>
                          </div>
                        )}
                      </div>

                      {/* Remove */}
                      <button
                        onClick={() => clearSlot(slot.key)}
                        className="text-xs text-[var(--text-muted)] hover:text-red-400 flex-shrink-0"
                      >
                        ✕
                      </button>
                    </div>
                  ) : (
                    /* Empty — search */
                    <div>
                      {!isSearching ? (
                        <button
                          onClick={() => setSearchingSlot(slot.key)}
                          className="w-full py-4 border-2 border-dashed border-[var(--border)] rounded-lg text-sm text-[var(--text-muted)] hover:border-[var(--spotify-green)] hover:text-[var(--spotify-green)] transition-colors"
                        >
                          + เลือกเพลง
                        </button>
                      ) : (
                        <div className="space-y-2">
                          <input
                            type="text"
                            value={searchQuery}
                            onChange={(e) => handleSearchInput(e.target.value)}
                            placeholder="🔍 ค้นหาเพลง..."
                            className="w-full px-3 py-2 bg-[var(--bg-main)] border border-[var(--border)] rounded-lg text-sm focus:outline-none focus:border-[var(--spotify-green)]"
                            autoFocus
                          />
                          {(searching || searchResults.length > 0) && (
                            <div className="max-h-40 overflow-y-auto space-y-1 border border-[var(--border)] rounded-lg p-2 bg-[var(--bg-main)]">
                              {searching && <p className="text-xs text-center text-[var(--text-muted)] py-2">Searching...</p>}
                              {searchResults.map((track) => (
                                <button
                                  key={track.id}
                                  onClick={() => pickTrack(slot.key, track)}
                                  className="w-full flex items-center gap-3 p-2 rounded-lg hover:bg-[var(--bg-hover)] text-left transition-colors"
                                >
                                  <img src={track.album_cover} alt="" className="w-9 h-9 rounded object-cover" />
                                  <div className="flex-1 min-w-0">
                                    <p className="text-sm font-medium truncate">{track.name}</p>
                                    <p className="text-xs text-[var(--text-muted)] truncate">{track.artist}</p>
                                  </div>
                                </button>
                              ))}
                            </div>
                          )}
                          <button
                            onClick={() => { setSearchingSlot(null); setSearchQuery(''); setSearchResults([]) }}
                            className="text-xs text-[var(--text-muted)] hover:text-white"
                          >
                            ← Cancel
                          </button>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {/* Bottom save */}
      <div className="sticky bottom-4 flex justify-center">
        <button
          onClick={saveAll}
          disabled={saving}
          className="px-8 py-3 bg-[var(--spotify-green)] text-black font-bold rounded-full text-sm shadow-lg hover:bg-[var(--spotify-green-hover)] transition-all disabled:opacity-50"
        >
          {saving ? 'Saving...' : '💾 Save Recap'}
        </button>
      </div>
    </div>
  )
}
