package domain

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user account in the system
type User struct {
	ID           int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Role         int            `gorm:"not null;default:2" json:"role"` // 1=Admin, 2=Author
	FullName     string         `gorm:"type:varchar(100);not null" json:"full_name"`
	Email        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	PhotoURI     *string        `gorm:"type:varchar(255)" json:"photo_uri,omitempty"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == 1
}
