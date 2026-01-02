package service

import (
	"errors"

	"ssh-terminal-app/internal/config"
	"ssh-terminal-app/internal/models"
	"ssh-terminal-app/internal/repository"
	"ssh-terminal-app/internal/utils"
)

type SSHService interface {
	Create(userID uint, req SSHConnectionRequest) (*models.SSHConnection, error)
	List(userID uint) ([]models.SSHConnection, error)
	Get(id, userID uint) (*models.SSHConnection, error)
	Update(id, userID uint, req SSHConnectionRequest) (*models.SSHConnection, error)
	Delete(id, userID uint) error
	// DecryptCredentials helps retrieving raw password/key for connection
	GetDecryptedCredentials(id, userID uint) (string, string, *models.SSHConnection, error)
}

type sshService struct {
	repo repository.SSHRepository
	cfg  *config.Config
}

func NewSSHService(repo repository.SSHRepository, cfg *config.Config) SSHService {
	return &sshService{
		repo: repo,
		cfg:  cfg,
	}
}

// SSHConnectionRequest DTO
type SSHConnectionRequest struct {
	Name       string
	Host       string
	Port       int
	Username   string
	Password   string
	PrivateKey string
	AuthType   string
}

func (s *sshService) Create(userID uint, req SSHConnectionRequest) (*models.SSHConnection, error) {
	if req.Name == "" || req.Host == "" || req.Username == "" {
		return nil, errors.New("name, host, and username are required")
	}

	if req.Port == 0 {
		req.Port = 22
	}
	if req.AuthType == "" {
		req.AuthType = "password"
	}

	var encryptedPassword, encryptedKey string
	var err error

	if req.Password != "" {
		encryptedPassword, err = utils.Encrypt(req.Password, s.cfg.EncryptionKey)
		if err != nil {
			return nil, err
		}
	}

	if req.PrivateKey != "" {
		encryptedKey, err = utils.Encrypt(req.PrivateKey, s.cfg.EncryptionKey)
		if err != nil {
			return nil, err
		}
	}

	conn := &models.SSHConnection{
		UserID:     userID,
		Name:       req.Name,
		Host:       req.Host,
		Port:       req.Port,
		Username:   req.Username,
		Password:   encryptedPassword,
		PrivateKey: encryptedKey,
		AuthType:   req.AuthType,
	}

	if err := s.repo.Create(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

func (s *sshService) List(userID uint) ([]models.SSHConnection, error) {
	return s.repo.ListByUserID(userID)
}

func (s *sshService) Get(id, userID uint) (*models.SSHConnection, error) {
	return s.repo.GetByID(id, userID)
}

func (s *sshService) Update(id, userID uint, req SSHConnectionRequest) (*models.SSHConnection, error) {
	conn, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		conn.Name = req.Name
	}
	if req.Host != "" {
		conn.Host = req.Host
	}
	if req.Port != 0 {
		conn.Port = req.Port
	}
	if req.Username != "" {
		conn.Username = req.Username
	}
	if req.AuthType != "" {
		conn.AuthType = req.AuthType
	}

	if req.Password != "" {
		encrypted, err := utils.Encrypt(req.Password, s.cfg.EncryptionKey)
		if err != nil {
			return nil, err
		}
		conn.Password = encrypted
	}

	if req.PrivateKey != "" {
		encrypted, err := utils.Encrypt(req.PrivateKey, s.cfg.EncryptionKey)
		if err != nil {
			return nil, err
		}
		conn.PrivateKey = encrypted
	}

	if err := s.repo.Update(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

func (s *sshService) Delete(id, userID uint) error {
	return s.repo.Delete(id, userID)
}

func (s *sshService) GetDecryptedCredentials(id, userID uint) (string, string, *models.SSHConnection, error) {
	conn, err := s.repo.GetByID(id, userID)
	if err != nil {
		return "", "", nil, err
	}

	var password, privateKey string

	if conn.Password != "" {
		password, err = utils.Decrypt(conn.Password, s.cfg.EncryptionKey)
		if err != nil {
			return "", "", nil, err
		}
	}

	if conn.PrivateKey != "" {
		privateKey, err = utils.Decrypt(conn.PrivateKey, s.cfg.EncryptionKey)
		if err != nil {
			return "", "", nil, err
		}
	}

	return password, privateKey, conn, nil
}
