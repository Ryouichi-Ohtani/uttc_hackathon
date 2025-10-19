import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { User } from '@/types'
import { authService } from '@/services/auth'
import toast from 'react-hot-toast'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  loading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (data: {
    email: string
    username: string
    password: string
    display_name: string
  }) => Promise<void>
  logout: () => void
  setUser: (user: User) => void
}

// Initialize from localStorage immediately
const getInitialState = () => {
  if (typeof window === 'undefined') {
    return { user: null, token: null, isAuthenticated: false }
  }

  const token = localStorage.getItem('auth_token')
  const storedState = localStorage.getItem('auth-storage')

  if (token && storedState) {
    try {
      const parsed = JSON.parse(storedState)
      if (parsed.state?.user && parsed.state?.token) {
        return {
          user: parsed.state.user,
          token: parsed.state.token,
          isAuthenticated: true,
        }
      }
    } catch (e) {
      console.error('Failed to parse auth storage:', e)
    }
  }

  return {
    user: null,
    token: null,
    isAuthenticated: false,
  }
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => {
      const initial = getInitialState()
      return {
        ...initial,
        loading: false,

        login: async (email, password) => {
          console.log('AuthStore: Starting login...')
          set({ loading: true })
          try {
            const response = await authService.login({ email, password })
            console.log('AuthStore: Login response received:', { user: response.user, hasToken: !!response.token })
            localStorage.setItem('auth_token', response.token)
            console.log('AuthStore: Token saved to localStorage')
            set({
              user: response.user,
              token: response.token,
              isAuthenticated: true,
              loading: false,
            })
            console.log('AuthStore: State updated, isAuthenticated:', true)
            toast.success('Login successful!')
          } catch (error: any) {
            console.error('AuthStore: Login failed:', error)
            console.error('Error response:', error.response)
            set({ loading: false })

            // Improved error message handling
            let errorMessage = 'Invalid email or password'
            if (error.response?.data?.error?.message) {
              errorMessage = error.response.data.error.message
            } else if (error.response?.data?.message) {
              errorMessage = error.response.data.message
            } else if (error.message) {
              errorMessage = error.message
            }

            toast.error(errorMessage)
            throw error
          }
        },

        register: async (data) => {
          console.log('AuthStore: Starting registration...')
          set({ loading: true })
          try {
            const response = await authService.register(data)
            console.log('AuthStore: Registration response received:', { user: response.user, hasToken: !!response.token })
            localStorage.setItem('auth_token', response.token)
            console.log('AuthStore: Token saved to localStorage')
            set({
              user: response.user,
              token: response.token,
              isAuthenticated: true,
              loading: false,
            })
            console.log('AuthStore: State updated, isAuthenticated:', true)
            toast.success('Registration successful!')
          } catch (error: any) {
            console.error('AuthStore: Registration failed:', error)
            console.error('Error response:', error.response)
            set({ loading: false })

            // Improved error message handling
            let errorMessage = 'Registration failed. Please try again.'
            if (error.response?.data?.error?.message) {
              errorMessage = error.response.data.error.message
            } else if (error.response?.data?.message) {
              errorMessage = error.response.data.message
            } else if (error.message) {
              errorMessage = error.message
            }

            toast.error(errorMessage)
            throw error
          }
        },

        logout: () => {
          localStorage.removeItem('auth_token')
          set({
            user: null,
            token: null,
            isAuthenticated: false,
          })
          toast.success('Logged out successfully')
        },

        setUser: (user) => {
          set({ user })
        },
      }
    },
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
