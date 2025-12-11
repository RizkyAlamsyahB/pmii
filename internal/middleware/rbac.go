package middleware

import (
	"net/http"
	"strconv"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/gin-gonic/gin"
)

// RequireRole middleware untuk validasi role user
// Harus dipanggil setelah AuthMiddleware() karena depend on user_role di context
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil user_role dari context (di-set oleh AuthMiddleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "Akses ditolak: Token tidak valid"))
			c.Abort()
			return
		}

		// Convert ke string
		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Terjadi kesalahan sistem"))
			c.Abort()
			return
		}

		// Validasi role
		if role != requiredRole {
			c.JSON(http.StatusForbidden, responses.ErrorResponse(403, "Forbidden: Anda tidak memiliki hak akses untuk resource ini"))
			c.Abort()
			return
		}

		// Role valid, lanjut ke handler
		c.Next()
	}
}

// RequireAnyRole middleware untuk validasi multiple roles (OR condition)
// User harus memiliki salah satu dari roles yang di-allowed
func RequireAnyRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil user_role dari context
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "Akses ditolak: Token tidak valid"))
			c.Abort()
			return
		}

		// Convert ke string
		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Terjadi kesalahan sistem"))
			c.Abort()
			return
		}

		// Cek apakah role user termasuk dalam allowed roles
		hasAccess := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, responses.ErrorResponse(403, "Forbidden: Anda tidak memiliki hak akses untuk resource ini"))
			c.Abort()
			return
		}

		// Role valid, lanjut ke handler
		c.Next()
	}
}

// RequireOwnerOrAdmin middleware untuk validasi akses resource milik sendiri atau admin
// Admin (role "1") dapat mengakses semua resource
// User lain hanya dapat mengakses resource milik sendiri (berdasarkan URL param)
// paramName adalah nama URL parameter yang berisi ID resource (contoh: "id")
func RequireOwnerOrAdmin(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user_role from context (di-set oleh AuthMiddleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "Akses ditolak: Token tidak valid"))
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Terjadi kesalahan sistem"))
			c.Abort()
			return
		}

		// Admin (role "1") can access any resource
		if role == "1" {
			c.Next()
			return
		}

		// For non-admin: check if requested ID matches user's own ID
		requestedID, err := strconv.Atoi(c.Param(paramName))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID tidak valid"))
			c.Abort()
			return
		}

		tokenUserID, _ := c.Get("user_id")
		currentUserID, ok := tokenUserID.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Terjadi kesalahan sistem"))
			c.Abort()
			return
		}

		if currentUserID != requestedID {
			c.JSON(http.StatusForbidden, responses.ErrorResponse(403, "Anda tidak memiliki akses ke data user ini"))
			c.Abort()
			return
		}

		// Access valid, lanjut ke handler
		c.Next()
	}
}
