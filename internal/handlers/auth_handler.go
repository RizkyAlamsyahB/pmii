package handlers

import (
	"net/http"
	"strings"

	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthHandler handles HTTP requests untuk authentication
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler constructor untuk AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req requests.LoginRequest

	// Validasi Input
	if err := c.ShouldBindJSON(&req); err != nil {
		// Parse validation errors
		validationErrors := make(map[string][]string)

		if validationErr, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErr {
				fieldName := strings.ToLower(e.Field())
				var message string

				switch e.Tag() {
				case "required":
					switch fieldName {
					case "email":
						message = "Email wajib diisi"
					case "password":
						message = "Password wajib diisi"
					}
				case "email":
					message = "Format email tidak valid"
				}

				validationErrors[fieldName] = append(validationErrors[fieldName], message)
			}
		}

		c.JSON(http.StatusBadRequest, responses.ValidationErrorResponse(validationErrors))
		return
	}

	// Call service layer
	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "Email atau password salah"))
		return
	}

	// Convert domain.User → dto.UserDTO
	response := responses.LoginResponse{
		Token: token,
		User: responses.UserDTO{
			ID:       user.ID,
			FullName: user.FullName,
			Email:    user.Email,
			Role:     getRoleName(user.Role),
		},
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Login berhasil", response))
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get token dari header Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "Token tidak valid atau sesi telah berakhir"))
		return
	}

	// Parse token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "Token tidak valid atau sesi telah berakhir"))
		return
	}

	token := parts[1]

	// Logout (blacklist token)
	if err := h.authService.Logout(token); err != nil {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "Token tidak valid atau sesi telah berakhir"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Logout berhasil", nil))
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req requests.ChangePasswordRequest

	// Validasi input
	err := c.ShouldBindJSON(&req)
	if err != nil {
		validationErrors := make(map[string][]string)

		validationErr, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range validationErr {
				fieldName := strings.ToLower(e.Field())
				var message string

				switch e.Tag() {
				case "required":
					switch fieldName {
					case "oldpassword":
						message = "Password lama wajib diisi"
					case "newpassword":
						message = "Password baru wajib diisi"
					}
				case "min":
					message = "Password baru minimal 8 karakter"
				case "containsany":
					message = "Password baru harus mengandung minimal satu karakter spesial (!@#$%^&*)"
				}

				validationErrors[fieldName] = append(validationErrors[fieldName], message)
			}
		}

		c.JSON(http.StatusBadRequest, responses.ValidationErrorResponse(validationErrors))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "User tidak terautentikasi"))
		return
	}

	// Panggil service layer dengan userID langsung
	err = h.authService.ChangePassword(userID.(int), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Password berhasil diubah", nil))
}
