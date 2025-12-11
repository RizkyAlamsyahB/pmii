package domain

// PostTag represents the many-to-many relationship between posts and tags
type PostTag struct {
	PostID int `gorm:"primaryKey" json:"post_id"`
	TagID  int `gorm:"primaryKey" json:"tag_id"`
}

// TableName specifies the table name for PostTag
func (PostTag) TableName() string {
	return "post_tags"
}
