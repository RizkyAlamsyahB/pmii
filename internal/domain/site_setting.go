package domain

import "time"

// SiteSetting represents global site configuration (Basic tab in Settings)
type SiteSetting struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	SiteName        *string   `gorm:"type:varchar(100)" json:"site_name,omitempty"`
	SiteTitle       *string   `gorm:"type:varchar(100)" json:"site_title,omitempty"`
	SiteDescription *string   `gorm:"type:text" json:"site_description,omitempty"`
	Favicon         *string   `gorm:"type:varchar(255)" json:"favicon,omitempty"`
	LogoHeader      *string   `gorm:"type:varchar(255)" json:"logo_header,omitempty"`
	LogoBig         *string   `gorm:"type:varchar(255)" json:"logo_big,omitempty"`
	FacebookURL     *string   `gorm:"type:varchar(255)" json:"facebook_url,omitempty"`
	TwitterURL      *string   `gorm:"type:varchar(255)" json:"twitter_url,omitempty"`
	LinkedinURL     *string   `gorm:"type:varchar(255)" json:"linkedin_url,omitempty"`
	InstagramURL    *string   `gorm:"type:varchar(255)" json:"instagram_url,omitempty"`
	YoutubeURL      *string   `gorm:"type:varchar(255)" json:"youtube_url,omitempty"`
	GithubURL       *string   `gorm:"type:varchar(255)" json:"github_url,omitempty"`
	UpdatedAt       time.Time `gorm:"default:now()" json:"updated_at"`
}

// TableName specifies the table name for SiteSetting
func (SiteSetting) TableName() string {
	return "site_settings"
}
