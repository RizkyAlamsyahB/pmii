package domain

import "time"

// SiteSetting represents global site configuration
type SiteSetting struct {
	ID              int            `gorm:"primaryKey;autoIncrement" json:"id"`
	SiteName        *string        `gorm:"type:varchar(100)" json:"site_name,omitempty"`
	SiteDescription *string        `gorm:"type:text" json:"site_description,omitempty"`
	LogoHeader      *string        `gorm:"type:varchar(255)" json:"logo_header,omitempty"`
	LogoFooter      *string        `gorm:"type:varchar(255)" json:"logo_footer,omitempty"`
	Favicon         *string        `gorm:"type:varchar(255)" json:"favicon,omitempty"`
	SocialLinks     map[string]any `gorm:"type:jsonb" json:"social_links,omitempty"`
	ContactInfo     map[string]any `gorm:"type:jsonb" json:"contact_info,omitempty"`
	UpdatedAt       time.Time      `gorm:"default:now()" json:"updated_at"`
}

// TableName specifies the table name for SiteSetting
func (SiteSetting) TableName() string {
	return "site_settings"
}
