package domain

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID            int            `gorm:"primaryKey;column:id" json:"id"`
	UserID        int            `gorm:"column:user_id;not null" json:"authorId"`
	CategoryID    int            `gorm:"column:category_id;not null" json:"categoryId"`
	Category      Category       `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
	Title         string         `gorm:"column:title;size:255;not null" json:"title"`
	Slug          string         `gorm:"column:slug;size:255;not null" json:"slug"`
	Excerpt       string         `gorm:"column:excerpt;type:text" json:"excerpt"`
	Content       string         `gorm:"column:content;type:text" json:"content"`
	FeaturedImage string         `gorm:"column:featured_image;size:255" json:"imageUrl"`
	Views         int            `gorm:"column:views;default:0" json:"views"`
	Status        int            `gorm:"column:status;default:1" json:"status"`
	Tags          []Tag          `gorm:"many2many:post_tags;foreignKey:ID;joinForeignKey:PostID;References:ID;joinReferences:TagID" json:"tags"`
	PublishedAt   time.Time      `gorm:"column:published_at" json:"publishedAt"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Post) TableName() string { return "posts" }
