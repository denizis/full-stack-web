import { createContext, useContext, useState, useEffect } from 'react'
import api from '../services/api'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
    const [user, setUser] = useState(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        checkAuth()
    }, [])

    const checkAuth = async () => {
        const token = localStorage.getItem('token')
        if (!token) {
            setLoading(false)
            return
        }

        try {
            const response = await api.get('/auth/me')
            setUser(response.data)
        } catch (err) {
            localStorage.removeItem('token')
        } finally {
            setLoading(false)
        }
    }

    const login = async (email, password) => {
        setError(null)
        try {
            const response = await api.post('/auth/login', { email, password })
            localStorage.setItem('token', response.data.token)
            setUser(response.data.user)
            return response.data
        } catch (err) {
            const message = err.response?.data || 'Login failed'
            setError(message)
            throw new Error(message)
        }
    }

    const register = async (email, password, name) => {
        setError(null)
        try {
            const response = await api.post('/auth/register', { email, password, name })
            localStorage.setItem('token', response.data.token)
            setUser(response.data.user)
            return response.data
        } catch (err) {
            const message = err.response?.data || 'Registration failed'
            setError(message)
            throw new Error(message)
        }
    }

    const loginWithGoogle = (mode = 'login') => {
        window.location.href = `/api/auth/google?mode=${mode}`
    }

    const handleGoogleCallback = async (token) => {
        localStorage.setItem('token', token)
        await checkAuth()
    }

    const logout = () => {
        localStorage.removeItem('token')
        setUser(null)
    }

    return (
        <AuthContext.Provider value={{
            user,
            loading,
            error,
            login,
            register,
            loginWithGoogle,
            handleGoogleCallback,
            logout
        }}>
            {children}
        </AuthContext.Provider>
    )
}

export function useAuth() {
    const context = useContext(AuthContext)
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider')
    }
    return context
}
