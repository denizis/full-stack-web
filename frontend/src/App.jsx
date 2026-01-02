import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from './context/AuthContext'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import Terminal from './pages/Terminal'
import AuthCallback from './pages/AuthCallback'
import Navbar from './components/Navbar'

function ProtectedRoute({ children }) {
    const { user, loading } = useAuth()

    if (loading) {
        return (
            <div className="auth-page">
                <div className="spinner" style={{ width: 40, height: 40 }}></div>
            </div>
        )
    }

    if (!user) {
        return <Navigate to="/login" replace />
    }

    return children
}

function PublicRoute({ children }) {
    const { user, loading } = useAuth()

    if (loading) {
        return (
            <div className="auth-page">
                <div className="spinner" style={{ width: 40, height: 40 }}></div>
            </div>
        )
    }

    if (user) {
        return <Navigate to="/dashboard" replace />
    }

    return children
}

function App() {
    const { user } = useAuth()

    return (
        <div className="app-container">
            {user && <Navbar />}
            <Routes>
                <Route path="/login" element={
                    <PublicRoute>
                        <Login />
                    </PublicRoute>
                } />
                <Route path="/register" element={
                    <PublicRoute>
                        <Register />
                    </PublicRoute>
                } />
                <Route path="/auth/callback" element={<AuthCallback />} />
                <Route path="/dashboard" element={
                    <ProtectedRoute>
                        <Dashboard />
                    </ProtectedRoute>
                } />
                <Route path="/terminal/:id" element={
                    <ProtectedRoute>
                        <Terminal />
                    </ProtectedRoute>
                } />
                <Route path="/" element={<Navigate to="/dashboard" replace />} />
            </Routes>
        </div>
    )
}

export default App
