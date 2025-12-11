package domain

import "time"

// Category represents a blog post category
type Category struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "categories"
}
