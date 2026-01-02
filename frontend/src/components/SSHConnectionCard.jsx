function SSHConnectionCard({ connection, onConnect, onEdit, onDelete }) {

    const getGradient = (id) => {
        const hues = [
            ['#10b981', '#059669'],
            ['#6366f1', '#4f46e5'],
            ['#ec4899', '#db2777'],
            ['#f59e0b', '#d97706'],
            ['#3b82f6', '#2563eb'],
        ];
        const [c1, c2] = hues[id % hues.length];
        return `linear-gradient(135deg, ${c1}, ${c2})`;
    }

    return (
        <div className="connection-card" onClick={onConnect}>
            <div className="connection-info">
                <div
                    className="connection-icon"
                    style={{ background: getGradient(connection.id), color: 'white', border: 'none' }}
                >
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <polyline points="4 17 10 11 4 5"></polyline>
                        <line x1="12" y1="19" x2="20" y2="19"></line>
                    </svg>
                </div>
                <div className="connection-details">
                    <h4>{connection.name}</h4>
                    <p>{connection.username}@{connection.host}</p>
                </div>
            </div>

            <div className="connection-actions" onClick={(e) => e.stopPropagation()}>
                <button onClick={onConnect} className="btn btn-primary btn-icon" title="Connect" style={{ width: 36, height: 36, padding: 0, borderRadius: '50%' }}>
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                        <polygon points="5 3 19 12 5 21 5 3"></polygon>
                    </svg>
                </button>
                <div style={{ width: 1, height: 24, background: 'rgba(255,255,255,0.1)', margin: '0 4px' }}></div>
                <button onClick={onEdit} className="btn btn-secondary btn-icon" title="Edit" style={{ width: 36, height: 36, padding: 0, borderRadius: '50%' }}>
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                    </svg>
                </button>
                <button onClick={onDelete} className="btn btn-danger btn-icon" title="Delete" style={{ width: 36, height: 36, padding: 0, borderRadius: '50%' }}>
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <polyline points="3 6 5 6 21 6"></polyline>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                    </svg>
                </button>
            </div>
        </div>
    )
}

export default SSHConnectionCard
