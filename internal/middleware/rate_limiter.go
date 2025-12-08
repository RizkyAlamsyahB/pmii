package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Visitor menyimpan rate limiter per IP
type Visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter menyimpan semua visitor
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter membuat instance rate limiter baru
// r = jumlah request per detik, b = burst capacity
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     r,
		burst:    b,
	}

	// Cleanup visitor yang sudah tidak aktif setiap 5 menit
	go rl.cleanupVisitors()

	return rl
}

// getVisitor mendapatkan atau membuat visitor baru untuk IP
func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	visitor, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = &Visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	visitor.lastSeen = time.Now()
	return visitor.limiter
}

// cleanupVisitors menghapus visitor yang tidak aktif lebih dari 5 menit
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(5 * time.Minute)

		rl.mu.Lock()
		for ip, visitor := range rl.visitors {
			if time.Since(visitor.lastSeen) > 5*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Limit adalah middleware untuk rate limiting
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getVisitor(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, responses.ErrorResponse(429, "Terlalu banyak request. Silakan coba lagi nanti"))
			c.Abort()
			return
		}

		c.Next()
	}
}
