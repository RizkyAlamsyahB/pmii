package responses

import "time"

// PublicAboutPageResponse adalah response untuk halaman About publik
// Menggabungkan data About dan Members per department
type PublicAboutPageResponse struct {
	About       PublicAboutResponse         `json:"about"`
	Departments []DepartmentMembersResponse `json:"departments"`
}

// PublicAboutResponse adalah data about untuk public
type PublicAboutResponse struct {
	History  *string `json:"history,omitempty"`
	Vision   *string `json:"vision,omitempty"`
	Mission  *string `json:"mission,omitempty"`
	ImageUrl string  `json:"imageUrl,omitempty"`
	VideoURL *string `json:"videoUrl,omitempty"`
}

// DepartmentMembersResponse adalah response members per department
type DepartmentMembersResponse struct {
	Department      string                `json:"department"`
	DepartmentLabel string                `json:"departmentLabel"`
	Members         PublicMembersResponse `json:"members"`
}

// PublicMemberResponse adalah data member untuk public
type PublicMemberResponse struct {
	ID          int            `json:"id"`
	FullName    string         `json:"fullName"`
	Position    string         `json:"position"`
	Photo       string         `json:"photo"`
	SocialLinks map[string]any `json:"socialLinks,omitempty"`
}

// PublicMembersResponse adalah list members dengan pagination
type PublicMembersResponse struct {
	Data       []PublicMemberResponse `json:"data"`
	Pagination PaginationMeta         `json:"pagination"`
}

// PaginationMeta metadata pagination
type PaginationMeta struct {
	Page     int   `json:"page"`
	Limit    int   `json:"limit"`
	Total    int64 `json:"total"`
	LastPage int   `json:"lastPage"`
}

// PublicMemberListResponse adalah response untuk list members saja (endpoint terpisah)
type PublicMemberListResponse struct {
	ID          int            `json:"id"`
	FullName    string         `json:"fullName"`
	Position    string         `json:"position"`
	Department  string         `json:"department"`
	Photo       string         `json:"photo"`
	SocialLinks map[string]any `json:"socialLinks,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
}
