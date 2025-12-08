package middleware

import (
	"net/http"
	"strings"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware memverifikasi JWT token dari header Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "Token tidak ditemukan"))
			c.Abort()
			return
		}

		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "Format token tidak valid"))
			c.Abort()
			return
		}

		token := parts[1]

		// Cek apakah token sudah di-blacklist (logout)
		if utils.IsBlacklisted(token) {
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "Token tidak valid atau sesi telah berakhir"))
			c.Abort()
			return
		}

		// Validasi token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse(401, "Token tidak valid atau kadaluarsa"))
			c.Abort()
			return
		}

		// Set user info ke context untuk digunakan di handler
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminOnly middleware untuk route yang hanya boleh diakses admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != "1" {
			c.JSON(http.StatusForbidden, responses.ErrorResponse(403, "Akses ditolak. Hanya admin yang diizinkan"))
			c.Abort()
			return
		}

		c.Next()
	}
}
