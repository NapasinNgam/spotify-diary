export default function HistoryPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">📚 History Bookshelf</h1>
        <p className="text-sm text-[var(--text-muted)]">Your past records archived by category</p>
      </div>

      <div className="space-y-6">
        <div>
          <h3 className="text-xs uppercase tracking-wider text-[var(--text-muted)] font-semibold mb-3">
            📅 Diary Calendars
          </h3>
          <div className="grid grid-cols-3 gap-3">
            <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-lg p-4 cursor-pointer hover:border-[var(--spotify-green)] transition-all">
              <div className="text-xl mb-2">📅</div>
              <h4 className="text-sm font-semibold">No entries yet</h4>
              <p className="text-xs text-[var(--text-muted)]">Start logging to see history</p>
            </div>
          </div>
        </div>

        <div>
          <h3 className="text-xs uppercase tracking-wider text-[var(--text-muted)] font-semibold mb-3">
            📊 Monthly Records
          </h3>
          <div className="grid grid-cols-3 gap-3">
            <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-lg p-4 cursor-pointer hover:border-[var(--spotify-green)] transition-all">
              <div className="text-xl mb-2">📊</div>
              <h4 className="text-sm font-semibold">No records yet</h4>
              <p className="text-xs text-[var(--text-muted)]">Generated on the 1st of each month</p>
            </div>
          </div>
        </div>

        <div>
          <h3 className="text-xs uppercase tracking-wider text-[var(--text-muted)] font-semibold mb-3">
            🏆 Half-Year Recaps
          </h3>
          <div className="grid grid-cols-3 gap-3">
            <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-lg p-4 cursor-pointer hover:border-[var(--spotify-green)] transition-all">
              <div className="text-xl mb-2">🏆</div>
              <h4 className="text-sm font-semibold">No recaps yet</h4>
              <p className="text-xs text-[var(--text-muted)]">Created every 6 months</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
