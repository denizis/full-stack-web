import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { sshApi } from '../services/api'
import SSHConnectionCard from '../components/SSHConnectionCard'
import SSHConnectionForm from '../components/SSHConnectionForm'

function Dashboard() {
    const navigate = useNavigate()
    const [connections, setConnections] = useState([])
    const [loading, setLoading] = useState(true)
    const [showForm, setShowForm] = useState(false)
    const [editConnection, setEditConnection] = useState(null)
    const [error, setError] = useState('')

    useEffect(() => {
        fetchConnections()
    }, [])

    const fetchConnections = async () => {
        try {
            const response = await sshApi.list()
            setConnections(response.data || [])
        } catch (err) {
            setError('Failed to load connections')
        } finally {
            setLoading(false)
        }
    }

    const handleConnect = (id) => {
        navigate(`/terminal/${id}`)
    }

    const handleEdit = (connection) => {
        setEditConnection(connection)
        setShowForm(true)
    }

    const handleDelete = async (id) => {
        if (!confirm('Are you sure you want to delete this connection?')) return

        try {
            await sshApi.delete(id)
            setConnections(connections.filter(c => c.id !== id))
        } catch (err) {
            setError('Failed to delete connection')
        }
    }

    const handleSave = async (data) => {
        try {
            if (editConnection) {
                await sshApi.update(editConnection.id, data)
            } else {
                await sshApi.create(data)
            }
            setShowForm(false)
            setEditConnection(null)
            fetchConnections()
        } catch (err) {
            throw err
        }
    }

    const handleCloseForm = () => {
        setShowForm(false)
        setEditConnection(null)
    }

    return (
        <main className="main-content">
            <div className="dashboard-header">
                <h1 className="dashboard-title">SSH Connections</h1>
                <button onClick={() => setShowForm(true)} className="btn btn-primary">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <line x1="12" y1="5" x2="12" y2="19"></line>
                        <line x1="5" y1="12" x2="19" y2="12"></line>
                    </svg>
                    Add Connection
                </button>
            </div>

            {error && (
                <div className="alert alert-error">
                    {error}
                    <button onClick={() => setError('')} style={{ marginLeft: 'auto', background: 'none', border: 'none', color: 'inherit', cursor: 'pointer' }}>×</button>
                </div>
            )}

            {loading ? (
                <div style={{ display: 'flex', justifyContent: 'center', padding: '4rem' }}>
                    <div className="spinner" style={{ width: 40, height: 40 }}></div>
                </div>
            ) : connections.length === 0 ? (
                <div className="empty-state">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                        <rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect>
                        <line x1="8" y1="21" x2="16" y2="21"></line>
                        <line x1="12" y1="17" x2="12" y2="21"></line>
                    </svg>
                    <h3>No SSH connections yet</h3>
                    <p>Add your first SSH connection to get started</p>
                    <button onClick={() => setShowForm(true)} className="btn btn-primary" style={{ marginTop: '1rem' }}>
                        Add Connection
                    </button>
                </div>
            ) : (
                <div className="connections-grid">
                    {connections.map((connection, index) => (
                        <div
                            key={connection.id}
                            style={{
                                animation: `fadeIn 0.5s ease-out forwards`,
                                animationDelay: `${index * 0.1}s`,
                                opacity: 0
                            }}
                        >
                            <SSHConnectionCard
                                connection={connection}
                                onConnect={() => handleConnect(connection.id)}
                                onEdit={() => handleEdit(connection)}
                                onDelete={() => handleDelete(connection.id)}
                            />
                        </div>
                    ))}
                </div>
            )}

            {showForm && (
                <SSHConnectionForm
                    connection={editConnection}
                    onSave={handleSave}
                    onClose={handleCloseForm}
                />
            )}
        </main>
    )
}

export default Dashboard
