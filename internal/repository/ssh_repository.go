package repository

import (
	"ssh-terminal-app/internal/models"

	"gorm.io/gorm"
)

// SSHRepository defines the interface for SSH connection data access
type SSHRepository interface {
	Create(conn *models.SSHConnection) error
	ListByUserID(userID uint) ([]models.SSHConnection, error)
	GetByID(id uint, userID uint) (*models.SSHConnection, error)
	Update(conn *models.SSHConnection) error
	Delete(id uint, userID uint) error
}

// sshRepository implements SSHRepository using GORM
type sshRepository struct {
	db *gorm.DB
}

// NewSSHRepository creates a new SSHRepository instance
func NewSSHRepository(db *gorm.DB) SSHRepository {
	return &sshRepository{db: db}
}

func (r *sshRepository) Create(conn *models.SSHConnection) error {
	return r.db.Create(conn).Error
}

func (r *sshRepository) ListByUserID(userID uint) ([]models.SSHConnection, error) {
	var connections []models.SSHConnection
	err := r.db.Where("user_id = ?", userID).Find(&connections).Error
	return connections, err
}

func (r *sshRepository) GetByID(id uint, userID uint) (*models.SSHConnection, error) {
	var conn models.SSHConnection
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&conn).Error
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

func (r *sshRepository) Update(conn *models.SSHConnection) error {
	return r.db.Save(conn).Error
}

func (r *sshRepository) Delete(id uint, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.SSHConnection{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
