package handlers

import (
	"net/http"
	"strconv"

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
