package domain

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID            int            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int            `gorm:"not null" json:"user_id"`
	CategoryID    int            `gorm:"not null" json:"category_id"`
	Title         string         `gorm:"type:varchar(255);not null" json:"title"`
	Slug          string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Excerpt       *string        `gorm:"type:text" json:"excerpt,omitempty"`
	Content       string         `gorm:"type:text;not null" json:"content"`
	FeaturedImage *string        `gorm:"type:varchar(255)" json:"featured_image,omitempty"`
	Status        PostStatus     `gorm:"type:post_status;not null;default:'draft'" json:"status"`
	PublishedAt   *time.Time     `json:"published_at,omitempty"`
	CreatedAt     time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags     []Tag    `gorm:"many2many:post_tags" json:"tags,omitempty"`
}

// TableName specifies the table name for Post
func (Post) TableName() string {
	return "posts"
}
