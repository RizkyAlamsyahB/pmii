package domain

type Tag struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:50;not null" json:"name"`
	Slug string `gorm:"column:slug;size:50;not null" json:"slug"`
}

func (Tag) TableName() string { return "tags" }
