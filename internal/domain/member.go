package domain

import "time"

// Member represents an organization member (e.g., team member, staff)
type Member struct {
	ID          int            `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName    string         `gorm:"type:varchar(100);not null" json:"full_name"`
	Position    string         `gorm:"type:varchar(100);not null" json:"position"`
	PhotoURI    *string        `gorm:"type:varchar(255)" json:"photo_uri,omitempty"`
	SocialLinks map[string]any `gorm:"type:jsonb;serializer:json" json:"social_links,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `gorm:"default:now()" json:"created_at"`
}

// TableName specifies the table name for Member
func (Member) TableName() string {
	return "members"
}
