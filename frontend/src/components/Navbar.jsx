import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

function Navbar() {
    const { user, logout } = useAuth()

    return (
        <nav className="navbar">
            <div className="navbar-content">
                <Link to="/dashboard" className="navbar-brand">
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <polyline points="4 17 10 11 4 5"></polyline>
                        <line x1="12" y1="19" x2="20" y2="19"></line>
                    </svg>
                    SSH Terminal
                </Link>

                {user && (
                    <div className="navbar-user">
                        <div className="user-info">
                            <div className="user-avatar">
                                {user.name?.charAt(0).toUpperCase()}
                            </div>
                            <span>{user.name}</span>
                        </div>
                        <button onClick={logout} className="btn btn-ghost">
                            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
                                <polyline points="16 17 21 12 16 7"></polyline>
                                <line x1="21" y1="12" x2="9" y2="12"></line>
                            </svg>
                            Logout
                        </button>
                    </div>
                )}
            </div>
        </nav>
    )
}

export default Navbar
