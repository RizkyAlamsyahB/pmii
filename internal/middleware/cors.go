package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS middleware untuk menghandle Cross-Origin Resource Sharing
func CORS(allowedOrigins string) gin.HandlerFunc {
	// Check if wildcard is used
	allowAll := strings.TrimSpace(allowedOrigins) == "*"

	// Parse comma-separated origins into slice
	origins := strings.Split(allowedOrigins, ",")
	// Trim whitespace from each origin
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return func(c *gin.Context) {
		requestOrigin := c.Request.Header.Get("Origin")

		// Check if request origin is in allowed list or if wildcard is used
		allowed := allowAll
		if !allowed {
			for _, origin := range origins {
				if origin == requestOrigin {
					allowed = true
					break
				}
			}
		}

		// Set CORS headers only for allowed origins
		if allowed {
			if allowAll {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				c.Writer.Header().Set("Access-Control-Allow-Origin", requestOrigin)
			}
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
