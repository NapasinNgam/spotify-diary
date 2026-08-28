export default function MonthlyPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">📊 Monthly Records</h1>
          <p className="text-sm text-[var(--text-muted)]">Top 10 most played from Spotify history</p>
        </div>
        <button className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-gradient-to-r from-purple-600 via-red-500 to-orange-400 text-white text-xs font-semibold">
          📷 Share to IG
        </button>
      </div>

      <div className="bg-[var(--bg-card)] border border-[var(--border)] rounded-xl p-6">
        <p className="text-[var(--text-muted)] text-center py-12">
          Monthly records will appear here after your first month of syncing.
        </p>
      </div>
    </div>
  )
}
