package requests

// CreateUserRequest adalah DTO untuk request create user
type CreateUserRequest struct {
	FullName string `json:"full_name" form:"full_name" binding:"required,min=2,max=100"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=8"`
	// PhotoURI akan dihandle terpisah menggunakan c.FormFile("photo")
}

// UpdateUserRequest adalah DTO untuk request update user (Admin only)
// Semua field opsional - hanya field yang ada di request yang akan di-update
type UpdateUserRequest struct {
	FullName *string `json:"full_name" form:"full_name" binding:"omitempty,min=2,max=100"`
	Email    *string `json:"email" form:"email" binding:"omitempty,email"`
	Role     *int    `json:"role" form:"role" binding:"omitempty,oneof=1 2"`
	IsActive *bool   `json:"is_active" form:"is_active"`
	Password *string `json:"password,omitempty" form:"password" binding:"omitempty,min=8"`
	// PhotoURI akan dihandle terpisah menggunakan c.FormFile("photo")
}
