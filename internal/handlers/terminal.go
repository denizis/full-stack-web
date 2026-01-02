package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"ssh-terminal-app/internal/config"
	"ssh-terminal-app/internal/service"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type TerminalHandler struct {
	service service.TerminalService
	cfg     *config.Config
}

func NewTerminalHandler(service service.TerminalService, cfg *config.Config) *TerminalHandler {
	return &TerminalHandler{
		service: service,
		cfg:     cfg,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *TerminalHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		h.sendError(ws, "Unauthorized: No token provided")
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		h.sendError(ws, "Unauthorized: Invalid token")
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		h.sendError(ws, "Unauthorized: Invalid token claims")
		return
	}
	userID := uint(userIDFloat)

	vars := mux.Vars(r)
	connID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.sendError(ws, "Invalid connection ID")
		return
	}

	if err := h.service.StartSession(ws, uint(connID), userID); err != nil {
		h.sendError(ws, fmt.Sprintf("Session error: %v", err))
		return
	}
}

func (h *TerminalHandler) sendError(ws *websocket.Conn, message string) {
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\x1b[31mError: %s\r\n\x1b[0m", message)))
	ws.Close()
}
