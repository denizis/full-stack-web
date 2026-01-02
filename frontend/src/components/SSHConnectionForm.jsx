import { useState } from 'react'

function SSHConnectionForm({ connection, onSave, onClose }) {
    const [formData, setFormData] = useState({
        name: connection?.name || '',
        host: connection?.host || '',
        port: connection?.port || 22,
        username: connection?.username || '',
        password: '',
        private_key: '',
        auth_type: connection?.auth_type || 'password'
    })
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState('')

    const handleChange = (e) => {
        const { name, value } = e.target
        setFormData(prev => ({
            ...prev,
            [name]: name === 'port' ? parseInt(value) || 22 : value
        }))
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        setError('')
        setLoading(true)

        try {
            await onSave(formData)
        } catch (err) {
            setError(err.response?.data || 'Failed to save connection')
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="modal-overlay" onClick={onClose}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
                <div className="modal-header">
                    <h2 className="modal-title">
                        {connection ? 'Edit Connection' : 'New SSH Connection'}
                    </h2>
                    <button onClick={onClose} className="modal-close">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <line x1="18" y1="6" x2="6" y2="18"></line>
                            <line x1="6" y1="6" x2="18" y2="18"></line>
                        </svg>
                    </button>
                </div>

                {error && (
                    <div className="alert alert-error">
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label className="form-label">Connection Name</label>
                        <input
                            type="text"
                            name="name"
                            className="form-input"
                            placeholder="My Server"
                            value={formData.name}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '1rem' }}>
                        <div className="form-group">
                            <label className="form-label">Host</label>
                            <input
                                type="text"
                                name="host"
                                className="form-input"
                                placeholder="192.168.1.1 or example.com"
                                value={formData.host}
                                onChange={handleChange}
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label className="form-label">Port</label>
                            <input
                                type="number"
                                name="port"
                                className="form-input"
                                value={formData.port}
                                onChange={handleChange}
                                min="1"
                                max="65535"
                            />
                        </div>
                    </div>

                    <div className="form-group">
                        <label className="form-label">Username</label>
                        <input
                            type="text"
                            name="username"
                            className="form-input"
                            placeholder="root"
                            value={formData.username}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label className="form-label">Authentication Type</label>
                        <select
                            name="auth_type"
                            className="form-input"
                            value={formData.auth_type}
                            onChange={handleChange}
                        >
                            <option value="password">Password</option>
                            <option value="key">Private Key</option>
                        </select>
                    </div>

                    {formData.auth_type === 'password' ? (
                        <div className="form-group">
                            <label className="form-label">
                                Password {connection && '(leave blank to keep current)'}
                            </label>
                            <input
                                type="password"
                                name="password"
                                className="form-input"
                                placeholder="••••••••"
                                value={formData.password}
                                onChange={handleChange}
                                required={!connection}
                            />
                        </div>
                    ) : (
                        <div className="form-group">
                            <label className="form-label">
                                Private Key {connection && '(leave blank to keep current)'}
                            </label>
                            <textarea
                                name="private_key"
                                className="form-input"
                                placeholder="-----BEGIN RSA PRIVATE KEY-----&#10;...&#10;-----END RSA PRIVATE KEY-----"
                                value={formData.private_key}
                                onChange={handleChange}
                                rows={6}
                                style={{ fontFamily: 'var(--font-mono)', fontSize: '0.75rem' }}
                                required={!connection}
                            />
                        </div>
                    )}

                    <div style={{ display: 'flex', gap: '1rem', marginTop: '1.5rem' }}>
                        <button type="button" onClick={onClose} className="btn btn-secondary" style={{ flex: 1 }}>
                            Cancel
                        </button>
                        <button type="submit" className="btn btn-primary" style={{ flex: 1 }} disabled={loading}>
                            {loading ? <span className="spinner"></span> : (connection ? 'Update' : 'Create')}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    )
}

export default SSHConnectionForm
