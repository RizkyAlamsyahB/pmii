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
			FullName: user.Name,
			Email:    user.Email,
			Role:     getRoleName(user.Level),
			Status:   getStatusName(user.Status),
		})
	}

	// Response dengan total count
	response := responses.UserListResponse{
		Users: userList,
		Total: len(userList),
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Data user berhasil diambil", response))
}
