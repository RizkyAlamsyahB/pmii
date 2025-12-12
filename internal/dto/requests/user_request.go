package requests

// CreateUserRequest adalah DTO untuk request create user
type CreateUserRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	PhotoURI string `json:"photo_uri,omitempty"`
}
