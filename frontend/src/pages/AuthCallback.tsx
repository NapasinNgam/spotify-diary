import { useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

export default function AuthCallback() {
  const [searchParams] = useSearchParams()
  const { setToken } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    const token = searchParams.get('token')
    if (token) {
      setToken(token)
      navigate('/news', { replace: true })
    } else {
      navigate('/', { replace: true })
    }
  }, [searchParams, setToken, navigate])

  return (
    <div className="min-h-screen flex items-center justify-center bg-[var(--bg-primary)]">
      <div className="text-center space-y-4">
        <div className="animate-spin w-8 h-8 border-2 border-[var(--spotify-green)] border-t-transparent rounded-full mx-auto"></div>
        <p className="text-[var(--text-secondary)]">Logging you in...</p>
      </div>
    </div>
  )
}
