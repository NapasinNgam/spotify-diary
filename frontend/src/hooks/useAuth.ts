import { create } from 'zustand'
import { api } from '../lib/api'

interface User {
  id: number
  spotify_id: string
  display_name: string
  email: string
  avatar_url: string
}

interface AuthState {
  token: string | null
  user: User | null
  isAuthenticated: boolean
  login: () => Promise<void>
  logout: () => void
  setToken: (token: string) => void
  fetchUser: () => Promise<void>
}

export const useAuth = create<AuthState>((set, get) => ({
  token: localStorage.getItem('token'),
  user: null,
  isAuthenticated: !!localStorage.getItem('token'),

  login: async () => {
    try {
      const response = await api.get('/api/auth/login')
      window.location.href = response.data.url
    } catch (error) {
      console.error('Login failed:', error)
    }
  },

  logout: () => {
    localStorage.removeItem('token')
    set({ token: null, user: null, isAuthenticated: false })
  },

  setToken: (token: string) => {
    localStorage.setItem('token', token)
    set({ token, isAuthenticated: true })
    get().fetchUser()
  },

  fetchUser: async () => {
    try {
      const response = await api.get('/api/me')
      set({ user: response.data.user })
    } catch (error) {
      console.error('Failed to fetch user:', error)
      get().logout()
    }
  },
}))
