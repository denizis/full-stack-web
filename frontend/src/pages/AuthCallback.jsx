import { useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

function AuthCallback() {
    const navigate = useNavigate()
    const [searchParams] = useSearchParams()
    const { handleGoogleCallback } = useAuth()

    useEffect(() => {
        const token = searchParams.get('token')
        if (token) {
            handleGoogleCallback(token).then(() => {
                navigate('/dashboard')
            })
        } else {
            navigate('/login')
        }
    }, [searchParams, handleGoogleCallback, navigate])

    return (
        <div className="auth-page">
            <div className="spinner" style={{ width: 40, height: 40 }}></div>
        </div>
    )
}

export default AuthCallback
