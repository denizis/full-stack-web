package main

import (
	"log"
	"net/http"
	"os"

	"ssh-terminal-app/internal/config"
	"ssh-terminal-app/internal/database"
	"ssh-terminal-app/internal/handlers"
	"ssh-terminal-app/internal/middleware"
	"ssh-terminal-app/internal/repository"
	"ssh-terminal-app/internal/service"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	// 1. Load configuration
	cfg := config.Load()

	// 2. Initialize database
	db, err := database.Initialize(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 3. Configure Google OAuth
	var googleOAuth *oauth2.Config
	if cfg.GoogleClientID != "" {
		googleOAuth = &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
	}

	// 4. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	sshRepo := repository.NewSSHRepository(db)

	// 5. Initialize Services
	authService := service.NewAuthService(userRepo, cfg, googleOAuth)
	sshService := service.NewSSHService(sshRepo, cfg)
	terminalService := service.NewTerminalService(sshService)

	// 6. Initialize Handlers with Services
	authHandler := handlers.NewAuthHandler(authService, cfg)
	sshHandler := handlers.NewSSHHandler(sshService, cfg)
	terminalHandler := handlers.NewTerminalHandler(terminalService, cfg)

	// 7. Setup Router
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.CORS)

	// Public routes
	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/google", authHandler.GoogleLogin).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/auth/google/callback", authHandler.GoogleCallback).Methods("GET", "OPTIONS")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.Auth(cfg.JWTSecret))

	protected.HandleFunc("/ssh", sshHandler.List).Methods("GET", "OPTIONS")
	protected.HandleFunc("/ssh", sshHandler.Create).Methods("POST", "OPTIONS")
	protected.HandleFunc("/ssh/{id}", sshHandler.Get).Methods("GET", "OPTIONS")
	protected.HandleFunc("/ssh/{id}", sshHandler.Update).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/ssh/{id}", sshHandler.Delete).Methods("DELETE", "OPTIONS")
	protected.HandleFunc("/auth/me", authHandler.Me).Methods("GET", "OPTIONS")

	// WebSocket route for terminal (handshakes auth internally via query token)
	r.HandleFunc("/ws/terminal/{id}", terminalHandler.HandleWebSocket)

	// Serve static files for frontend
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/dist")))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default changed to 3000 to match previous setup
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
