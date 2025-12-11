package domain

import "time"

// Testimonial represents a testimonial from a user or organization
type Testimonial struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Organization *string   `gorm:"type:varchar(100)" json:"organization,omitempty"`
	Position     *string   `gorm:"type:varchar(100)" json:"position,omitempty"`
	Content      string    `gorm:"type:text;not null" json:"content"`
	PhotoURI     *string   `gorm:"type:varchar(255)" json:"photo_uri,omitempty"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"default:now()" json:"created_at"`
}

// TableName specifies the table name for Testimonial
func (Testimonial) TableName() string {
	return "testimonials"
}
