package domain

import "time"

// PostView tracks individual views of a post
type PostView struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    int       `gorm:"not null;index" json:"post_id"`
	IPAddress *string   `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent *string   `gorm:"type:text" json:"user_agent,omitempty"`
	ViewedAt  time.Time `gorm:"default:now()" json:"viewed_at"`

	// Relationships
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

// TableName specifies the table name for PostView
func (PostView) TableName() string {
	return "post_views"
}
