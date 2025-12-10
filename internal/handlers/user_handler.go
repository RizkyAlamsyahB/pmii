package handlers

import (
	"net/http"

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

// GetProfile handles GET /user/profile (Authenticated User)
// Menampilkan profil user yang sedang login
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Ambil user_id dari context (di-set oleh AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "Token tidak valid"))
		return
	}

	// Convert ke int
	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Terjadi kesalahan sistem"))
		return
	}

	// Get user by ID dari service
	user, err := h.userService.GetUserByID(id)
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
