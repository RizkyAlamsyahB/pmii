package domain

import (
	"time"

	"gorm.io/gorm"
)

// Post merepresentasikan tabel tbl_post
type Post struct {
	ID          int       `gorm:"primaryKey;column:post_id" json:"id"` // Map ke post_id
	Title       string    `gorm:"column:post_title" json:"title"`
	Description string    `gorm:"column:post_description;type:text" json:"excerpt"` // Map description ke json excerpt
	Content     string    `gorm:"column:post_contents;type:text" json:"content"`
	Image       string    `gorm:"column:post_image" json:"imageUrl"` // Map ke imageUrl sesuai OpenAPI
	Date        time.Time `gorm:"column:post_date;autoCreateTime" json:"publishedAt"`
	UpdatedAt   time.Time `gorm:"column:post_last_update;autoUpdateTime" json:"updatedAt"`
	CategoryID  int       `gorm:"column:post_category_id" json:"categoryId"`
	Tags        string    `gorm:"column:post_tags" json:"tags"`
	Slug        string    `gorm:"column:post_slug" json:"slug"`
	Status      int       `gorm:"column:post_status;default:1" json:"status"`
	Views       int       `gorm:"column:post_views;default:0" json:"views"`
	UserID      int       `gorm:"column:post_user_id" json:"authorId"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName memaksa GORM menggunakan nama tabel 'tbl_post' bukan 'posts'
func (Post) TableName() string {
	return "tbl_post"
}
