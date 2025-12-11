package handlers

import (
	"net/http"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

// AdminHandler handles HTTP requests untuk admin endpoints
type AdminHandler struct {
	userService service.UserService
}

// NewAdminHandler constructor untuk AdminHandler
func NewAdminHandler(userService service.UserService) *AdminHandler {
	return &AdminHandler{userService: userService}
}

// GetDashboard handles GET /admin/dashboard (Admin Only)
// Menampilkan welcome message untuk admin dashboard
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	// Get user info dari context (di-set oleh AuthMiddleware)
	userID, _ := c.Get("user_id")

	response := gin.H{
		"id":      userID,
		"role":    "admin",
		"message": "Welcome to Admin Dashboard",
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Dashboard admin berhasil diakses", response))
}

// GetAllUsers handles GET /admin/users (Admin Only)
// Menampilkan list semua user di sistem
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
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
