package handlers

import (
	"net/http"
	"strconv"

	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests untuk user endpoints
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler constructor untuk UserHandler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetDashboard handles GET /user/dashboard (Authenticated User)
// Menampilkan welcome message untuk user dashboard
func (h *UserHandler) GetDashboard(c *gin.Context) {
	// Get user info dari context (di-set oleh AuthMiddleware)
	userID, _ := c.Get("user_id")

	response := gin.H{
		"id":      userID,
		"role":    "user",
		"message": "Welcome to User Dashboard",
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Dashboard user berhasil diakses", response))
}

// GetAllUsers handles GET /users (Admin Only)
// Menampilkan list semua user di sistem
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// Get all users dari service
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data user"))
		return
	}

	// Convert domain.User ke UserListItem DTO
	userList := make([]responses.UserListItem, 0, len(users))
	for _, user := range users {
		userList = append(userList, responses.UserListItem{
			ID:       user.ID,
			FullName: user.FullName,
			Email:    user.Email,
			Role:     getRoleName(user.Role),
			Status:   getStatusName(user.IsActive),
		})
	}

	// Response dengan total count
	response := responses.UserListResponse{
		Users: userList,
		Total: len(userList),
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Data user berhasil diambil", response))
}

// GetUserByID handles GET /users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Parse URL param :id
	requestedID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID tidak valid"))
		return
	}

	// Get user by ID dari service (access control already done by middleware)
	user, err := h.userService.GetUserByID(requestedID)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "User tidak ditemukan"))
		return
	}

	// Convert domain.User ke UserProfileResponse DTO
	photo := ""
	if user.PhotoURI != nil {
		photo = *user.PhotoURI
	}
	profile := responses.UserProfileResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     getRoleName(user.Role),
		Status:   getStatusName(user.IsActive),
		Photo:    photo,
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Profil berhasil diambil", profile))
}

// CreateUser handles POST /users (Admin Only)
// Membuat user baru di sistem
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req requests.CreateUserRequest

	// Bind dan validasi request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ValidationErrorResponse(err.Error()))
		return
	}

	// Create user via service
	user, err := h.userService.CreateUser(&req)
	if err != nil {
		// Handle specific errors
		switch err.Error() {
		case "email sudah terdaftar":
			c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Email sudah terdaftar"))
		case "password harus kombinasi huruf dan angka":
			c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Password harus kombinasi huruf dan angka"))
		default:
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal membuat user"))
		}
		return
	}

	// Convert domain.User ke response DTO
	response := responses.UserListItem{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     getRoleName(user.Role),
		Status:   getStatusName(user.IsActive),
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(201, "User berhasil dibuat", response))
}

// UpdateUserByID handles PUT /users/:id (Admin Only)
// Mengupdate data user berdasarkan ID
func (h *UserHandler) UpdateUserByID(c *gin.Context) {
	// Parse URL param :id
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID tidak valid"))
		return
	}

	var req requests.UpdateUserRequest

	// Bind dan validasi request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ValidationErrorResponse(err.Error()))
		return
	}

	// Update user via service
	user, err := h.userService.UpdateUser(userID, &req)
	if err != nil {
		// Handle specific errors
		switch err.Error() {
		case "user tidak ditemukan":
			c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "User tidak ditemukan"))
		case "email sudah digunakan user lain":
			c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Email sudah digunakan user lain"))
		case "password harus kombinasi huruf dan angka":
			c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Password harus kombinasi huruf dan angka"))
		default:
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengupdate user"))
		}
		return
	}

	// Convert domain.User ke response DTO
	response := responses.UserListItem{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     getRoleName(user.Role),
		Status:   getStatusName(user.IsActive),
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "User berhasil diupdate", response))
}
