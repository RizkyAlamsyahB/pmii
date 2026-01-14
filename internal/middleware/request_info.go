package middleware

import (
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
	"github.com/gin-gonic/gin"
)

// RequestInfoMiddleware injects IP address and user agent into gin context
// These values can later be retrieved in handlers/services for activity logging
func RequestInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP (handles X-Forwarded-For for proxied requests)
		clientIP := c.ClientIP()

		// Get User-Agent header
		userAgent := c.GetHeader("User-Agent")

		// Store in gin context for later retrieval
		c.Set(string(utils.ContextKeyIPAddress), clientIP)
		c.Set(string(utils.ContextKeyUserAgent), userAgent)

		c.Next()
	}
}
