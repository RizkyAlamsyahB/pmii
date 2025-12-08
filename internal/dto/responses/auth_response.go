package responses

// LoginResponse adalah DTO untuk response login
type LoginResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

// UserDTO adalah DTO untuk data user (tanpa data sensitif)
type UserDTO struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
