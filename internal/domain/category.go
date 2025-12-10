package domain

import "time"

type Category struct {
	ID          int       `gorm:"primaryKey;column:id" json:"id"`
	Name        string    `gorm:"column:name;size:100;not null" json:"name"`
	Slug        string    `gorm:"column:slug;size:100;not null" json:"slug"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (Category) TableName() string { return "categories" }
