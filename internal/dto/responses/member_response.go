package responses

// MemberResponse adalah DTO untuk response member list (admin table)
type MemberResponse struct {
	ID          int            `json:"id"`
	FullName    string         `json:"fullName"`
	Position    string         `json:"position"`
	Photo       string         `json:"photo"`
	SocialLinks map[string]any `json:"socialLinks,omitempty"`
}

// MemberDetailResponse adalah DTO untuk response member detail (admin edit form)
type MemberDetailResponse struct {
	ID          int            `json:"id"`
	FullName    string         `json:"fullName"`
	Position    string         `json:"position"`
	Department  string         `json:"department"`
	Photo       string         `json:"photo"`
	SocialLinks map[string]any `json:"socialLinks,omitempty"`
	IsActive    bool           `json:"isActive"`
}
