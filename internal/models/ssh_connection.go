package models

import (
	"time"

	"gorm.io/gorm"
)

type SSHConnection struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint   `gorm:"not null;index" json:"user_id"`
	Name   string `gorm:"not null" json:"name"`
	Host   string `gorm:"not null" json:"host"`
	Port   int    `gorm:"not null;default:22" json:"port"`

	Username   string `gorm:"not null" json:"username"`
	Password   string `gorm:"" json:"-"`      // Encrypted, not exposed in JSON
	PrivateKey string `gorm:"" json:"-"`      // Encrypted, not exposed in JSON
	AuthType   string `gorm:"not null" json:"auth_type"` // "password" or "key"
}
