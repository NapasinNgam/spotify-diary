import { Outlet, NavLink } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'

const navItems = [
  { path: '/news', label: '📰 News' },
  { path: '/calendar', label: '📅 Calendar' },
  { path: '/monthly', label: '📊 Monthly' },
  { path: '/recap', label: '🏆 Recap' },
  { path: '/history', label: '📚 History' },
]

export default function Layout() {
  const { user, logout } = useAuth()

  return (
    <div className="min-h-screen bg-[var(--bg-primary)]">
      {/* Top Navigation */}
      <nav className="fixed top-0 left-0 right-0 bg-[rgba(9,9,11,0.85)] backdrop-blur-xl border-b border-[var(--border)] z-50 px-6 h-14 flex items-center justify-between">
        <div className="flex items-center gap-3 font-bold text-base">
          <span className="text-[var(--spotify-green)] text-xl">♪</span>
          <span>Music Diary</span>
        </div>

        <div className="flex gap-1 bg-[var(--bg-card)] p-1 rounded-lg">
          {navItems.map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) =>
                `px-3 py-1.5 rounded-md text-xs font-medium transition-all ${
                  isActive
                    ? 'bg-[var(--bg-hover)] text-[var(--text-primary)]'
                    : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)]'
                }`
              }
            >
              {item.label}
            </NavLink>
          ))}
        </div>

        <div className="flex items-center gap-3">
          {user && (
            <span className="text-xs text-[var(--text-secondary)]">
              {user.display_name}
            </span>
          )}
          <button
            onClick={logout}
            className="text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-colors"
          >
            Logout
          </button>
        </div>
      </nav>

      {/* Main Content */}
      <main className="pt-14 max-w-5xl mx-auto px-6 py-8">
        <Outlet />
      </main>
    </div>
  )
}
