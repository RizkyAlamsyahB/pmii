package middleware

import (
	"net/http"

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
