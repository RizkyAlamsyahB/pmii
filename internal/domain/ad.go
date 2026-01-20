package domain

import "time"

// AdPage represents the page where ad is displayed
type AdPage string

const (
	AdPageLanding       AdPage = "landing"
	AdPageNews          AdPage = "news"
	AdPageOpini         AdPage = "opini"
	AdPageLifeAtPMII    AdPage = "life_at_pmii"
	AdPageIslamic       AdPage = "islamic"
	AdPageDetailArticle AdPage = "detail_article"
)

// ValidAdPages returns all valid ad pages
func ValidAdPages() []AdPage {
	return []AdPage{
		AdPageLanding,
		AdPageNews,
		AdPageOpini,
		AdPageLifeAtPMII,
		AdPageIslamic,
		AdPageDetailArticle,
	}
}

// IsValidAdPage checks if the given page is valid
func IsValidAdPage(page string) bool {
	for _, p := range ValidAdPages() {
		if string(p) == page {
			return true
		}
	}
	return false
}

// Ad represents an advertisement slot in the system
type Ad struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Page       AdPage    `gorm:"type:ad_page;not null" json:"page"`
	Slot       int       `gorm:"not null" json:"slot"`
	ImageURL   *string   `gorm:"type:varchar(500)" json:"image_url,omitempty"`
	Resolution string    `gorm:"type:varchar(20);not null" json:"resolution"`
	CreatedAt  time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:now()" json:"updated_at"`
}

// TableName specifies the table name for Ad
func (Ad) TableName() string {
	return "ads"
}

// GetPageDisplayName returns human-readable page name
func (a *Ad) GetPageDisplayName() string {
	switch a.Page {
	case AdPageLanding:
		return "Landing Page"
	case AdPageNews:
		return "News Page"
	case AdPageOpini:
		return "Opini Page"
	case AdPageLifeAtPMII:
		return "Life at PMII Page"
	case AdPageIslamic:
		return "Islamic Page"
	case AdPageDetailArticle:
		return "Detail Article Page"
	default:
		return string(a.Page)
	}
}

// GetSlotName returns the slot display name (e.g., "ADS 1 Landing Page")
func (a *Ad) GetSlotName() string {
	return "ADS " + string(rune('0'+a.Slot)) + " " + a.GetPageDisplayName()
}
