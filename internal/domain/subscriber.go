package domain

import "time"

// Subscriber represents an email newsletter subscriber
type Subscriber struct {
	ID                int              `gorm:"primaryKey;autoIncrement" json:"id"`
	Email             string           `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Status            SubscriberStatus `gorm:"type:subscriber_status;default:'active'" json:"status"`
	VerificationToken *string          `gorm:"type:varchar(64)" json:"verification_token,omitempty"`
	VerifiedAt        *time.Time       `json:"verified_at,omitempty"`
	CreatedAt         time.Time        `gorm:"default:now()" json:"created_at"`
}

// TableName specifies the table name for Subscriber
func (Subscriber) TableName() string {
	return "subscribers"
}
