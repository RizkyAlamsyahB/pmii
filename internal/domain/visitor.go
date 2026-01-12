package domain

import "time"

// Visitor tracks unique daily visitors by IP address
// Each IP is recorded only once per day
type Visitor struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	IPAddress string    `gorm:"type:varchar(50);not null" json:"ip_address"`
	Platform  *string   `gorm:"type:varchar(255)" json:"platform,omitempty"`
	VisitedAt time.Time `gorm:"default:now()" json:"visited_at"`
}

// TableName specifies the table name for Visitor
func (Visitor) TableName() string {
	return "visitors"
}
