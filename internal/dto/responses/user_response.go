package responses

// UserListItem adalah DTO untuk item dalam list users (admin endpoint)
type UserListItem struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

// UserListResponse adalah response untuk GET /admin/users
type UserListResponse struct {
	Users []UserListItem `json:"users"`
	Total int            `json:"total"`
}

// UserProfileResponse adalah response untuk GET /user/profile
type UserProfileResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	Photo    string `json:"photo,omitempty"`
}
