import { useState, useEffect, useRef, useLayoutEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { WebLinksAddon } from 'xterm-addon-web-links'
import 'xterm/css/xterm.css'
import { sshApi } from '../services/api'

function TerminalPage() {
    const { id } = useParams()
    const navigate = useNavigate()
    const terminalRef = useRef(null)
    const termRef = useRef(null)
    const fitAddonRef = useRef(null)
    const wsRef = useRef(null)
    const [connection, setConnection] = useState(null)
    const [status, setStatus] = useState('connecting')
    const [error, setError] = useState('')
    const [isTerminalReady, setIsTerminalReady] = useState(false)

    // Fetch connection info
    useEffect(() => {
        sshApi.get(id)
            .then(res => setConnection(res.data))
            .catch(() => {
                setError('Connection not found')
                navigate('/dashboard')
            })
    }, [id, navigate])

    // Initialize terminal with useLayoutEffect for immediate DOM access
    useLayoutEffect(() => {
        if (!terminalRef.current) {
            console.log('terminalRef.current is null')
            return
        }

        if (termRef.current) {
            console.log('Terminal already created')
            return
        }

        console.log('Creating xterm instance...')

        const term = new Terminal({
            theme: {
                background: '#0a0a0f',
                foreground: '#e4e4e7',
                cursor: '#10b981',
                cursorAccent: '#10b981',
                selection: 'rgba(16, 185, 129, 0.3)',
                black: '#18181b',
                red: '#ef4444',
                green: '#10b981',
                yellow: '#f59e0b',
                blue: '#3b82f6',
                magenta: '#a855f7',
                cyan: '#06b6d4',
                white: '#e4e4e7',
                brightBlack: '#52525b',
                brightRed: '#f87171',
                brightGreen: '#34d399',
                brightYellow: '#fbbf24',
                brightBlue: '#60a5fa',
                brightMagenta: '#c084fc',
                brightCyan: '#22d3ee',
                brightWhite: '#ffffff'
            },
            fontFamily: '"JetBrains Mono", "Fira Code", Consolas, monospace',
            fontSize: 14,
            lineHeight: 1.2,
            cursorBlink: true,
            cursorStyle: 'block',
            scrollback: 10000,
            convertEol: true
        })

        const fitAddon = new FitAddon()
        term.loadAddon(fitAddon)
        term.loadAddon(new WebLinksAddon())

        console.log('Opening terminal in container...')
        term.open(terminalRef.current)

        // Wait a tiny bit for DOM to settle, then fit
        setTimeout(() => {
            fitAddon.fit()
            console.log('Terminal fitted, cols:', term.cols, 'rows:', term.rows)
        }, 50)

        termRef.current = term
        fitAddonRef.current = fitAddon

        // Handle input
        term.onData((data) => {
            if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
                wsRef.current.send(JSON.stringify({ type: 'input', data }))
            }
        })

        // Handle window resize
        const handleResize = () => {
            fitAddon.fit()
            if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
                wsRef.current.send(JSON.stringify({
                    type: 'resize',
                    cols: term.cols,
                    rows: term.rows
                }))
            }
        }

        window.addEventListener('resize', handleResize)

        setIsTerminalReady(true)
        console.log('Terminal ready!')

        return () => {
            window.removeEventListener('resize', handleResize)
            term.dispose()
            termRef.current = null
        }
    }, [])

    // Connect WebSocket when terminal is ready
    useEffect(() => {
        if (!isTerminalReady) return

        const connect = () => {
            if (wsRef.current) {
                wsRef.current.close()
            }

            const token = localStorage.getItem('token')
            const wsUrl = `ws://localhost:3000/ws/terminal/${id}?token=${token}`

            console.log('Connecting WebSocket to:', wsUrl)
            setStatus('connecting')
            setError('')

            const ws = new WebSocket(wsUrl)
            wsRef.current = ws

            ws.onopen = () => {
                console.log('WebSocket opened')
                setStatus('connected')

                if (fitAddonRef.current) {
                    fitAddonRef.current.fit()
                }

                if (termRef.current) {
                    ws.send(JSON.stringify({
                        type: 'resize',
                        cols: termRef.current.cols,
                        rows: termRef.current.rows
                    }))
                }
            }

            ws.onmessage = (event) => {
                console.log('WS message:', event.data.substring(0, 100))
                try {
                    const data = JSON.parse(event.data)
                    if (data.type === 'output' && termRef.current) {
                        termRef.current.write(data.data)
                    } else if (data.type === 'error') {
                        setError(data.data)
                        setStatus('disconnected')
                    }
                } catch (e) {
                    if (termRef.current) {
                        termRef.current.write(event.data)
                    }
                }
            }

            ws.onclose = () => {
                console.log('WebSocket closed')
                setStatus('disconnected')
            }

            ws.onerror = (e) => {
                console.error('WebSocket error:', e)
                setError('WebSocket connection failed')
                setStatus('disconnected')
            }
        }

        connect()

        return () => {
            wsRef.current?.close()
        }
    }, [id, isTerminalReady])

    const handleDisconnect = () => {
        wsRef.current?.close()
        navigate('/dashboard')
    }

    const handleReconnect = () => {
        if (termRef.current) {
            termRef.current.clear()
        }

        const token = localStorage.getItem('token')
        const wsUrl = `ws://localhost:3000/ws/terminal/${id}?token=${token}`

        setStatus('connecting')
        setError('')

        const ws = new WebSocket(wsUrl)
        wsRef.current = ws

        ws.onopen = () => {
            setStatus('connected')
            if (fitAddonRef.current) fitAddonRef.current.fit()
            if (termRef.current) {
                ws.send(JSON.stringify({
                    type: 'resize',
                    cols: termRef.current.cols,
                    rows: termRef.current.rows
                }))
            }
        }

        ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data)
                if (data.type === 'output' && termRef.current) {
                    termRef.current.write(data.data)
                } else if (data.type === 'error') {
                    setError(data.data)
                    setStatus('disconnected')
                }
            } catch (e) {
                if (termRef.current) {
                    termRef.current.write(event.data)
                }
            }
        }

        ws.onclose = () => setStatus('disconnected')
        ws.onerror = () => {
            setError('WebSocket connection failed')
            setStatus('disconnected')
        }
    }

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: 'calc(100vh - 60px)' }}>
            <div className="terminal-container" style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
                <div className="terminal-header">
                    <div className="terminal-title">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <polyline points="4 17 10 11 4 5"></polyline>
                            <line x1="12" y1="19" x2="20" y2="19"></line>
                        </svg>
                        {connection?.name || 'Terminal'} • {connection?.username}@{connection?.host}
                    </div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                        <div className={`terminal-status ${status}`}>
                            <span className={`status-dot ${status}`}></span>
                            {status === 'connecting' && 'Connecting...'}
                            {status === 'connected' && 'Connected'}
                            {status === 'disconnected' && 'Disconnected'}
                        </div>
                        {status === 'disconnected' && (
                            <button onClick={handleReconnect} className="btn btn-primary" style={{ padding: '0.375rem 0.75rem', fontSize: '0.75rem' }}>
                                Reconnect
                            </button>
                        )}
                        <button onClick={handleDisconnect} className="btn btn-ghost" style={{ padding: '0.375rem 0.75rem', fontSize: '0.75rem' }}>
                            Close
                        </button>
                    </div>
                </div>

                {error && (
                    <div className="alert alert-error" style={{ margin: '0.5rem', borderRadius: '4px' }}>
                        {error}
                    </div>
                )}

                <div
                    className="terminal-body"
                    ref={terminalRef}
                    style={{ flex: 1, minHeight: '400px' }}
                ></div>
            </div>
        </div>
    )
}

export default TerminalPage
