package domain

import "time"

type Inbox struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID   int       `gorm:"not null" json:"sender_id"`
	ReceiverID int       `gorm:"not null" json:"receiver_id"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	IsRead     bool      `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time `gorm:"default:now()" json:"created_at"`

	// Relasi ke User
	Sender   User `gorm:"foreignKey:SenderID" json:"sender"`
	Receiver User `gorm:"foreignKey:ReceiverID" json:"receiver"`
}

func (Inbox) TableName() string {
	return "inbox"
}
