package domain

import "time"

// Contact represents contact information (Contact tab in Settings)
// This is a singleton table - only one record exists
type Contact struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Address       *string   `gorm:"type:text" json:"address,omitempty"`
	Email         *string   `gorm:"type:varchar(100)" json:"email,omitempty"`
	Phone         *string   `gorm:"type:varchar(50)" json:"phone,omitempty"`
	GoogleMapsURL *string   `gorm:"type:varchar(500)" json:"google_maps_url,omitempty"`
	UpdatedAt     time.Time `gorm:"default:now()" json:"updated_at"`
}

// TableName specifies the table name for Contact
func (Contact) TableName() string {
	return "contacts"
}
