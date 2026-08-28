export default function CalendarPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">📅 Diary Calendar</h1>
          <p className="text-sm text-[var(--text-muted)]">Log 1 song per day</p>
        </div>
        <button className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-gradient-to-r from-purple-600 via-red-500 to-orange-400 text-white text-xs font-semibold">
          📷 Share to IG
        </button>
      </div>

      {/* Calendar grid placeholder */}
      <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-xl p-6">
        <p className="text-[var(--text-muted)] text-center py-12">
          Calendar view will be rendered here.<br/>
          Connect your Spotify account to start logging songs.
        </p>
      </div>
    </div>
  )
}
