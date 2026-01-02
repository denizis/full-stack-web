package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type TerminalService interface {
	StartSession(ws *websocket.Conn, connID uint, userID uint) error
}

type terminalService struct {
	sshService SSHService
}

func NewTerminalService(sshService SSHService) TerminalService {
	return &terminalService{
		sshService: sshService,
	}
}

func (s *terminalService) StartSession(ws *websocket.Conn, connID uint, userID uint) error {
	// Get credentials and connection info
	password, privateKey, conn, err := s.sshService.GetDecryptedCredentials(connID, userID)
	if err != nil {
		return fmt.Errorf("failed to get connection credentials: %v", err)
	}

	// Create SSH client config
	sshConfig := &ssh.ClientConfig{
		User:            conn.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, verify host key
		Timeout:         10 * time.Second,
	}

	if conn.AuthType == "key" && privateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(privateKey))
		if err != nil {
			return fmt.Errorf("failed to parse private key: %v", err)
		}
		sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else if password != "" {
		sshConfig.Auth = []ssh.AuthMethod{ssh.Password(password)}
	} else {
		return fmt.Errorf("no authentication credentials provided")
	}

	// Connect to SSH server
	addr := fmt.Sprintf("%s:%d", conn.Host, conn.Port)
	log.Printf("TerminalService: Connecting to %s as %s", addr, conn.Username)

	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("ssh connection failed: %v", err)
	}
	defer sshClient.Close()

	// Create session
	session, err := sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create ssh session: %v", err)
	}
	defer session.Close()

	// Request PTY
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", 40, 80, modes); err != nil {
		return fmt.Errorf("request for pty failed: %v", err)
	}

	// Pipes
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("unable to setup stdin: %v", err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to setup stdout: %v", err)
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("unable to setup stderr: %v", err)
	}

	// Start shell
	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %v", err)
	}

	// Handle I/O
	errorChan := make(chan error, 3)

	// Custom Writer for WebSocket
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				if err != io.EOF {
					errorChan <- err
				}
				break
			}
			if err := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				errorChan <- err
				break
			}
		}
	}()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				if err != io.EOF {
					errorChan <- err
				}
				break
			}
			if err := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				errorChan <- err
				break
			}
		}
	}()

	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				errorChan <- err
				break
			}

			// Try to parse as JSON
			var wsMsg struct {
				Type string `json:"type"`
				Data string `json:"data"`
				Cols int    `json:"cols"`
				Rows int    `json:"rows"`
			}
			if err := json.Unmarshal(msg, &wsMsg); err == nil {
				switch wsMsg.Type {
				case "input":
					if _, err := stdin.Write([]byte(wsMsg.Data)); err != nil {
						errorChan <- err
						return
					}
				case "resize":
					session.WindowChange(wsMsg.Rows, wsMsg.Cols)
				}
			} else {
				// Raw message, write directly
				if _, err := stdin.Write(msg); err != nil {
					errorChan <- err
					break
				}
			}
		}
	}()

	// Wait for errors or completion
	select {
	case err := <-errorChan:
		log.Printf("TerminalService: Connection closed with error: %v", err)
		return err // Or nil if just closed
	}
}
