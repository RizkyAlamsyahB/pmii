package domain

import (
	"time"

	"gorm.io/gorm"
)

// Comment represents a comment on a blog post
type Comment struct {
	ID         int            `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID     int            `gorm:"not null" json:"post_id"`
	ParentID   *int           `json:"parent_id,omitempty"` // For nested comments
	UserID     *int           `json:"user_id,omitempty"`   // Null if guest comment
	GuestName  *string        `gorm:"type:varchar(100)" json:"guest_name,omitempty"`
	GuestEmail *string        `gorm:"type:varchar(100)" json:"guest_email,omitempty"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Status     CommentStatus  `gorm:"type:comment_status;default:'pending'" json:"status"`
	CreatedAt  time.Time      `gorm:"default:now()" json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Post    Post       `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Parent  *Comment   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	User    *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Replies []Comment  `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// TableName specifies the table name for Comment
func (Comment) TableName() string {
	return "comments"
}
