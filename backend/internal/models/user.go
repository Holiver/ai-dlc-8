package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	FullName          string    `gorm:"size:100;not null" json:"full_name"`
	Email             string    `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone             string    `gorm:"size:20;not null" json:"phone"`
	PasswordHash      string    `gorm:"size:255;not null" json:"-"`
	Role              string    `gorm:"type:enum('employee','admin');default:'employee'" json:"role"`
	PointsBalance     int       `gorm:"default:0" json:"points_balance"`
	IsFirstLogin      bool      `gorm:"default:true" json:"is_first_login"`
	IsActive          bool      `gorm:"default:true" json:"is_active"`
	PreferredLanguage string    `gorm:"size:10;default:'zh'" json:"preferred_language"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
