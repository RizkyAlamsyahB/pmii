package domain

import "time"

// Message represents a contact form message
type Message struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);not null" json:"email"`
	Subject   *string   `gorm:"type:varchar(200)" json:"subject,omitempty"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
}

// TableName specifies the table name for Message
func (Message) TableName() string {
	return "messages"
}
