package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"ssh-terminal-app/internal/config"
	"ssh-terminal-app/internal/models"
	"ssh-terminal-app/internal/repository"
	"ssh-terminal-app/internal/utils"

	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

type AuthService interface {
	Register(email, password, name string) (*models.User, error)
	Login(email, password string) (string, error)
	GoogleLogin(state string) string
	GoogleCallback(ctx context.Context, code, mode string) (string, error)
	GenerateToken(userID uint) (string, error)
}

type authService struct {
	repo        repository.UserRepository
	cfg         *config.Config
	googleOAuth *oauth2.Config
}

func NewAuthService(repo repository.UserRepository, cfg *config.Config, googleOAuth *oauth2.Config) AuthService {
	return &authService{
		repo:        repo,
		cfg:         cfg,
		googleOAuth: googleOAuth,
	}
}

func (s *authService) Register(email, password, name string) (*models.User, error) {
	// ... (same as before)
	existingUser, _ := s.repo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    email,
		Password: hashedPassword,
		Name:     name,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
	// ... (same as before)
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return s.GenerateToken(user.ID)
}

func (s *authService) GoogleLogin(state string) string {
	return s.googleOAuth.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *authService) GoogleCallback(ctx context.Context, code, mode string) (string, error) {
	token, err := s.googleOAuth.Exchange(ctx, code)
	if err != nil {
		return "", errors.New("failed to exchange token")
	}

	client := s.googleOAuth.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", errors.New("failed to get user info")
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return "", errors.New("failed to decode user info")
	}

	user, err := s.repo.FindByGoogleID(googleUser.ID)
	if err != nil {
		// User not found by Google ID
		user, err = s.repo.FindByEmail(googleUser.Email)
		if err != nil {
			// User does not exist
			if mode == "login" {
				return "", errors.New("User not found. Please register first.")
			}

			// Register mode
			newUser := &models.User{
				Email:    googleUser.Email,
				Name:     googleUser.Name,
				GoogleID: googleUser.ID,
			}
			if err := s.repo.Create(newUser); err != nil {
				return "", err
			}
			user = newUser
		} else {
			// Link account
			user.GoogleID = googleUser.ID
			if err := s.repo.Update(user); err != nil {
				return "", err
			}
		}
	}

	return s.GenerateToken(user.ID)
}

func (s *authService) GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
