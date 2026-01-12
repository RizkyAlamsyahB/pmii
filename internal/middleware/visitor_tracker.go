package middleware

import (
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/gin-gonic/gin"
)

// VisitorTracker middleware untuk mencatat unique visitor per hari
// Setiap IP hanya dicatat 1x per hari
func VisitorTracker(visitorRepo repository.VisitorRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP (support proxy/load balancer)
		clientIP := c.ClientIP()

		// Get platform/user agent
		platform := c.GetHeader("User-Agent")

		// Record visit (akan skip jika sudah ada record hari ini)
		go func() {
			_ = visitorRepo.RecordVisit(clientIP, platform)
		}()

		c.Next()
	}
}
