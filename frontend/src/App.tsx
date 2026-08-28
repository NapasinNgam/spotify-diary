import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from './hooks/useAuth'
import Layout from './components/layout/Layout'
import LoginPage from './pages/LoginPage'
import NewsPage from './pages/NewsPage'
import CalendarPage from './pages/CalendarPage'
import MonthlyPage from './pages/MonthlyPage'
import RecapPage from './pages/RecapPage'
import HistoryPage from './pages/HistoryPage'
import AuthCallback from './pages/AuthCallback'

function App() {
  const { isAuthenticated } = useAuth()

  if (!isAuthenticated) {
    return (
      <Routes>
        <Route path="/auth/callback" element={<AuthCallback />} />
        <Route path="*" element={<LoginPage />} />
      </Routes>
    )
  }

  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<Navigate to="/news" replace />} />
        <Route path="news" element={<NewsPage />} />
        <Route path="calendar" element={<CalendarPage />} />
        <Route path="monthly" element={<MonthlyPage />} />
        <Route path="recap" element={<RecapPage />} />
        <Route path="history" element={<HistoryPage />} />
      </Route>
      <Route path="/auth/callback" element={<AuthCallback />} />
    </Routes>
  )
}

export default App
