package domain

// Tag represents a tag that can be attached to posts
type Tag struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(50);not null" json:"name"`
	Slug string `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`
}

// TableName specifies the table name for Tag
func (Tag) TableName() string {
	return "tags"
}
