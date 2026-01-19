package domain

import "time"

// About represents the about page content
type About struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     *string   `gorm:"type:varchar(255)" json:"title,omitempty"`     // Page title
	History   *string   `gorm:"type:text" json:"history,omitempty"`           // Brief history
	Vision    *string   `gorm:"type:text" json:"vision,omitempty"`            // Vision statement
	Mission   *string   `gorm:"type:text" json:"mission,omitempty"`           // Mission/Goals
	VideoURL  *string   `gorm:"type:varchar(255)" json:"video_url,omitempty"` // YouTube profile link
	UpdatedAt time.Time `gorm:"default:now()" json:"updated_at"`
}

// TableName specifies the table name for About
func (About) TableName() string {
	return "about"
}
