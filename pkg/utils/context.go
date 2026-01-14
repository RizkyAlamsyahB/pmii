package utils

import "context"

// Context keys for activity logging
type contextKey string

const (
	// ContextKeyUserID is the key for the authenticated user ID in context
	ContextKeyUserID contextKey = "user_id"
	// ContextKeyIPAddress is the key for the client IP address in context
	ContextKeyIPAddress contextKey = "ip_address"
	// ContextKeyUserAgent is the key for the client user agent in context
	ContextKeyUserAgent contextKey = "user_agent"
)

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

// WithRequestInfo adds IP address and user agent to context
func WithRequestInfo(ctx context.Context, ipAddress, userAgent string) context.Context {
	ctx = context.WithValue(ctx, ContextKeyIPAddress, ipAddress)
	ctx = context.WithValue(ctx, ContextKeyUserAgent, userAgent)
	return ctx
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(ContextKeyUserID).(int)
	return userID, ok
}

// GetIPAddress retrieves IP address from context
func GetIPAddress(ctx context.Context) string {
	if ip, ok := ctx.Value(ContextKeyIPAddress).(string); ok {
		return ip
	}
	return ""
}

// GetUserAgent retrieves user agent from context
func GetUserAgent(ctx context.Context) string {
	if ua, ok := ctx.Value(ContextKeyUserAgent).(string); ok {
		return ua
	}
	return ""
}
