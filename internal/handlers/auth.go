package handlers

import (
	"encoding/json"
	"net/http"

	"ssh-terminal-app/internal/config"
	"ssh-terminal-app/internal/service"
)

type AuthHandler struct {
	service service.AuthService
	cfg     *config.Config
}

func NewAuthHandler(service service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		service: service,
		cfg:     cfg,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, "Email, password and name are required", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(req.Email, req.Password, req.Name)
	if err != nil {
		if err.Error() == "user already exists" {
			http.Error(w, "User already exists", http.StatusConflict)
		} else {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if h.cfg.GoogleClientID == "" {
		http.Error(w, "Google OAuth not configured", http.StatusNotImplemented)
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "login"
	}

	url := h.service.GoogleLogin(mode)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in callback", http.StatusBadRequest)
		return
	}

	mode := r.URL.Query().Get("state")

	token, err := h.service.GoogleCallback(r.Context(), code, mode)
	if err != nil {
		// Handle specific errors like "User not found"
		if err.Error() == "User not found. Please register first." {
			frontendURL := h.cfg.FrontendURL + "/login?error=User not found. Please register first."
			http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	frontendURL := h.cfg.FrontendURL + "/auth/callback?token=" + token
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Placeholder for now. To be fully layered, need GetProfile in AuthService.
	// For this refactoring step, we focus on auth flow.
	w.WriteHeader(http.StatusOK)
}
