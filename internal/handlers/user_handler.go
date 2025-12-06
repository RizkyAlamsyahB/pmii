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

// GetProfile handles GET /user/profile (Authenticated User)
// Menampilkan profil user yang sedang login
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Ambil user_id dari context (di-set oleh AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "Token tidak valid"))
		return
	}

	// Convert ke uint
	id, ok := userID.(uint)
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
	profile := responses.UserProfileResponse{
		ID:       user.ID,
		FullName: user.Name,
		Email:    user.Email,
		Role:     getRoleName(user.Level),
		Status:   getStatusName(user.Status),
		Photo:    user.Photo, // Bisa di-transform ke URL jika diperlukan
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Profil berhasil diambil", profile))
}
