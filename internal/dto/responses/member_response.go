package responses

import "time"

// MemberResponse adalah DTO untuk response member
type MemberResponse struct {
	ID          int            `json:"id"`
	FullName    string         `json:"fullName"`
	Position    string         `json:"position"`
	Department  string         `json:"department"`
	Photo       string         `json:"photo"`
	SocialLinks map[string]any `json:"socialLinks,omitempty"`
	IsActive    bool           `json:"isActive"`
	CreatedAt   time.Time      `json:"createdAt"`
}
