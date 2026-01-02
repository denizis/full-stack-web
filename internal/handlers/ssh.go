package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ssh-terminal-app/internal/config"
	"ssh-terminal-app/internal/middleware"
	"ssh-terminal-app/internal/service"

	"github.com/gorilla/mux"
)

type SSHHandler struct {
	service service.SSHService
	cfg     *config.Config
}

func NewSSHHandler(service service.SSHService, cfg *config.Config) *SSHHandler {
	return &SSHHandler{
		service: service,
		cfg:     cfg,
	}
}

// Response structs (DTOs)
type SSHConnectionResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	AuthType  string `json:"auth_type"`
	CreatedAt string `json:"created_at"`
}

func (h *SSHHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	connections, err := h.service.List(userID)
	if err != nil {
		http.Error(w, "Error fetching connections", http.StatusInternalServerError)
		return
	}

	response := make([]SSHConnectionResponse, len(connections))
	for i, conn := range connections {
		response[i] = SSHConnectionResponse{
			ID:        conn.ID,
			Name:      conn.Name,
			Host:      conn.Host,
			Port:      conn.Port,
			Username:  conn.Username,
			AuthType:  conn.AuthType,
			CreatedAt: conn.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *SSHHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid connection ID", http.StatusBadRequest)
		return
	}

	conn, err := h.service.Get(uint(id), userID)
	if err != nil {
		http.Error(w, "Connection not found", http.StatusNotFound)
		return
	}

	response := SSHConnectionResponse{
		ID:        conn.ID,
		Name:      conn.Name,
		Host:      conn.Host,
		Port:      conn.Port,
		Username:  conn.Username,
		AuthType:  conn.AuthType,
		CreatedAt: conn.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *SSHHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	var req service.SSHConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	conn, err := h.service.Create(userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := SSHConnectionResponse{
		ID:        conn.ID,
		Name:      conn.Name,
		Host:      conn.Host,
		Port:      conn.Port,
		Username:  conn.Username,
		AuthType:  conn.AuthType,
		CreatedAt: conn.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *SSHHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid connection ID", http.StatusBadRequest)
		return
	}

	var req service.SSHConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	conn, err := h.service.Update(uint(id), userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := SSHConnectionResponse{
		ID:        conn.ID,
		Name:      conn.Name,
		Host:      conn.Host,
		Port:      conn.Port,
		Username:  conn.Username,
		AuthType:  conn.AuthType,
		CreatedAt: conn.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *SSHHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid connection ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(uint(id), userID); err != nil {
		http.Error(w, "Error deleting connection", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
