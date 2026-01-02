package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"" json:"-"`
	Name     string `gorm:"not null" json:"name"`
	GoogleID string `gorm:"uniqueIndex" json:"-"`

	SSHConnections []SSHConnection `gorm:"foreignKey:UserID" json:"-"`
}
