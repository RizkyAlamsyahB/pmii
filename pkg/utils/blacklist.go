package utils

import (
	"sync"
	"time"
)

// TokenBlacklist menyimpan token yang sudah di-logout
type TokenBlacklist struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

var blacklist *TokenBlacklist

// InitBlacklist inisialisasi token blacklist
func InitBlacklist() {
	blacklist = &TokenBlacklist{
		tokens: make(map[string]time.Time),
	}

	// Cleanup expired tokens setiap 1 jam
	go blacklist.cleanupExpired()
}

// AddToBlacklist menambahkan token ke blacklist
func AddToBlacklist(token string, expiresAt time.Time) {
	blacklist.mu.Lock()
	defer blacklist.mu.Unlock()
	blacklist.tokens[token] = expiresAt
}

// IsBlacklisted mengecek apakah token sudah di-blacklist
func IsBlacklisted(token string) bool {
	blacklist.mu.Lock()
	defer blacklist.mu.Unlock()

	expiresAt, exists := blacklist.tokens[token]
	if !exists {
		return false
	}

	// Jika token sudah expired, hapus dari blacklist
	if time.Now().After(expiresAt) {
		delete(blacklist.tokens, token)
		return false
	}

	return true
}

// cleanupExpired menghapus token yang sudah expired dari blacklist
func (b *TokenBlacklist) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		b.mu.Lock()
		now := time.Now()
		for token, expiresAt := range b.tokens {
			if now.After(expiresAt) {
				delete(b.tokens, token)
			}
		}
		b.mu.Unlock()
	}
}
