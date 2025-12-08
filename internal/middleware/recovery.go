package middleware

import (
	"net/http"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Recovery middleware untuk handle panic dan mengembalikan 500 response
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log error
				logger.Error.Printf("Panic recovered: %v", err)

				// Return 500 response
				c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Terjadi kesalahan server. Silakan coba lagi nanti"))
				c.Abort()
			}
		}()

		c.Next()
	}
}
