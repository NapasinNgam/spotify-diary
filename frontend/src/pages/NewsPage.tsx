export default function NewsPage() {
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

      {/* Stats will be populated from API */}
      <div className="grid grid-cols-3 gap-3">
        <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-lg p-4 text-center">
          <div className="text-2xl font-bold text-[var(--spotify-green)] font-mono">--</div>
          <div className="text-xs text-[var(--text-muted)] mt-1">Tracks Played</div>
        </div>
        <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-lg p-4 text-center">
          <div className="text-2xl font-bold text-[var(--spotify-green)] font-mono">--</div>
          <div className="text-xs text-[var(--text-muted)] mt-1">Listening Time</div>
        </div>
        <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-lg p-4 text-center">
          <div className="text-2xl font-bold text-[var(--spotify-green)] font-mono">--</div>
          <div className="text-xs text-[var(--text-muted)] mt-1">Artists</div>
        </div>
      </div>

      <p className="text-[var(--text-muted)] text-center py-8">
        Start syncing your listening history to see stats here.
      </p>
    </div>
  )
}
