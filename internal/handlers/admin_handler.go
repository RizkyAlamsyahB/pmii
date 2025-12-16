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
