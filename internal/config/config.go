package config

import "os"

type Config struct {
	DatabasePath       string
	JWTSecret          string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	EncryptionKey      string
	FrontendURL        string
}

func Load() *Config {
	return &Config{
		DatabasePath:       getEnv("DATABASE_PATH", "./ssh_terminal.db"),
		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
		EncryptionKey:      getEnv("ENCRYPTION_KEY", "a-32-byte-encryption-key-here!!"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
