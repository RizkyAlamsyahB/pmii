package requests

// CreateUserRequest adalah DTO untuk request create user
type CreateUserRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	PhotoURI string `json:"photo_uri,omitempty"`
}

// UpdateUserRequest adalah DTO untuk request update user (Admin only)
type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Role     int    `json:"role" binding:"required,oneof=1 2"`
	PhotoURI string `json:"photo_uri,omitempty"`
	IsActive bool   `json:"is_active"`
	Password string `json:"password,omitempty" binding:"omitempty,min=8"`
}
